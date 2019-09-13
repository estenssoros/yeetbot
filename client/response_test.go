package client

import "testing"

func TestClientGetUserResponses(t *testing.T) {
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	client := config.NewClient(config.Reports[0])
	mockUser := "berto"
	user, err := client.GetUserByName(mockUser)
	if err != nil {
		t.Fatal(err)
	}
	mockUserReports := map[string][]*Report{mockUser: []*Report{config.Reports[0]}}
	report, err := FindReportByUser(user, mockUserReports)
	if err != nil {
		t.Fatal(err)
	}
	client.Report = report
	response, err := client.GetUserResponse(user)
	if err != nil {
		t.Fatal(err)
	}
	var expected *Response
	if want, have := response, expected; want != have {
		t.Fatalf("have: %v, want: %v", have, want)
	}
}

func TestClientRecordAndCompleteResponses(t *testing.T) {
	// TODO: RecordResponse
	// TODO: CompleteResponse
}
