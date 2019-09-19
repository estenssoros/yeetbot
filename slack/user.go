package slack

import (
	"bytes"
	"encoding/json"
	"text/template"

	"github.com/pkg/errors"
)

type User struct {
	ID                string  `json:"id"`
	TeamID            string  `json:"team_id"`
	Name              string  `json:"name"`
	Deleted           bool    `json:"deleted"`
	Color             string  `json:"color"`
	RealName          string  `json:"real_name"`
	TZ                string  `json:"tz,omitempty"`
	TZLabel           string  `json:"tz_label"`
	TZOffset          int     `json:"tz_offset"`
	Profile           Profile `json:"profile"`
	IsBot             bool    `json:"is_bot"`
	IsAdmin           bool    `json:"is_admin"`
	IsOwner           bool    `json:"is_owner"`
	IsPrimaryOwner    bool    `json:"is_primary_owner"`
	IsRestricted      bool    `json:"is_restricted"`
	IsUltraRestricted bool    `json:"is_ultra_restricted"`
	IsStranger        bool    `json:"is_stranger"`
	IsAppUser         bool    `json:"is_app_user"`
	IsInvitedUser     bool    `json:"is_invited_user"`
	Has2FA            bool    `json:"has_2fa"`
	HasFiles          bool    `json:"has_files"`
	Presence          string  `json:"presence"`
	Locale            string  `json:"locale"`
}

func (u User) String() string {
	ju, _ := json.MarshalIndent(u, "", " ")
	return string(ju)
}

type Profile struct {
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	RealName              string `json:"real_name"`
	RealNameNormalized    string `json:"real_name_normalized"`
	DisplayName           string `json:"display_name"`
	DisplayNameNormalized string `json:"display_name_normalized"`
	Email                 string `json:"email"`
	Skype                 string `json:"skype"`
	Phone                 string `json:"phone"`
	Image24               string `json:"image_24"`
	Image32               string `json:"image_32"`
	Image48               string `json:"image_48"`
	Image72               string `json:"image_72"`
	Image192              string `json:"image_192"`
	ImageOriginal         string `json:"image_original"`
	Title                 string `json:"title"`
	BotID                 string `json:"bot_id,omitempty"`
	ApiAppID              string `json:"api_app_id,omitempty"`
	StatusText            string `json:"status_text,omitempty"`
	StatusEmoji           string `json:"status_emoji,omitempty"`
	StatusExpiration      int    `json:"status_expiration"`
	Team                  string `json:"team"`
}

func (p Profile) String() string {
	ju, _ := json.MarshalIndent(p, "", " ")
	return string(ju)
}

// Template templates a user struct onto a string
func (u *User) Template(t string) (string, error) {
	tmpl := template.Must(template.New("").Parse(t))
	var b bytes.Buffer
	if err := tmpl.Execute(&b, u); err != nil {
		return "", errors.Wrap(err, "tmpl execute")
	}
	return b.String(), nil
}

type userInfoResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
	User  *User  `json:"user"`
}

// GetUserByID searches for a user by iD
func (a *API) GetUserByID(userID string) (*User, error) {
	data, err := a.newRequest(UsersInfo).
		addParam("user", userID).
		addParam("include_local", "true").
		addParam("token", a.botToken).
		Get()
	if err != nil {
		return nil, errors.Wrap(err, "api request")
	}
	resp := &userInfoResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	if !resp.OK {
		return nil, errors.New(resp.Error)
	}
	return resp.User, nil
}
