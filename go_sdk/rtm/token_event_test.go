//go:build test

package agorartm

import "testing"

func TestTokenEventFromC(t *testing.T) {
	ev := newCTestTokenEvent(t)

	goEvent := NewTokenEvent()
	goEvent.fromC(ev)

	if goEvent.EventType != RtmTokenEventType(2) {
		t.Fatalf("EventType=%v", goEvent.EventType)
	}
	if goEvent.Reason != "permission revoked" {
		t.Fatalf("Reason=%q", goEvent.Reason)
	}
	if goEvent.Timestamp != 123456789 {
		t.Fatalf("Timestamp=%d", goEvent.Timestamp)
	}
	if got := goEvent.AffectedResources.MessageChannels; len(got) != 2 || got[0] != "alpha" || got[1] != "beta" {
		t.Fatalf("MessageChannels=%v", got)
	}
}

func TestTokenEventFromC_Empty(t *testing.T) {
	ev := newCTestEmptyTokenEvent(t)

	goEvent := NewTokenEvent()
	goEvent.fromC(ev)

	if goEvent.Reason != "" {
		t.Fatalf("Reason=%q", goEvent.Reason)
	}
	if goEvent.AffectedResources.MessageChannels == nil {
		t.Fatal("MessageChannels is nil")
	}
	if len(goEvent.AffectedResources.MessageChannels) != 0 {
		t.Fatalf("MessageChannels=%v", goEvent.AffectedResources.MessageChannels)
	}
}
