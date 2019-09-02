package main

import (
	"fmt"
	"testing"

	"github.com/estenssoros/yeetbot/client"
)

func TestClient(t *testing.T) {
	c, err := client.ClientFromConfig("test-config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if err := c.AddTeamFromFile("test-team.yaml"); err != nil {
		t.Fatal(err)
	}
	fmt.Println(c)
	for _, user := range c.Team.Users {
		if err := c.SendGreeting(user); err != nil {
			t.Fatal(err)
		}
		if err := c.Run(user); err != nil {
			t.Fatal(err)
		}
	}
}
