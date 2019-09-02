package slack

import (
	"encoding/json"
)

type Response struct {
	OK       bool              `json:"ok"`
	Error    string            `json:"error"`
	Warning  string            `json:"warning"`
	MetaData *ResponseMetaData `json:"response_metadata"`
}

func (r Response) String() string {
	ju, _ := json.Marshal(r)
	return string(ju)
}

type ResponseMetaData struct {
	Warnings []string `json:"warnings"`
}
