package agorartm

import "testing"

func TestOnWhoNowResult_UserStateList(t *testing.T) {
	var gotRequestID uint64
	var gotUsers []*UserState
	var gotNextPage string
	var gotErrorCode int

	client := &IRtmClient{
		handler: &RtmEventHandler{
			OnWhoNowResult: func(requestId uint64, userStateList []*UserState, nextPage string, errorCode int) {
				gotRequestID = requestId
				gotUsers = userStateList
				gotNextPage = nextPage
				gotErrorCode = errorCode
			},
		},
	}

	handler := newCTestEventHandler(t, client)
	cList, count := newCTestUserStateList(t,
		userStateSpec{userID: "who-1"},
		userStateSpec{userID: "who-2", states: []StateItem{{Key: "role", Value: "member"}}},
	)
	nextPage := newCTestCString(t, "page-2")

	cgo_RtmEventHandlerBridge_onWhoNowResult(handler, 42, cList, count, nextPage, 0)

	if gotRequestID != 42 {
		t.Fatalf("requestId=%d", gotRequestID)
	}
	if len(gotUsers) != 2 {
		t.Fatalf("len(userStateList)=%d", len(gotUsers))
	}
	if gotUsers[0] == nil || gotUsers[0].UserId != "who-1" {
		t.Fatalf("userStateList[0]=%#v", gotUsers[0])
	}
	if gotUsers[1] == nil || gotUsers[1].UserId != "who-2" {
		t.Fatalf("userStateList[1]=%#v", gotUsers[1])
	}
	if len(gotUsers[1].States) != 1 || gotUsers[1].States[0].Key != "role" {
		t.Fatalf("userStateList[1].States=%v", gotUsers[1].States)
	}
	if gotNextPage != "page-2" {
		t.Fatalf("nextPage=%q", gotNextPage)
	}
	if gotErrorCode != 0 {
		t.Fatalf("errorCode=%d", gotErrorCode)
	}
}

func TestOnWhoNowResult_EmptyList(t *testing.T) {
	var gotUsers []*UserState

	client := &IRtmClient{
		handler: &RtmEventHandler{
			OnWhoNowResult: func(_ uint64, userStateList []*UserState, _ string, _ int) {
				gotUsers = userStateList
			},
		},
	}

	handler := newCTestEventHandler(t, client)
	cgo_RtmEventHandlerBridge_onWhoNowResult(handler, 0, nil, 0, nil, 0)

	if gotUsers == nil {
		t.Fatal("expected non-nil empty slice")
	}
	if len(gotUsers) != 0 {
		t.Fatalf("len(userStateList)=%d, want 0", len(gotUsers))
	}
}
