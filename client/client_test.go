package client

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/estenssoros/yeetbot/slack"
)

var (
	testIndex = "yeetbot-test"
)

func MockNewClient(t *testing.T) *Client {
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	client := config.NewClient(config.Reports[0])
	client.ElasticIndex = testIndex
	return client
}

func MockNewUser() *slack.User {
	return &slack.User{ID: "berto"}
}

func TestReadConfig(t *testing.T) {
	path := os.Getenv("YEETBOT_CONFIG")
	if path == "" {
		t.Fatal("missing YEETBOT_CONFIG env")
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, err = ConfigFromReader(f)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadConfigEnv(t *testing.T) {
	if _, err := ConfigFromEnv(); err != nil {
		t.Fatal(err)
	}
}

func TestReportTodayTime(t *testing.T) {
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	r := config.Reports[0]
	fmt.Println(r.TodayTime())
}

func TestClientListUsers(t *testing.T) {
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	client := config.NewClient(config.Reports[0])
	users, err := client.ListUsers()
	if err != nil {
		t.Fatal(err)
	}
	if len(users) == 0 {
		t.Fatal("no users returned")
	}
}

func TestClientGetUser(t *testing.T) {
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	client := config.NewClient(config.Reports[0])
	userByName, err := client.GetUserByName("sebastian")
	if err != nil {
		t.Fatal(err)
	}
	userByID, err := client.GetUserByID(userByName.ID)
	if err != nil {
		t.Fatal(err)
	}
	if userByName.Name != userByID.Name {
		t.Fatal("user names don't match")
	}
	if userByID.ID != userByName.ID {
		t.Fatal("user ids don't match")
	}
}

func TestClientGetuserReportTime(t *testing.T) {
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	client := config.NewClient(config.Reports[0])
	user, err := client.GetUserByName("sebastian")
	if err != nil {
		t.Fatal(err)
	}
	userTime, err := client.UserReportTime(user)
	if err != nil {
		t.Fatal(err)
	}
	todayTime, err := client.TodayTime()
	if err != nil {
		t.Fatal(err)
	}
	loc, err := time.LoadLocation("America/Denver")
	if err != nil {
		t.Fatal(err)
	}
	todayTime = todayTime.In(loc)
	if want, have := todayTime.Format("15:04"), userTime.Format("15:04"); want != have {
		t.Fatalf("have: %v, want: %v", have, want)
	}
}
