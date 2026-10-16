//go:build agora_video_observer_test

package agoraservice

import (
	"fmt"
	"testing"
)

func newTestVideoObserver(t *testing.T, observer *VideoFrameObserver) *RtcConnection {
	t.Helper()
	testVideoObserverReset()
	conn := testVideoObserverConnection()
	t.Cleanup(func() {
		testVideoObserverFail(0, 0, 0)
		testVideoObserverFreeConnection(conn)
	})
	if ret := conn.RegisterVideoFrameObserver(observer); ret != 0 {
		t.Fatalf("register: %d", ret)
	}
	return conn
}

func TestVideoFrameObserverDefaults(t *testing.T) {
	conn := newTestVideoObserver(t, &VideoFrameObserver{})
	want := []int{0, int(VideoPixelI420), 0, 0, 3}
	for slot, value := range want {
		if got := testVideoObserverQuery(conn.cVideoObserver, slot); got != value {
			t.Errorf("query %d = %d, want %d", slot, got, value)
		}
	}
	for slot := 0; slot < 5; slot++ {
		if ret, _ := testVideoObserverFrame(conn.cVideoObserver, slot, 0, false, VideoPixelI420); ret != 1 {
			t.Errorf("unset frame callback %d returned %d", slot, ret)
		}
		if ret, _ := testVideoObserverFrame(conn.cVideoObserver, slot, 0, true, VideoPixelI420); ret != 0 {
			t.Errorf("nil frame %d returned %d", slot, ret)
		}
	}
}

func TestVideoFrameObserverCallbacks(t *testing.T) {
	calls := [5]int{}
	accept := false
	checkFrame := func(slot int, frame *VideoFrame) bool {
		calls[slot]++
		if frame.Width != 2 || frame.Height != 3 || len(frame.YBuffer) != 12 || len(frame.UBuffer) != 4 || len(frame.VBuffer) != 4 || frame.AlphaMode != 2 {
			t.Errorf("incorrect frame conversion: %+v", frame)
		}
		return accept
	}
	queryCalls := 0
	observer := &VideoFrameObserver{
		OnCaptureVideoFrame: func(source int, frame *VideoFrame) bool {
			if source != 4 {
				t.Errorf("source = %d", source)
			}
			return checkFrame(0, frame)
		},
		OnPreEncodeVideoFrame: func(source int, frame *VideoFrame) bool {
			if source != 4 {
				t.Errorf("source = %d", source)
			}
			return checkFrame(1, frame)
		},
		OnMediaPlayerVideoFrame: func(frame *VideoFrame, player int) bool {
			if player != 42 {
				t.Errorf("player = %d", player)
			}
			return checkFrame(2, frame)
		},
		OnRenderVideoFrame: func(channel string, uid uint32, frame *VideoFrame) bool {
			if channel != "test-channel" || (uid != 0 && uid != ^uint32(0)) {
				t.Errorf("render arguments: %s, %d", channel, uid)
			}
			return checkFrame(3, frame)
		},
		OnTranscodedVideoFrame:   func(frame *VideoFrame) bool { return checkFrame(4, frame) },
		GetVideoFrameProcessMode: func() VideoFrameProcessMode { queryCalls++; return VideoFrameProcessModeReadOnly },
		GetVideoFormatPreference: func() VideoPixelFormat { return VideoPixelRGBA },
		GetRotationApplied:       func() bool { return true },
		GetMirrorApplied:         func() bool { return true },
		GetObservedFramePosition: func() VideoModulePosition { return VideoPositionPreEncoder },
	}
	conn := newTestVideoObserver(t, observer)
	if queryCalls != 1 {
		t.Fatalf("configuration was unavailable during set: %d calls", queryCalls)
	}
	for slot, want := range []int{0, 4, 1, 1, 4} {
		if got := testVideoObserverQuery(conn.cVideoObserver, slot); got != want {
			t.Errorf("query %d: %d, want %d", slot, got, want)
		}
	}
	for _, value := range []bool{false, true} {
		accept = value
		want := 0
		if accept {
			want = 1
		}
		for slot := 0; slot < 5; slot++ {
			if got, _ := testVideoObserverFrame(conn.cVideoObserver, slot, ^uint32(0), false, VideoPixelI420); got != want {
				t.Errorf("frame %d: %d, want %d", slot, got, want)
			}
		}
	}
	// UID zero must reach the new numeric render callback, too.
	testVideoObserverFrame(conn.cVideoObserver, 3, 0, false, VideoPixelI420)
	if calls != [5]int{2, 2, 2, 3, 2} {
		t.Errorf("callback counts: %v", calls)
	}
	handle := conn.cVideoObserver
	agoraService.deleteConFromHandle(handle, ConTypeCVideoObserver)
	defer agoraService.setConFromHandle(handle, conn, ConTypeCVideoObserver)
	if ret, _ := testVideoObserverFrame(handle, 3, 1, false, VideoPixelI420); ret != 0 {
		t.Errorf("unknown handle: %d", ret)
	}
}

func TestVideoFrameObserverFrameCopies(t *testing.T) {
	for _, mode := range []VideoFrameProcessMode{VideoFrameProcessModeReadOnly, VideoFrameProcessModeReadWrite} {
		for slot := 0; slot < 5; slot++ {
			t.Run(fmt.Sprintf("mode%d/slot%d", mode, slot), func(t *testing.T) {
				var retained *VideoFrame
				edit := func(frame *VideoFrame) bool {
					retained = frame
					frame.YBuffer[0] = 99
					frame.UBuffer[0] = 66
					frame.VBuffer[0] = 55
					frame.MetadataBuffer[0] = 88
					frame.AlphaBuffer[0] = 77
					frame.Rotation = VideoOrientation90
					return true
				}
				observer := &VideoFrameObserver{
					OnCaptureVideoFrame:      func(_ int, frame *VideoFrame) bool { return edit(frame) },
					OnPreEncodeVideoFrame:    func(_ int, frame *VideoFrame) bool { return edit(frame) },
					OnMediaPlayerVideoFrame:  func(frame *VideoFrame, _ int) bool { return edit(frame) },
					OnRenderVideoFrame:       func(_ string, _ uint32, frame *VideoFrame) bool { return edit(frame) },
					OnTranscodedVideoFrame:   edit,
					GetVideoFrameProcessMode: func() VideoFrameProcessMode { return mode },
				}
				conn := newTestVideoObserver(t, observer)
				ret, native := testVideoObserverFrame(conn.cVideoObserver, slot, 1, false, VideoPixelI420)
				if ret != 1 {
					t.Fatalf("callback: %d", ret)
				}
				if native.YBuffer[0] != 7 || native.UBuffer[0] != 7 || native.VBuffer[0] != 7 ||
					native.MetadataBuffer[0] != 7 || native.AlphaBuffer[0] != 7 || native.Rotation != VideoOrientation0 {
					t.Errorf("Go edits changed the native frame: %+v", native)
				}
				if retained.YBuffer[0] != 99 {
					t.Error("Go frame did not survive native frame release")
				}
			})
		}
	}
}

func TestVideoFrameObserverEditsAreLocal(t *testing.T) {
	for name, edit := range map[string]func(*VideoFrame) bool{
		"layout":         func(frame *VideoFrame) bool { frame.Height++; return true },
		"short-plane":    func(frame *VideoFrame) bool { frame.UBuffer = frame.UBuffer[:1]; return true },
		"large-metadata": func(frame *VideoFrame) bool { frame.MetadataBuffer = make([]byte, 20); return true },
		"rejected-frame": func(frame *VideoFrame) bool { return false },
	} {
		t.Run(name, func(t *testing.T) {
			conn := newTestVideoObserver(t, &VideoFrameObserver{
				GetVideoFrameProcessMode: func() VideoFrameProcessMode { return VideoFrameProcessModeReadWrite },
				OnTranscodedVideoFrame:   func(frame *VideoFrame) bool { frame.YBuffer[0] = 99; return edit(frame) },
			})
			want := 1
			if name == "rejected-frame" {
				want = 0
			}
			ret, native := testVideoObserverFrame(conn.cVideoObserver, 4, 1, false, VideoPixelI420)
			if ret != want || native.YBuffer[0] != 7 || native.Height != 3 || len(native.UBuffer) != 4 || len(native.MetadataBuffer) != 2 {
				t.Error("callback result or native frame changed after editing the Go copy")
			}
		})
	}
}

func TestVideoFrameObserverPackedFormat(t *testing.T) {
	conn := newTestVideoObserver(t, &VideoFrameObserver{
		OnTranscodedVideoFrame: func(frame *VideoFrame) bool {
			if len(frame.YBuffer) != 24 || len(frame.UBuffer) != 0 || len(frame.VBuffer) != 0 {
				t.Errorf("packed buffers: %d/%d/%d", len(frame.YBuffer), len(frame.UBuffer), len(frame.VBuffer))
			}
			return true
		},
	})
	testVideoObserverFrame(conn.cVideoObserver, 4, 1, false, VideoPixelRGBA)
}

func TestVideoFrameObserverPixelFormats(t *testing.T) {
	for _, test := range []struct {
		format  VideoPixelFormat
		y, u, v int
	}{
		{VideoPixelI420, 12, 4, 4},
		{VideoPixelI422, 12, 6, 6},
		{VideoPixelNV12, 12, 4, 0},
		{VideoPixelNV21, 12, 4, 0},
		{VideoPixelI010, 24, 8, 8},
		{VideoPixelBGRA, 24, 0, 0},
		{VideoPixelRGBA, 24, 0, 0},
		{VideoTexture2D, 0, 0, 0},
	} {
		t.Run(fmt.Sprintf("format%d", test.format), func(t *testing.T) {
			conn := newTestVideoObserver(t, &VideoFrameObserver{
				GetVideoFrameProcessMode: func() VideoFrameProcessMode { return VideoFrameProcessModeReadWrite },
				OnTranscodedVideoFrame: func(frame *VideoFrame) bool {
					if frame.Type != test.format || len(frame.YBuffer) != test.y || len(frame.UBuffer) != test.u || len(frame.VBuffer) != test.v {
						t.Errorf("format %d: plane lengths %d/%d/%d", frame.Type, len(frame.YBuffer), len(frame.UBuffer), len(frame.VBuffer))
					}
					if len(frame.YBuffer) > 0 {
						frame.YBuffer[0] = 33
					}
					return true
				},
			})
			ret, native := testVideoObserverFrame(conn.cVideoObserver, 4, 1, false, test.format)
			if ret != 1 || (len(native.YBuffer) > 0 && native.YBuffer[0] != 7) {
				t.Error("format callback changed native storage")
			}
		})
	}
}

func TestVideoFrameObserverProcessModeCallback(t *testing.T) {
	mode := VideoFrameProcessModeReadWrite
	calls := 0
	conn := newTestVideoObserver(t, &VideoFrameObserver{
		GetVideoFrameProcessMode: func() VideoFrameProcessMode { calls++; return mode },
	})
	if got := testVideoObserverQuery(conn.cVideoObserver, 0); got != int(VideoFrameProcessModeReadWrite) {
		t.Errorf("read/write mode = %d", got)
	}
	mode = VideoFrameProcessModeReadOnly
	if got := testVideoObserverQuery(conn.cVideoObserver, 0); got != int(VideoFrameProcessModeReadOnly) {
		t.Errorf("read-only mode = %d", got)
	}
	if calls != 3 {
		t.Errorf("mode callback invoked %d times, want 3", calls)
	}
}

func TestVideoFrameObserverLifecycle(t *testing.T) {
	observer := &VideoFrameObserver{}
	conn := newTestVideoObserver(t, observer)
	handle := conn.cVideoObserver
	if ret := conn.RegisterVideoFrameObserver(observer); ret != 0 || conn.cVideoObserver != handle || testVideoObserverCount(0) != 1 {
		t.Fatal("same observer was registered twice")
	}
	testVideoObserverFail(0, 0, -7)
	if ret := conn.unregisterVideoFrameObserver(); ret != -7 || conn.cVideoObserver != handle || testVideoObserverCount(1) != 0 {
		t.Fatal("failed unset destroyed a live observer")
	}
	if ret := conn.RegisterVideoFrameObserver(&VideoFrameObserver{}); ret != -7 || conn.videoObserver.Load() != observer {
		t.Fatal("replacement ignored unset failure")
	}
	testVideoObserverFail(0, -8, 0)
	if ret := conn.RegisterVideoFrameObserver(&VideoFrameObserver{}); ret != -8 || conn.cVideoObserver != nil || conn.videoObserver.Load() != nil || testVideoObserverCount(1) != 2 {
		t.Fatal("set failure leaked a handle or ignored native error")
	}
	if agoraService.getConFromHandle(handle, ConTypeCVideoObserver) != nil {
		t.Fatal("stale handle mapping")
	}
	testVideoObserverFail(1, 0, 0)
	if ret := conn.RegisterVideoFrameObserver(observer); ret == 0 || conn.videoObserver.Load() != nil {
		t.Fatal("creation failure was accepted")
	}
	testVideoObserverFail(0, 0, 0)
	if ret := conn.RegisterVideoFrameObserver(observer); ret != 0 {
		t.Fatalf("retry registration: %d", ret)
	}
	handle = conn.cVideoObserver
	if ret := conn.unregisterVideoFrameObserver(); ret != 0 {
		t.Fatalf("unregister: %d", ret)
	}
	count := testVideoObserverCount(1)
	if ret := conn.unregisterVideoFrameObserver(); ret != 0 || testVideoObserverCount(1) != count {
		t.Fatal("repeated unregister was not a no-op")
	}
	if agoraService.getConFromHandle(handle, ConTypeCVideoObserver) != nil {
		t.Fatal("unregister retained handle mapping")
	}
}

func TestVideoFrameObserverInvalidConnection(t *testing.T) {
	var conn *RtcConnection
	if conn.RegisterVideoFrameObserver(&VideoFrameObserver{}) == 0 {
		t.Error("nil connection registered")
	}
	if conn.unregisterVideoFrameObserver() != 0 {
		t.Error("nil connection unregister is not idempotent")
	}
	conn = &RtcConnection{}
	if conn.RegisterVideoFrameObserver(&VideoFrameObserver{}) == 0 {
		t.Error("released connection registered")
	}
	conn = testVideoObserverConnection()
	defer testVideoObserverFreeConnection(conn)
	if conn.RegisterVideoFrameObserver(nil) == 0 {
		t.Error("nil observer registered")
	}
}

func TestVideoFrameObserverRelease(t *testing.T) {
	for _, unsetResult := range []int{0, -7} {
		t.Run(fmt.Sprintf("unset%d", unsetResult), func(t *testing.T) {
			conn := newTestVideoObserver(t, &VideoFrameObserver{})
			handle := conn.cVideoObserver
			// A failed unset must bypass delayed connection destruction.
			oldIdleMode := agoraService.idleMode
			agoraService.idleMode = unsetResult != 0
			defer func() { agoraService.idleMode = oldIdleMode }()
			testVideoObserverFail(0, 0, unsetResult)
			conn.Release()
			if conn.cConnection != nil || conn.cVideoObserver != nil || conn.videoObserver.Load() != nil {
				t.Fatal("release retained native or Go observer state")
			}
			if testVideoObserverCount(3) != 1 || testVideoObserverCount(1) != 1 || testVideoObserverCount(4) != 1 || testVideoObserverCount(5) != 0 {
				t.Fatal("release freed an attached observer or failed to destroy native objects")
			}
			if agoraService.getConFromHandle(handle, ConTypeCVideoObserver) != nil {
				t.Fatal("release retained handle mapping")
			}
			conn.Release()
			if testVideoObserverCount(1) != 1 {
				t.Fatal("repeated release destroyed observer twice")
			}
		})
	}
}
