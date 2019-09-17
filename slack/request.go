package slack

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type apiRequest struct {
	url     string
	suffix  string
	params  map[string]string
	body    []byte
	headers map[string]string
}

func newRequest(suffix string) *apiRequest {
	return &apiRequest{
		url:     "https://slack.com/api/",
		suffix:  suffix,
		params:  map[string]string{},
		headers: map[string]string{},
	}
}

func (r *apiRequest) addParam(key, value string) *apiRequest {
	r.params[key] = value
	return r
}

func (r *apiRequest) addHeader(key, value string) *apiRequest {
	r.headers[key] = value
	return r
}

func (r *apiRequest) addBody(v interface{}) *apiRequest {
	ju, _ := json.Marshal(v)
	r.body = ju
	return r
}

func (r *apiRequest) Get() ([]byte, error) {
	return r.Do(http.MethodGet)
}

func (r *apiRequest) Post() ([]byte, error) {
	return r.Do(http.MethodPost)
}

func (r *apiRequest) craftURL() (string, error) {
	u, err := url.Parse(r.url + r.suffix)
	if err != nil {
		return "", errors.Wrap(err, "url parse")
	}
	q := u.Query()
	for k, v := range r.params {
		q.Set(k, v)

	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (r *apiRequest) Do(method string) ([]byte, error) {
	url, err := r.craftURL()
	if err != nil {
		return nil, errors.Wrap(err, "craft url")
	}
	logrus.Println(method, url)
	var (
		req *http.Request
	)
	if r.body != nil {
		req, err = http.NewRequest(method, url, bytes.NewReader(r.body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, errors.Wrap(err, "http new request")
	}
	{
		req.Header.Set("Content-Type", "application/json;charset=iso-8859-1")
	}
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
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
	return data, nil
}
