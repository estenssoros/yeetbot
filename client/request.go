package client

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type apiRequest struct {
	url    string
	suffix string
	params map[string]string
}

func newAPIRequest(suffix string) *apiRequest {
	return &apiRequest{
		url:    "https://slack.com/api/",
		suffix: suffix,
		params: map[string]string{},
	}
}

func (r *apiRequest) addParam(key, value string) *apiRequest {
	r.params[key] = value
	return r
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

func (r *apiRequest) Get() ([]byte, error) {
	return r.Do(http.MethodGet)
}

func (r *apiRequest) Do(method string) ([]byte, error) {
	url, err := r.craftURL()
	if err != nil {
		return nil, errors.Wrap(err, "craft url")
	}
	log.Println(method, url)
	req, err := http.NewRequest(method, url, nil)
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
	return data, nil
}
