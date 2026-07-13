/**
 * author: Wei Hongqi
 * date: 2026-07-13
 * how to run:
 cd go_sdk/rtc
CGO_ENABLED=1 go test -vet=off -v -run 'RemoteAudioTrackAPM|InheritServiceAPM|SetRemoteAudioTrackAPMModel|IsEnable3A' .

// whole file:
cd go_sdk/rtc
CGO_ENABLED=1 go test -vet=off -v -run 'RemoteAudioTrackAPM|InheritServiceAPM|SetRemoteAudioTrackAPMModel|IsEnable3A' .
*/

package agoraservice

import (
	"testing"
	"unsafe"
)

func distinctAPMConfig() *APMConfig {
	cfg := NewAPMConfig()
	cfg.EnableDump = true
	return cfg
}

func newTestRtcConnectionForAPM() *RtcConnection {
	dummy := unsafe.Pointer(uintptr(1))
	return &RtcConnection{
		cConnection: dummy,
		localUser:   &LocalUser{cLocalUser: dummy},
	}
}

// inheritServiceAPMToConnection mirrors NewRtcConnection APM initialization (rtc_connection.go L535-536).
func inheritServiceAPMToConnection(conn *RtcConnection) {
	conn.remoteAudioTrackAPMModel = agoraService.apmModel
	conn.remoteAudioTrackAPMConfig = agoraService.apmConfig
}

func withServiceAPMState(t *testing.T, model int, config *APMConfig, fn func()) {
	t.Helper()
	oldModel := agoraService.apmModel
	oldConfig := agoraService.apmConfig
	agoraService.apmModel = model
	agoraService.apmConfig = config
	t.Cleanup(func() {
		agoraService.apmModel = oldModel
		agoraService.apmConfig = oldConfig
	})
	fn()
}

func TestInheritServiceAPM_WithoutSetter(t *testing.T) {
	svcCfg := distinctAPMConfig()

	cases := []struct {
		name       string
		svcModel   int
		svcConfig  *APMConfig
		wantModel  int
		wantConfig *APMConfig
		wantEnable bool
	}{
		{
			name:       "service_off",
			svcModel:   ApmModeOff,
			svcConfig:  nil,
			wantModel:  ApmModeOff,
			wantConfig: nil,
			wantEnable: false,
		},
		{
			name:       "service_on_with_config",
			svcModel:   ApmModeOn,
			svcConfig:  svcCfg,
			wantModel:  ApmModeOn,
			wantConfig: svcCfg,
			wantEnable: true,
		},
		{
			name:       "service_on_nil_config",
			svcModel:   ApmModeOn,
			svcConfig:  nil,
			wantModel:  ApmModeOn,
			wantConfig: nil,
			wantEnable: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withServiceAPMState(t, tc.svcModel, tc.svcConfig, func() {
				conn := newTestRtcConnectionForAPM()
				inheritServiceAPMToConnection(conn)

				if conn.remoteAudioTrackAPMModel != tc.wantModel {
					t.Fatalf("remoteAudioTrackAPMModel=%d, want %d", conn.remoteAudioTrackAPMModel, tc.wantModel)
				}
				if conn.remoteAudioTrackAPMConfig != tc.wantConfig {
					t.Fatalf("remoteAudioTrackAPMConfig=%p, want %p", conn.remoteAudioTrackAPMConfig, tc.wantConfig)
				}
				if got := conn.isEnable3A(); got != tc.wantEnable {
					t.Fatalf("isEnable3A()=%v, want %v", got, tc.wantEnable)
				}
			})
		})
	}
}

func TestSetRemoteAudioTrackAPMModel_OverridesService(t *testing.T) {
	svcCfg := distinctAPMConfig()
	customCfg := NewAPMConfig()
	customCfg.EnableDump = false

	t.Run("override_off_when_service_on", func(t *testing.T) {
		withServiceAPMState(t, ApmModeOn, svcCfg, func() {
			conn := newTestRtcConnectionForAPM()
			inheritServiceAPMToConnection(conn)

			if ret := conn.SetRemoteAudioTrackAPMModel(ApmModeOff, nil); ret != 0 {
				t.Fatalf("SetRemoteAudioTrackAPMModel()=%d, want 0", ret)
			}
			if conn.remoteAudioTrackAPMModel != ApmModeOff {
				t.Fatalf("remoteAudioTrackAPMModel=%d, want %d", conn.remoteAudioTrackAPMModel, ApmModeOff)
			}
			if conn.remoteAudioTrackAPMConfig != nil {
				t.Fatalf("remoteAudioTrackAPMConfig=%p, want nil", conn.remoteAudioTrackAPMConfig)
			}
			if conn.isEnable3A() {
				t.Fatal("isEnable3A()=true, want false")
			}
		})
	})

	t.Run("override_on_when_service_off", func(t *testing.T) {
		withServiceAPMState(t, ApmModeOff, nil, func() {
			conn := newTestRtcConnectionForAPM()
			inheritServiceAPMToConnection(conn)

			if ret := conn.SetRemoteAudioTrackAPMModel(ApmModeOn, customCfg); ret != 0 {
				t.Fatalf("SetRemoteAudioTrackAPMModel()=%d, want 0", ret)
			}
			if conn.remoteAudioTrackAPMModel != ApmModeOn {
				t.Fatalf("remoteAudioTrackAPMModel=%d, want %d", conn.remoteAudioTrackAPMModel, ApmModeOn)
			}
			if conn.remoteAudioTrackAPMConfig != customCfg {
				t.Fatalf("remoteAudioTrackAPMConfig=%p, want %p", conn.remoteAudioTrackAPMConfig, customCfg)
			}
			if !conn.isEnable3A() {
				t.Fatal("isEnable3A()=false, want true")
			}
		})
	})

	t.Run("override_on_uses_custom_config_not_service", func(t *testing.T) {
		withServiceAPMState(t, ApmModeOn, svcCfg, func() {
			conn := newTestRtcConnectionForAPM()
			inheritServiceAPMToConnection(conn)

			if ret := conn.SetRemoteAudioTrackAPMModel(ApmModeOn, customCfg); ret != 0 {
				t.Fatalf("SetRemoteAudioTrackAPMModel()=%d, want 0", ret)
			}
			if conn.remoteAudioTrackAPMConfig != customCfg {
				t.Fatalf("remoteAudioTrackAPMConfig=%p, want %p", conn.remoteAudioTrackAPMConfig, customCfg)
			}
			if conn.remoteAudioTrackAPMConfig == svcCfg {
				t.Fatal("remoteAudioTrackAPMConfig should not point to service config after override")
			}
		})
	})

	t.Run("service_change_after_override_does_not_affect_connection", func(t *testing.T) {
		withServiceAPMState(t, ApmModeOff, nil, func() {
			conn := newTestRtcConnectionForAPM()
			inheritServiceAPMToConnection(conn)

			if ret := conn.SetRemoteAudioTrackAPMModel(ApmModeOn, customCfg); ret != 0 {
				t.Fatalf("SetRemoteAudioTrackAPMModel()=%d, want 0", ret)
			}

			agoraService.apmModel = ApmModeOff
			agoraService.apmConfig = nil

			if conn.remoteAudioTrackAPMModel != ApmModeOn {
				t.Fatalf("remoteAudioTrackAPMModel=%d, want %d after service change", conn.remoteAudioTrackAPMModel, ApmModeOn)
			}
			if conn.remoteAudioTrackAPMConfig != customCfg {
				t.Fatalf("remoteAudioTrackAPMConfig=%p, want %p after service change", conn.remoteAudioTrackAPMConfig, customCfg)
			}
			if !conn.isEnable3A() {
				t.Fatal("isEnable3A()=false, want true after service change")
			}
		})
	})
}

func TestSetRemoteAudioTrackAPMModel_Validation(t *testing.T) {
	cfg := distinctAPMConfig()

	t.Run("nil_connection", func(t *testing.T) {
		if ret := (*RtcConnection)(nil).SetRemoteAudioTrackAPMModel(ApmModeOn, cfg); ret != -2000 {
			t.Fatalf("SetRemoteAudioTrackAPMModel()=%d, want -2000", ret)
		}
	})

	t.Run("on_without_config", func(t *testing.T) {
		conn := newTestRtcConnectionForAPM()
		if ret := conn.SetRemoteAudioTrackAPMModel(ApmModeOn, nil); ret != -2002 {
			t.Fatalf("SetRemoteAudioTrackAPMModel()=%d, want -2002", ret)
		}
		if conn.remoteAudioTrackAPMModel != 0 {
			t.Fatalf("remoteAudioTrackAPMModel=%d, want unchanged 0 after -2002", conn.remoteAudioTrackAPMModel)
		}
		if conn.remoteAudioTrackAPMConfig != nil {
			t.Fatalf("remoteAudioTrackAPMConfig=%p, want unchanged nil after -2002", conn.remoteAudioTrackAPMConfig)
		}
	})

	t.Run("off_clears_config", func(t *testing.T) {
		conn := newTestRtcConnectionForAPM()
		conn.remoteAudioTrackAPMModel = ApmModeOn
		conn.remoteAudioTrackAPMConfig = cfg

		if ret := conn.SetRemoteAudioTrackAPMModel(ApmModeOff, cfg); ret != 0 {
			t.Fatalf("SetRemoteAudioTrackAPMModel()=%d, want 0", ret)
		}
		if conn.remoteAudioTrackAPMModel != ApmModeOff {
			t.Fatalf("remoteAudioTrackAPMModel=%d, want %d", conn.remoteAudioTrackAPMModel, ApmModeOff)
		}
		if conn.remoteAudioTrackAPMConfig != nil {
			t.Fatalf("remoteAudioTrackAPMConfig=%p, want nil", conn.remoteAudioTrackAPMConfig)
		}
	})

	t.Run("on_with_config", func(t *testing.T) {
		conn := newTestRtcConnectionForAPM()
		if ret := conn.SetRemoteAudioTrackAPMModel(ApmModeOn, cfg); ret != 0 {
			t.Fatalf("SetRemoteAudioTrackAPMModel()=%d, want 0", ret)
		}
		if conn.remoteAudioTrackAPMModel != ApmModeOn {
			t.Fatalf("remoteAudioTrackAPMModel=%d, want %d", conn.remoteAudioTrackAPMModel, ApmModeOn)
		}
		if conn.remoteAudioTrackAPMConfig != cfg {
			t.Fatalf("remoteAudioTrackAPMConfig=%p, want %p", conn.remoteAudioTrackAPMConfig, cfg)
		}
	})
}

func TestIsEnable3A_Boundary(t *testing.T) {
	cfg := distinctAPMConfig()

	cases := []struct {
		name       string
		model      int
		config     *APMConfig
		wantEnable bool
	}{
		{"off_nil_config", ApmModeOff, nil, false},
		{"off_with_config", ApmModeOff, cfg, false},
		{"on_nil_config", ApmModeOn, nil, false},
		{"on_with_config", ApmModeOn, cfg, true},
		{"inherit_with_config", ApmModeInherit, cfg, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conn := newTestRtcConnectionForAPM()
			conn.remoteAudioTrackAPMModel = tc.model
			conn.remoteAudioTrackAPMConfig = tc.config

			if got := conn.isEnable3A(); got != tc.wantEnable {
				t.Fatalf("isEnable3A()=%v, want %v", got, tc.wantEnable)
			}
		})
	}
}

func TestSetRemoteAudioTrackAPMModel_ApmModeInherit(t *testing.T) {
	svcCfg := distinctAPMConfig()

	withServiceAPMState(t, ApmModeOn, svcCfg, func() {
		conn := newTestRtcConnectionForAPM()
		inheritServiceAPMToConnection(conn)

		if ret := conn.SetRemoteAudioTrackAPMModel(ApmModeInherit, nil); ret != 0 {
			t.Fatalf("SetRemoteAudioTrackAPMModel()=%d, want 0", ret)
		}
		if conn.remoteAudioTrackAPMModel != ApmModeInherit {
			t.Fatalf("remoteAudioTrackAPMModel=%d, want %d", conn.remoteAudioTrackAPMModel, ApmModeInherit)
		}
		if conn.remoteAudioTrackAPMConfig != nil {
			t.Fatalf("remoteAudioTrackAPMConfig=%p, want nil", conn.remoteAudioTrackAPMConfig)
		}
		if conn.isEnable3A() {
			t.Fatal("isEnable3A()=true, want false for ApmModeInherit")
		}
	})
}
