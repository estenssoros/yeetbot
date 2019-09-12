package client

import "github.com/pkg/errors"

// User naive user data structure
type User struct {
	Name    string
	SlackID string `yaml:"slackID,omitempty"`
}

// UserMap for dealing with the user data map
type UserMap map[string][]*Report

// LatestReport gets the lastest report from the user data map
func (u UserMap) LatestReport(userName string) (*Report, error) {
	reports, ok := u[userName]
	if !ok {
		return nil, errors.Errorf("missing %s", userName)
	}
	switch len(reports) {
	case 0:
		return nil, errors.New("no reports for user")
	case 1:
		return reports[0], nil
	}
	report := reports[0]
	reportTime, err := report.TodayTime()
	if err != nil {
		return nil, err
	}
	for _, r := range reports {
		rTime, err := r.TodayTime()
		if err != nil {
			return nil, err
		}
		if rTime.After(reportTime) {
			report = r
			reportTime = rTime
		}
	}
	return report, nil
}

// HasUser does there username exist in the user map
func (u UserMap) HasUser(userName string) bool {
	_, ok := u[userName]
	return ok
}
