package client

import (
	"context"
	"fmt"
	"time"

	"github.com/estenssoros/yeetbot/slack"
	"github.com/olivere/elastic"
	uuid "github.com/satori/go.uuid"
	"github.com/seaspancode/services/elasticsvc"
)

// Response stored in elastic
type Response struct {
	ID              uuid.UUID `json:"id"`
	User            string    `json:"user"`
	Report          string    `json:"report"`
	Date            time.Time `json:"date"`
	PendingResponse bool      `json:"pending_response"`
	Responses       []string  `json:"responses"`
}

// GetResponsesByUser gets user response from elastic search
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

// AddUserResponse puts user response to elastic search
func (c *Client) AddUserResponse(response *Response) error {
	es := elasticsvc.New(context.Background())
	es.SetURL(c.ElasticURL)
	fmt.Println("adding")
	return es.PutOne(c.ElasticIndex, response)
}

// RecordResponse adds response to responses only if it has pending status
// sets pending status to false
// and returns total number of responses recorded
func (c *Client) RecordResponse(user *slack.User, response *Response, message string) (int, error) {
	return 0, nil
}

// CompleteResponse sets pending status to false from response
// and sends "thank you" message to user
func (c *Client) CompleteResponse(user *slack.User, response *Response) error {
	return nil
}
