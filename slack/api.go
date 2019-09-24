package slack

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type API struct {
	botToken string
	verbose  bool
}

func New(token string) *API {
	return &API{
		botToken: token,
	}
}

func (a *API) SetVerbose(v bool) {
	a.verbose = v
}

func (a *API) newRequest(suffix string) *apiRequest {
	r := newRequest(suffix)
	r.SetVerbose(a.verbose)
	return r
}

func (a *API) SendMessage(msg *Message) error {
	url := slackurl + "/" + ChatPostMessage
	if a.verbose {
		logrus.Infof("%s %s", http.MethodPost, url)
	}
	ju, _ := json.Marshal(msg)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(ju))
	req.Header.Set("Authorization", "Bearer "+a.botToken)
	req.Header.Set("Content-Type", "application/json;charset=iso-8859-1")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "http default client do")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read responses")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("bad status code: %d, %s", resp.StatusCode, string(b))
	}
	slackresp := &Response{}
	if err := json.Unmarshal(b, slackresp); err != nil {
		return errors.Wrap(err, "unmarshal slack resp")
	}
	if slackresp.Warning != "" {
		logrus.Warning(slackresp.Warning)
	}
	if !slackresp.OK {
		logrus.Error(string(b))
		return errors.Wrap(errors.New(slackresp.Error), "slack response")
	}
	return nil
}

type listUsersResponse struct {
	OK      bool    `json:"ok"`
	Members []*User `json:"members"`
	Error   string  `json:"error"`
}

// ListUsers lists users in workspace
func (a *API) ListUsers() ([]*User, error) {
	data, err := a.newRequest(UsersList).
		addParam("token", a.botToken).
		Get()
	if err != nil {
		return nil, errors.Wrap(err, "slack api request")
	}
	listResponse := &listUsersResponse{}
	if err := json.Unmarshal(data, &listResponse); err != nil {
		return nil, errors.Wrap(err, "unmarshal request: %s")
	}
	if !listResponse.OK {
		return nil, errors.New(listResponse.Error)
	}
	return listResponse.Members, nil
}
