package client

import (
	"context"
	"fmt"

	"github.com/coreos/etcd/client"
	"github.com/estenssoros/yeetbot/slack"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"github.com/seaspancode/services/elasticsvc"
)

// User naive user data structure
type User struct {
	Name string
	ID   string `yaml:"slackID,omitempty"`
}

func (c *Client) HasUser(user *slack.User) bool {
	for _, u := range c.Users {
		if u.Name == user.Name {
			return true
		}
	}
	return false
}

// HasUserStartedReport checks to see if a report has already been started today
func (c *Client) HasUserStartedReport(user *slack.User) (bool, error) {
	es := elasticsvc.New(context.Background())
	es.SetURL(c.ElasticURL)
	query := elastic.NewPrefixQuery("user_id", user.ID)
	responses := []*client.Response{}
	if err := es.GetMany(c.ElasticIndex, query, &responses); err != nil {
		return false, errors.Wrap(err, "failed to get many")
	}
	fmt.Println(responses)
	return false, nil
}
