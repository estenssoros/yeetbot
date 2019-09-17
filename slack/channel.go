package slack

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Channel struct {
	ID             string
	Name           string
	IsChannel      bool     `json:"is_channel"`
	Created        int      `json:"created"`
	IsArchived     bool     `json:"is_archived"`
	IsGeneral      bool     `json:"is_general"`
	Unlinked       int      `json:"unlinked"`
	Creator        string   `json:"creator"`
	NameNormalized string   `json:"name_normalized"`
	IsShared       bool     `json:"is_shared"`
	IsOrgShared    bool     `json:"is_org_shared"`
	IsMember       bool     `json:"is_member"`
	IsPrivate      bool     `json:"is_private"`
	IsMpim         bool     `json:"is_mpim"`
	Members        []string `json:"members"`
	IsIm           bool     `json:"is_im"`
	User           string   `json:"user"`
	IsUserDeleted  bool     `json:"is_user_deleted"`
	Priority       int      `json:"priority"`
	Topic          struct {
		Value   string `json:"value"`
		Creator string `json:"creator"`
		LastSet int    `json:"last_set"`
	} `json:"topic"`
	Purpose struct {
		Value   string `json:"value"`
		Creator string `json:"creator"`
		LastSet int    `json:"last_set"`
	} `json:"purpose"`
	PreviousNames []interface{} `json:"previous_names"`
	NumMembers    int           `json:"num_members"`
}

func (c Channel) String() string {
	ju, _ := json.MarshalIndent(c, "", " ")
	return string(ju)
}

type listChannelResponse struct {
	OK       bool       `json:"ok"`
	Channels []*Channel `json:"channels"`
	Error    string     `json:"error"`
}

func (a *API) ListChannels() ([]*Channel, error) {
	data, err := newRequest(ChannelsList).
		addParam("token", a.botToken).
		Get()
	if err != nil {
		return nil, errors.Wrap(err, "slack api request")
	}
	resp := &listChannelResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	if !resp.OK {
		return nil, errors.New(resp.Error)
	}
	return resp.Channels, nil
}
