package slack

import "encoding/json"

type User struct {
	ID                string   `json:"id" yaml:"id,omitempty"`
	TeamID            string   `json:"team_id" yaml:"team_id,omitempty"`
	Name              string   `json:"name" yaml:"name,omitempty"`
	Deleted           bool     `json:"deleted" yaml:"deleted,omitempty"`
	Color             string   `json:"color" yaml:"color,omitempty"`
	RealName          string   `json:"real_name" yaml:"real_name,omitempty"`
	Tz                string   `json:"tz" yaml:"tz,omitempty"`
	TzLabel           string   `json:"tz_label" yaml:"tz_label,omitempty"`
	TzOffset          int      `json:"tz_offset" yaml:"tz_offset,omitempty"`
	Profile           *Profile `json:"profile" yaml:"profile,omitempty"`
	IsAdmin           bool     `json:"is_admin" yaml:"is_admin,omitempty"`
	IsOwner           bool     `json:"is_owner" yaml:"is_owner,omitempty"`
	IsPrimaryOwner    bool     `json:"is_primary_owner" yaml:"is_primary_owner,omitempty"`
	IsRestricted      bool     `json:"is_restricted" yaml:"is_restricted,omitempty"`
	IsUltraRestricted bool     `json:"is_ultra_restricted" yaml:"is_ultra_restricted,omitempty"`
	IsBot             bool     `json:"is_bot" yaml:"is_bot,omitempty"`
	IsAppUser         bool     `json:"is_app_user" yaml:"is_app_user,omitempty"`
	Updated           int      `json:"updated" yaml:"updated,omitempty"`
}

type Profile struct {
	Title                 string `json:"title" yaml:"title,omitempty"`
	Phone                 string `json:"phone" yaml:"phone,omitempty"`
	Skype                 string `json:"skype" yaml:"skype,omitempty"`
	RealName              string `json:"real_name" yaml:"real_name,omitempty"`
	RealNameNormalized    string `json:"real_name_normalized" yaml:"real_name_normalized,omitempty"`
	DisplayName           string `json:"display_name" yaml:"display_name,omitempty"`
	DisplayNameNormalized string `json:"display_name_normalized" yaml:"display_name_normalized,omitempty"`
	StatusText            string `json:"status_text" yaml:"status_text,omitempty"`
	StatusEmoji           string `json:"status_emoji" yaml:"status_emoji,omitempty"`
	StatusExpiration      int    `json:"status_expiration" yaml:"status_expiration,omitempty"`
	AvatarHash            string `json:"avatar_hash" yaml:"avatar_hash,omitempty"`
	AlwaysActive          bool   `json:"always_active" yaml:"always_active,omitempty"`
	FirstName             string `json:"first_name" yaml:"first_name,omitempty"`
	LastName              string `json:"last_name" yaml:"last_name,omitempty"`
	Image24               string `json:"image_24" yaml:"image_24,omitempty"`
	Image32               string `json:"image_32" yaml:"image_32,omitempty"`
	Image48               string `json:"image_48" yaml:"image_48,omitempty"`
	Image72               string `json:"image_72" yaml:"image_72,omitempty"`
	Image192              string `json:"image_192" yaml:"image_192,omitempty"`
	Image512              string `json:"image_512" yaml:"image_512,omitempty"`
	StatusTextCanonical   string `json:"status_text_canonical" yaml:"status_text_canonical,omitempty"`
	Team                  string `json:"team" yaml:"team,omitempty"`
}

func (p Profile) String() string {
	ju, _ := json.MarshalIndent(p, "", " ")
	return string(ju)
}

func (u *User) Template(t string) (string, error) {
	return ``, nil
}
