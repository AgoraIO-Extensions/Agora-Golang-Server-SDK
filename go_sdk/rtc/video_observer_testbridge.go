//go:build agora_video_observer_test

package agoraservice

/*
#include <stdlib.h>
#include "agora_local_user.h"

void test_video_observer_reset(void);
void test_video_observer_fail(int create_fail, int set_result, int unset_result);
int test_video_observer_count(int operation);
int test_video_observer_query(void* observer, int slot);
video_frame* test_video_frame_create(int format);
int test_video_observer_frame(void* observer, int slot, unsigned int uid, video_frame* frame);
*/
import "C"
import "unsafe"

func testVideoObserverReset() { C.test_video_observer_reset() }
func testVideoObserverFail(createFail, setResult, unsetResult int) {
	C.test_video_observer_fail(C.int(createFail), C.int(setResult), C.int(unsetResult))
}
func testVideoObserverCount(operation int) int {
	return int(C.test_video_observer_count(C.int(operation)))
}
func testVideoObserverQuery(observer unsafe.Pointer, slot int) int {
	return int(C.test_video_observer_query(observer, C.int(slot)))
}
func testVideoObserverFrame(observer unsafe.Pointer, slot int, uid uint32, nilFrame bool, format VideoPixelFormat) (int, *VideoFrame) {
	var frame *C.video_frame
	if !nilFrame {
		frame = C.test_video_frame_create(C.int(format))
		defer C.free(unsafe.Pointer(frame))
	}
	ret := int(C.test_video_observer_frame(observer, C.int(slot), C.uint(uid), frame))
	if frame == nil {
		return ret, nil
	}
	return ret, GoVideoFrame(frame)
}
func testVideoObserverConnection() *RtcConnection {
	handle := C.malloc(1)
	return &RtcConnection{
		cConnection: handle,
		localUser:   &LocalUser{cLocalUser: handle},
	}
}
func testVideoObserverFreeConnection(conn *RtcConnection) {
	conn.unregisterVideoFrameObserver()
	C.free(conn.cConnection)
}
