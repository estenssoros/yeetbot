package client

import (
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
)

func TestClientGetUserResponses(t *testing.T) {
	// TODO: fix
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	client := config.NewClient(config.Reports[0])
	client.ElasticIndex = testIndex
	mockUser := "berto"
	user, err := client.GetUserByName(mockUser)
	if err != nil {
		t.Fatal(err)
	}
	mockUserReports := map[string][]*Report{mockUser: []*Report{config.Reports[0]}}
	report, err := client.FindReportByUser(user, mockUserReports)
	if err != nil {
		t.Fatal(err)
	}
	newResponse := &Response{
		ID:       uuid.Must(uuid.NewV4()),
		Report:   report.Name,
		Channel:  report.Channel,
		UserID:   user.ID,
		EventTS:  "timestamp",
		Date:     time.Now(),
		Question: report.Questions[0].Text,
		Text:     "test",
	}
	// err = deleteIndex(client.ElasticURL)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// err = client.RecordResponse(newResponse)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	client.Report = report
	response, err := client.GetUserResponse(user)
	if err != nil {
		t.Fatal(err)
	}
	if have, want := response, newResponse; want != have {
		t.Fatalf("have: %v, want: %v", have, want)
	}
}

func TestClientRecordAndCompleteResponses(t *testing.T) {
	// TODO: RecordResponse
	// TODO: CompleteResponse
}

// func deleteIndex(url string) error {
// 	svc := elasticsvc.New(context.Background())
// 	svc.SetURL(url)
// 	if err := svc.DeleteIndex(testIndex); err != nil {
// 		return err
// 	}
// 	if err := svc.CreateIndex(testIndex, testMapping); err != nil {
// 		return err
// 	}
// 	return nil
// }
