package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/estenssoros/yeetbot/models"
	"github.com/estenssoros/yeetbot/slack"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type client struct {
	url   string
	token string
}

func newClient() *client {
	return &client{
		url:   slackurl,
		token: token,
	}
}

func (c *client) craftURL(suffix string) (string, error) {
	u, err := url.Parse(c.url + "/" + suffix)
	if err != nil {
		return "", errors.Wrap(err, "url parse")
	}
	q := u.Query()
	q.Set("token", c.token)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

type listUsersResponse struct {
	OK      bool           `json:"ok"`
	Members []*models.User `json:"members"`
}

func (c *client) ListUsers() ([]*models.User, error) {
	url, err := c.craftURL(slack.UsersList)
	if err != nil {
		return nil, errors.Wrap(err, "client craft url")
	}
	logrus.Info(http.MethodGet, url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "http new request")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http do request")
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read resp")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("bad response %d: %s", resp.StatusCode, string(data))
	}
	listResponse := &listUsersResponse{}
	if err := json.Unmarshal(data, &listResponse); err != nil {
		return nil, errors.Wrap(err, "unmarshal request")
	}
	return listResponse.Members, nil
}
