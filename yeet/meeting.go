package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/seaspancode/services/elasticsvc"
	"github.com/sirupsen/logrus"

	"github.com/estenssoros/yeetbot/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	meetingCmd.AddCommand(newMeetingCmd)
	meetingCmd.AddCommand(meetingStartCmd)
	meetingCmd.AddCommand(meetingEndCmd)
	meetingCmd.AddCommand(deleteMeetingCmd)
}

var meetingCmd = &cobra.Command{
	Use:   "meeting",
	Short: "create new yeet stuff",
}

var newMeetingCmd = &cobra.Command{
	Use:   "new",
	Short: "create a new meeting",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("must supply config file")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := ioutil.ReadFile(args[0])
		if err != nil {
			return errors.Wrap(err, "readfile")
		}
		config := &client.Config{}
		if err := yaml.Unmarshal(data, config); err != nil {
			return errors.Wrap(err, "unmarshal")
		}
		if len(config.Meetings) == 0 {
			return nil
		}
		if len(args) == 2 {
			for _, m := range config.Meetings {
				if args[1] == m.Name {
					if err := config.NewClient().CreateMeeting(m); err != nil {
						return errors.Wrap(err, "create new meeting from report")
					}
				}
			}
			return nil
		}
		fmt.Printf("found %d meetings\n", len(config.Meetings))
		for _, m := range config.Meetings {
			fmt.Println(m.Name)
		}
		return nil
	},
}

var meetingStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start a meeting",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := client.ConfigFromEnv()
		if err != nil {
			return errors.Wrap(err, "config from env")
		}
		c := config.NewClient()
		meetings, err := c.ListPendingMeetings()
		if err != nil {
			return errors.Wrap(err, "list pending meetings")
		}
		if len(args) == 0 {
			fmt.Printf("found %d meetings\n", len(meetings))
			for _, m := range meetings {
				fmt.Println(m.ID)
			}
			return nil
		}
		for _, m := range meetings {
			if m.ID == args[0] {
				if err := client.StartMeeting(m); err != nil {
					return errors.Wrap(err, "client start meeting")
				}
				m.Started = true
				es := elasticsvc.New(context.Background())
				if err := es.PutOne(meetingIndex, m); err != nil {
					return errors.Wrap(err, "es put one")
				}
				logrus.Infof("started meeting %s", m.ID)
				return nil
			}
		}
		return errors.New("no matching meeting")
	},
}

var meetingEndCmd = &cobra.Command{
	Use:   "end",
	Short: "end a meeting",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := client.ConfigFromEnv()
		if err != nil {
			return errors.Wrap(err, "config from env")
		}
		c := config.NewClient()
		meetings, err := c.ListActiveMeetings()
		if err != nil {
			return errors.Wrap(err, "list active meetings")
		}
		if len(args) == 0 {
			fmt.Printf("found %d meetings in progress\n", len(meetings))
			for _, m := range meetings {
				fmt.Println(m.ID)
			}
			return nil
		}
		for _, m := range meetings {
			if m.ID == args[0] {
				m.Ended = true
				es := elasticsvc.New(context.Background())
				if err := es.PutOne(meetingIndex, m); err != nil {
					return errors.Wrap(err, "es put one")
				}
				logrus.Infof("ended meeting %s", m.ID)
				return nil
			}
		}
		return nil
	},
}

var deleteMeetingCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a meeting",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := client.ConfigFromEnv()
		if err != nil {
			return errors.Wrap(err, "config from env")
		}
		c := config.NewClient()
		meetings, err := c.ListAllMeetings()
		if err != nil {
			return errors.Wrap(err, "client list all meetings")
		}
		if len(args) == 0 {
			fmt.Printf("found %d meetings in progress\n", len(meetings))
			for _, m := range meetings {
				fmt.Println(m.ID)
			}
			return nil
		}
		for _, m := range meetings {
			if m.ID == args[0] {
				es := elasticsvc.New(context.Background())
				if err := es.DeleteOne(meetingIndex, m); err != nil {
					return errors.Wrap(err, "es put one")
				}
				logrus.Infof("ended meeting %s", m.ID)
				return nil
			}
		}
		return nil
	},
}
