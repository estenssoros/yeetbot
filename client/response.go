package client

import (
	"context"
	"fmt"
	"time"

	"github.com/estenssoros/yeetbot/slack"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/seaspancode/services/elasticsvc"
)

// Response stored in elastic
type Response struct {
	ID       uuid.UUID `json:"id"`
	Report   string    `json:"report"`
	Channel  string    `json:"channel"`
	UserID   string    `json:"user_id"`
	EventTS  int64     `json:"event_ts"`
	Date     time.Time `json:"date"`
	Question string    `json:"question"`
	Text     string    `json:"text"`
}

func (r Response) EsType() string {
	return `response`
}

type RecordResponseInput struct {
	Question *Question
	User     *slack.User
	Text     string
}

// GetUserResponse gets user response from elastic search
// list most recent responses on channel
func (c *Client) GetUserResponse(user *slack.User) (*Response, error) {
	es := elasticsvc.New(context.Background())
	es.SetURL(c.ElasticURL)
	reports := []*Response{}
	fmt.Println(user.RealName)
	query := elastic.NewPrefixQuery("user", user.RealName)
	// {
	// 	query = query.Must(elastic.NewQueryStringQuery(user.RealName).DefaultField("user"))
	// query = query.Must(elastic.NewQueryStringQuery(c.Report.Name).DefaultField("report"))
	// query = query.Must(elastic.NewQueryStringQuery(time.Now().String()).DefaultField("date"))
	// }
	fmt.Println(user.RealName, c.ElasticIndex, c.ElasticURL)
	es.Search(c.ElasticIndex, query, &reports)
	fmt.Println(reports)
	return reports[len(reports)-1], nil
}

// RecordResponse adds response to responses
func (c *Client) RecordResponse(input *RecordResponseInput) error {
	resp := &Response{
		Report:   c.Report.Name,
		Channel:  c.Report.Channel,
		UserID:   input.User.ID,
		EventTS:  time.Now().Unix(),
		Date:     time.Now(),
		Question: input.Question.Text,
		Text:     input.Text,
	}
	es := elasticsvc.New(context.Background())
	es.SetURL(c.ElasticURL)
	if err := es.PutOne(c.ElasticIndex, resp); err != nil {
		return errors.Wrap(err, "adding response")
	}
	return nil
}

// CompleteResponse sets pending status to false from response
// and sends "thank you" message to user
func (c *Client) CompleteResponse(user *slack.User, response *Response) error {
	// TODO: this
	return nil
}
