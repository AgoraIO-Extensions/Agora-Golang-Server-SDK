package agorartm

import "testing"

func TestPresenceEventFromC_IntervalUserLists(t *testing.T) {
	ev := newCTestPresenceEvent(t)
	ev.channelName = newCTestCString(t, "test-channel")
	ev.publisher = newCTestCString(t, "pub-1")
	ev.interval.joinUserList = newCTestUserList(t, "join-a", "join-b")
	ev.interval.leaveUserList = newCTestUserList(t, "leave-a")
	ev.interval.timeoutUserList = newCTestUserList(t, "timeout-a")

	goEvent := NewPresenceEvent()
	goEvent.fromC(ev)

	if goEvent.ChannelName != "test-channel" {
		t.Fatalf("ChannelName=%q", goEvent.ChannelName)
	}
	if goEvent.Interval == nil {
		t.Fatal("Interval is nil")
	}
	if len(goEvent.Interval.JoinUserList.Users) != 2 {
		t.Fatalf("JoinUserList=%v", goEvent.Interval.JoinUserList.Users)
	}
	if goEvent.Interval.JoinUserList.Users[0] != "join-a" {
		t.Fatalf("JoinUserList=%v", goEvent.Interval.JoinUserList.Users)
	}
	if len(goEvent.Interval.LeaveUserList.Users) != 1 || goEvent.Interval.LeaveUserList.Users[0] != "leave-a" {
		t.Fatalf("LeaveUserList=%v", goEvent.Interval.LeaveUserList.Users)
	}
	if len(goEvent.Interval.TimeoutUserList.Users) != 1 || goEvent.Interval.TimeoutUserList.Users[0] != "timeout-a" {
		t.Fatalf("TimeoutUserList=%v", goEvent.Interval.TimeoutUserList.Users)
	}
}

func TestPresenceEventFromC_IntervalUserStateList(t *testing.T) {
	ev := newCTestPresenceEvent(t)
	userStateList, count := newCTestUserStateList(t,
		userStateSpec{userID: "u1"},
		userStateSpec{userID: "u2", states: []StateItem{{Key: "status", Value: "online"}}},
	)
	ev.interval.userStateList = userStateList
	ev.interval.userStateCount = count

	goEvent := NewPresenceEvent()
	goEvent.fromC(ev)

	if goEvent.Interval == nil {
		t.Fatal("Interval is nil")
	}
	if len(goEvent.Interval.UserStateList) != 2 {
		t.Fatalf("UserStateList len=%d", len(goEvent.Interval.UserStateList))
	}
	if goEvent.Interval.UserStateCount != 2 {
		t.Fatalf("UserStateCount=%d", goEvent.Interval.UserStateCount)
	}
	if goEvent.Interval.UserStateList[1].UserId != "u2" {
		t.Fatalf("UserStateList[1]=%#v", goEvent.Interval.UserStateList[1])
	}
	if len(goEvent.Interval.UserStateList[1].States) != 1 {
		t.Fatalf("States=%v", goEvent.Interval.UserStateList[1].States)
	}
}

func TestPresenceEventFromC_SnapshotUserStateList(t *testing.T) {
	ev := newCTestPresenceEvent(t)
	userStateList, count := newCTestUserStateList(t, userStateSpec{userID: "snap-user"})
	ev.snapshot.userStateList = userStateList
	ev.snapshot.userCount = count

	goEvent := NewPresenceEvent()
	goEvent.fromC(ev)

	if goEvent.Snapshot == nil {
		t.Fatal("Snapshot is nil")
	}
	if len(goEvent.Snapshot.UserStateList) != 1 {
		t.Fatalf("UserStateList len=%d", len(goEvent.Snapshot.UserStateList))
	}
	if goEvent.Snapshot.UserCount != 1 {
		t.Fatalf("UserCount=%d", goEvent.Snapshot.UserCount)
	}
	if goEvent.Snapshot.UserStateList[0].UserId != "snap-user" {
		t.Fatalf("UserStateList[0]=%#v", goEvent.Snapshot.UserStateList[0])
	}
}

func TestPresenceEventFromC_StateItems(t *testing.T) {
	ev := newCTestPresenceEvent(t)
	stateItems, count := newCTestStateItems(t, StateItem{Key: "k1", Value: "v1"})
	ev.stateItems = stateItems
	ev.stateItemCount = count

	goEvent := NewPresenceEvent()
	goEvent.fromC(ev)

	if len(goEvent.StateItems) != 1 {
		t.Fatalf("StateItems len=%d", len(goEvent.StateItems))
	}
	if goEvent.StateItems[0].Key != "k1" || goEvent.StateItems[0].Value != "v1" {
		t.Fatalf("StateItems[0]=%#v", goEvent.StateItems[0])
	}
}
