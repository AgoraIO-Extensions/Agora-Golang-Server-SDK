//go:build test

package agorartm

import "testing"

func TestRtmConfigOptions(t *testing.T) {
	reconnectTimeout, parameters := readCTestRtmConfigOptions(t, &RtmConfig{
		ReconnectTimeout: 120,
		Parameters:       `{"region":"cn"}`,
	})

	if reconnectTimeout != 120 {
		t.Fatalf("reconnectTimeout=%d", reconnectTimeout)
	}
	if parameters != `{"region":"cn"}` {
		t.Fatalf("parameters=%q", parameters)
	}
}
