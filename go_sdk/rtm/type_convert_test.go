package agorartm

import "testing"

func TestCUserListToUserList_Nil(t *testing.T) {
	if got := CUserListToUserList(nil); got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
}

func TestCUserListToUserList_Empty(t *testing.T) {
	cList := newCTestUserList(t)
	if got := CUserListToUserList(&cList); got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
}

func TestCUserListToUserList_MultipleUsers(t *testing.T) {
	cList := newCTestUserList(t, "alice", "bob")
	got := CUserListToUserList(&cList)
	if got == nil {
		t.Fatal("expected non-nil UserList")
	}
	if len(got.Users) != 2 {
		t.Fatalf("len(Users)=%d, want 2", len(got.Users))
	}
	if got.Users[0] != "alice" || got.Users[1] != "bob" {
		t.Fatalf("Users=%v", got.Users)
	}
}

func TestCUserStateToUserState_Nil(t *testing.T) {
	if got := CUserStateToUserState(nil); got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
}

func TestCUserStateToUserState_NoStates(t *testing.T) {
	cList, _ := newCTestUserStateList(t, userStateSpec{userID: "user-1"})
	got := CUserStateToUserState(cList)
	if got == nil {
		t.Fatal("expected non-nil UserState")
	}
	if got.UserId != "user-1" {
		t.Fatalf("UserId=%q", got.UserId)
	}
	if len(got.States) != 0 {
		t.Fatalf("States=%v, want empty", got.States)
	}
}

func TestCUserStateToUserState_WithStates(t *testing.T) {
	cList, _ := newCTestUserStateList(t, userStateSpec{
		userID: "user-2",
		states: []StateItem{{Key: "role", Value: "host"}},
	})
	got := CUserStateToUserState(cList)
	if got == nil {
		t.Fatal("expected non-nil UserState")
	}
	if got.UserId != "user-2" {
		t.Fatalf("UserId=%q", got.UserId)
	}
	if len(got.States) != 1 {
		t.Fatalf("len(States)=%d", len(got.States))
	}
	if got.States[0].Key != "role" || got.States[0].Value != "host" {
		t.Fatalf("States=%v", got.States)
	}
}

func TestCUserStateListToUserStateList_Empty(t *testing.T) {
	if got := CUserStateListToUserStateList(nil, 0); got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
}

func TestCUserStateListToUserStateList_MultipleUsers(t *testing.T) {
	cList, count := newCTestUserStateList(t,
		userStateSpec{userID: "u1"},
		userStateSpec{userID: "u2", states: []StateItem{{Key: "k", Value: "v"}}},
	)

	got := CUserStateListToUserStateList(cList, count)
	if got == nil {
		t.Fatal("expected non-nil slice")
	}
	if len(got) != 2 {
		t.Fatalf("len=%d, want 2", len(got))
	}
	if got[0] == nil || got[0].UserId != "u1" {
		t.Fatalf("got[0]=%#v", got[0])
	}
	if got[1] == nil || got[1].UserId != "u2" {
		t.Fatalf("got[1]=%#v", got[1])
	}
	if len(got[1].States) != 1 || got[1].States[0].Key != "k" {
		t.Fatalf("got[1].States=%v", got[1].States)
	}
}
