package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
)

type slackReq struct {
	Type    string `json:"type"`
	Actions []struct {
		Name            string `json:"name"`
		Type            string `json:"type"`
		SelectedOptions []struct {
			Value string `json:"value"`
		} `json:"selected_options"`
	} `json:"actions"`
	CallbackID string `json:"callback_id"`
	Team       struct {
		ID     string `json:"id"`
		Domain string `json:"domain"`
	} `json:"team"`
	Channel struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"channel"`
	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	ActionTs        string `json:"action_ts"`
	MessageTs       string `json:"message_ts"`
	AttachmentID    string `json:"attachment_id"`
	Token           string `json:"token"`
	IsAppUnfurl     bool   `json:"is_app_unfurl"`
	OriginalMessage struct {
		BotID       string `json:"bot_id"`
		Type        string `json:"type"`
		Text        string `json:"text"`
		User        string `json:"user"`
		Ts          string `json:"ts"`
		Team        string `json:"team"`
		Attachments []struct {
			CallbackID string `json:"callback_id"`
			Text       string `json:"text"`
			ID         int    `json:"id"`
			Color      string `json:"color"`
			Actions    []struct {
				ID         string `json:"id"`
				Name       string `json:"name"`
				Text       string `json:"text"`
				Type       string `json:"type"`
				DataSource string `json:"data_source"`
				Options    []struct {
					Text  string `json:"text"`
					Value string `json:"value"`
				} `json:"options"`
			} `json:"actions"`
			Fallback string `json:"fallback"`
		} `json:"attachments"`
	} `json:"original_message"`
	ResponseURL string `json:"response_url"`
	TriggerID   string `json:"trigger_id"`
}

func (r slackReq) String() string {
	ju, _ := json.MarshalIndent(r, "", " ")
	return string(ju)
}

func handleLambda(ctx context.Context, r map[string]interface{}) error {
	body, ok := r["body"]
	if !ok {
		return errors.New("missing request body")
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "json marshal")
	}
	raw = raw[9 : len(raw)-1]
	jsonStr, err := url.QueryUnescape(string(raw))
	if err != nil {
		return errors.Wrap(err, "query unescape")
	}
	req := &slackReq{}
	if err := json.Unmarshal([]byte(jsonStr), req); err != nil {
		return errors.Wrap(err, "json unmarshal")
	}
	fmt.Println(req)
	return nil
}

func main() {
	lambda.Start(handleLambda)
}
