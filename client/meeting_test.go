package client

import (
	"testing"
	"time"

	"github.com/estenssoros/yeetbot/slack"
	"github.com/stretchr/testify/assert"
)

func TestMeetingStringer(t *testing.T) {
	m := &Meeting{ID: "asdf"}
	if s := m.String(); len(s) == 0 {
		t.Error("could not stringify meerintg")
	}
}

func TestMeetingSetID(t *testing.T) {
	m := &Meeting{
		Team:           "denver",
		Name:           "test",
		ScheduledStart: time.Date(2006, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	m.SetID()
	if want, have := "denver-test-20060101", m.ID; want != have {
		t.Errorf("have: %s, want %s", want, have)
	}
}

func TestMeetingESType(t *testing.T) {
	m := &Meeting{}
	if want, have := "meeting", m.EsType(); want != have {
		t.Errorf("have: %s, want %s", want, have)
	}
}

func TestMeetingHasUser(t *testing.T) {
	m := &Meeting{
		Users: []*slack.User{
			&slack.User{ID: "asdf"},
		},
	}
	if !m.HasUser(&slack.User{ID: "asdf"}) {
		t.Error("coult not find user")
	}
	if m.HasUser(&slack.User{ID: "fdsa"}) {
		t.Error("found not existing user")
	}
}

func TestMeetingCraftEvents(t *testing.T) {
	m := &Meeting{
		Questions: []*Question{
			&Question{},
		},
	}
	events := m.CraftEvents()
	assert.Equal(t, 1, len(events))
}

func TestMeetingCreateReports(t *testing.T) {
	m := &Meeting{
		Users: []*slack.User{
			&slack.User{},
		},
		Schedule: &Schedule{
			Cron: "* * * * * * *",
		},
	}
	reports := m.CreateReports()
	assert.Equal(t, len(reports), 1)
}
