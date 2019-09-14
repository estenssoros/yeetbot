package client

import (
	"testing"
)

func TestClientUserHasStartedReport(t *testing.T) {
	mockUser := MockNewUser()
	client := MockNewClient(t)
	got, err := client.HasUserStartedReport(mockUser)
	if err != nil {
		t.Fatal(err)
	}
	if have, want := got, false; want != have {
		t.Fatalf("have: %v, want: %v", have, want)
	}
	newResponseInput := &RecordResponseInput{
		Question: &Question{
			Text:  "test",
			Color: "pink",
		},
		User: mockUser,
		Text: "test",
	}
	client.RecordResponse(newResponseInput)
	got, err = client.HasUserStartedReport(mockUser)
	if err != nil {
		t.Fatal(err)
	}
	if have, want := got, true; want != have {
		t.Fatalf("have: %v, want: %v", have, want)
	}
}
