package client

import (
	"fmt"
	"os"
	"testing"
)

func TestReadConfig(t *testing.T) {
	path := os.Getenv("YEETBOT_CONFIG")
	if path == "" {
		t.Fatal("missing YEETBOT_CONFIG env")
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	_, err = ConfigFromReader(f)
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadConfigEnv(t *testing.T) {
	if _, err := ConfigFromEnv(); err != nil {
		t.Fatal(err)
	}
}

func TestReportTodayTime(t *testing.T) {
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	r := config.Reports[0]
	fmt.Println(r.TodayTime())
}

// TODO write a test for multiple reports and getting latest, etc...
// TODO how do we make sure that whatever this deploys to has the timezone database?
// TODO i'm not sure how to check a users local time. maybe a func (u *User) ReportTime(report) that looks at the reports schedule?
