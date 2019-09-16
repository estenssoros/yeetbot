package client

import (
	"os"
	"testing"

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
	client := config.NewClient()
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

func TestClientListUsers(t *testing.T) {
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	client := config.NewClient()
	users, err := client.ListUsers()
	if err != nil {
		t.Fatal(err)
	}
	if len(users) == 0 {
		t.Fatal("no users returned")
	}
	client.BotToken = ""
	if _, err := client.ListUsers(); err == nil {
		t.Fatal("should error")
	}
}

func TestClientGetUser(t *testing.T) {
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	client := config.NewClient()
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

func TestClientStriog(t *testing.T) {
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	client := config.NewClient()
	clientStr := client.String()
	if clientStr == "" {
		t.Fatal("error stringifying client")
	}
}

func TestClientSendMessage(t *testing.T) {
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	client := config.NewClient()
	{
		msg := &slack.Message{
			Channel: "#test",
			Text:    "test",
		}
		if err := client.SendMessage(msg); err != nil {
			t.Fatal(err)
		}
	}
	{
		msg := &slack.Message{
			Channel: "#test",
		}
		if err := client.SendMessage(msg); err == nil {
			t.Fatal("did not error")
		}
	}
}

func TestClientPostRequest(t *testing.T) {
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	client := config.NewClient()
	if err := client.postRequest("asdf", nil); err == nil {
		t.Fatal("should have returned error")
	}
	if err := client.postRequest("http://google.com", "asdf"); err == nil {
		t.Fatal("should be error")
	}
}

func TestClientGenericMessage(t *testing.T) {
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	client := config.NewClient()
	if err := client.GenericMessage(&slack.User{Name: "sebastian"}, "asdf"); err != nil {
		t.Fatal(err)
	}
}

func TestClientSendGreeting(t *testing.T) {
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	client := config.NewClient()
	meeting := &Meeting{
		IntroMessage: "{{.Name}} this is a test",
	}
	user := &slack.User{Name: "sebastian"}
	if err := client.SendGreeting(meeting, user); err != nil {
		t.Fatal(err)
	}
	meeting = &Meeting{IntroMessage: "{{.asdf}}"}
	if err := client.SendGreeting(meeting, user); err == nil {
		t.Fatal("should error")
	}
	client.UserToken = ""
	meeting = &Meeting{IntroMessage: ""}
	if err := client.SendGreeting(meeting, user); err == nil {
		t.Fatal("should error")
	}
	client.BotToken = ""
	meeting = &Meeting{IntroMessage: ""}
	if err := client.SendGreeting(meeting, user); err == nil {
		t.Fatal("should error")
	}
}

func TestClientListChannels(t *testing.T) {
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	client := config.NewClient()
	channels, err := client.ListChannels()
	if err != nil {
		t.Error(err)
	}
	if len(channels) == 0 {
		t.Error("no channels returned")
	}
	messages, err := client.ListMessages("DML1K88KV")
	if err != nil {
		t.Error(err)
	}
	if len(messages) == 0 {
		t.Error("no messages returned")
	}
	if _, err := client.ListMessages(""); err == nil {
		t.Error("should error")
	}
	messages, err = client.ListTodayMessages("DML1K88KV")
	if err != nil {
		t.Error(err)
	}
	if len(messages) == 0 {
		t.Error("no messages returned")
	}
	for _, m := range messages {
		if m.BotID != "" {
			if err := client.DeleteBotMessage("DML1K88KV", m.Ts); err != nil {
				t.Error(err)
			}
			break
		}
	}
	if _, err := client.ListTodayMessages(""); err == nil {
		t.Error("should error")
	}
	client.BotToken = ""
	if _, err := client.ListChannels(); err == nil {
		t.Error("should error")
	}
}

func TestClientListDirectMessageChannels(t *testing.T) {
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	client := config.NewClient()
	channels, err := client.ListDirectMessageChannels()
	if err != nil {
		t.Error(err)
	}
	if len(channels) == 0 {
		t.Error("no channels")
	}

	client.BotToken = ""
	if _, err := client.ListDirectMessageChannels(); err == nil {
		t.Error("should error")
	}
}
