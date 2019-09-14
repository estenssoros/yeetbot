package client

import (
	"context"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/seaspancode/services/elasticsvc"
)

var (
	testIndex   = "yeetbot-test"
	testMapping = `
	{
			"settings": {
					"number_of_shards": 1,
					"number_of_replicas": 1
			},
			"mappings": {
					"Response": {
							"properties": {
									"user":     					{"type": "text"},
									"report":  						{"type": "text"},
									"date":      					{"type": "date", "format": "dateOptionalTime"},
									"pending_response":   {"type": "boolean"},
									"responses":  				{"type": "text"}
							}
					}
			}
	}`
)

func TestClientGetUserResponses(t *testing.T) {
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
	report, err := FindReportByUser(user, mockUserReports)
	if err != nil {
		t.Fatal(err)
	}
	newResponse := &Response{
		ID:              uuid.NewV4(),
		User:            user.RealName,
		Report:          report.Name,
		Date:            time.Now(),
		Responses:       []string{},
		PendingResponse: true,
	}
	// err = deleteIndex(client.ElasticURL)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	err = client.AddUserResponse(newResponse)
	if err != nil {
		t.Fatal(err)
	}
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

func deleteIndex(url string) error {
	svc := elasticsvc.New(context.Background())
	svc.SetURL(url)
	if err := svc.DeleteIndex(testIndex); err != nil {
		return err
	}
	if err := svc.CreateIndex(testIndex, testMapping); err != nil {
		return err
	}
	return nil
}
