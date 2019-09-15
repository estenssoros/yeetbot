package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/estenssoros/yeetbot/client"
	"github.com/olivere/elastic"
	uuid "github.com/satori/go.uuid"
	"github.com/seaspancode/services/elasticsvc"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	elasticURL   = os.Getenv("ELASTIC_URL")
	elasticIndex = "yeetbot"
	testIndex    = "yeetbot-test"
	mapping      = `
{
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 1
    }, 
    "mappings": {
        "response": {
            "properties": {
				        "team":           {"type": "text"},
                "user_id":     		{"type": "text"},
                "event_ts":     	{"type": "integer"},
                "date":      			{"type": "date", "format": "dateOptionalTime"},
                "question":   		{"type": "text"},
                "text":  					{"type": "text"}
            }
        }
    }
}
`
)

func init() {
	elasticCmd.AddCommand(elasticDeleteCmd)
	elasticCmd.AddCommand(elasticCreateCmd)
	elasticCmd.AddCommand(elasticDeleteCreateCmd)
	elasticCmd.AddCommand(elasticSearchCmd)
	elasticCmd.AddCommand(elasticPutCmd)

	elasticCreateCmd.Flags().Bool("test", false, "Uses test index")
	elasticDeleteCreateCmd.Flags().Bool("test", false, "Uses test index")
	elasticSearchCmd.Flags().Bool("test", false, "Uses test index")
	elasticPutCmd.Flags().Bool("test", false, "Uses test index")

	elasticSearchCmd.Flags().StringP("name", "n", "berto", "User id to search or add")
	elasticPutCmd.Flags().StringP("name", "n", "berto", "User id to search or add")
	elasticPutCmd.Flags().StringP("report", "r", "daily_standup", "Report to add")
	elasticPutCmd.Flags().StringP("message", "m", "good", "Message to add")
}

var elasticCmd = &cobra.Command{
	Use:   "elastic",
	Short: "do elastic stuff",
}

var elasticDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete an elastic index",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := elasticsvc.New(context.Background())
		svc.SetURL(elasticURL)
		return errors.Wrap(svc.DeleteIndex(elasticIndex), "delete index")
	},
}

var elasticCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create an elastic index",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := elasticsvc.New(context.Background())
		svc.SetURL(elasticURL)
		index := getIndex(cmd)
		return errors.Wrap(svc.CreateIndex(index, mapping), "create index")
	},
}

var elasticDeleteCreateCmd = &cobra.Command{
	Use:   "delete-create",
	Short: "deletes and creates an elastic index",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := elasticsvc.New(context.Background())
		svc.SetURL(elasticURL)
		index := getIndex(cmd)
		if err := svc.DeleteIndex(index); err != nil {
			return errors.Wrap(err, "delete index")
		}
		return errors.Wrap(svc.CreateIndex(index, mapping), "create index")
	},
}

var elasticSearchCmd = &cobra.Command{
	Use: "search",
	RunE: func(cmd *cobra.Command, args []string) error {
		es := elasticsvc.New(context.Background())
		es.SetURL(elasticURL)
		name, _ := cmd.Flags().GetString("name")
		query := elastic.NewPrefixQuery("user_id", name)
		responses := []*client.Response{}
		index := getIndex(cmd)
		if err := es.GetMany(index, query, &responses); err != nil {
			return err
		}
		logrus.Infof("results for %s: ", index)
		for i := range responses {
			logrus.Infof("%+v", responses[i])
		}
		return nil
	},
}

var elasticPutCmd = &cobra.Command{
	Use: "put",
	RunE: func(cmd *cobra.Command, args []string) error {
		es := elasticsvc.New(context.Background())
		es.SetURL(elasticURL)

		name, _ := cmd.Flags().GetString("name")
		report, _ := cmd.Flags().GetString("report")
		response, _ := cmd.Flags().GetString("message")
		doc := &client.Response{
			ID:       uuid.Must(uuid.NewV4()),
			Report:   report,
			Channel:  "daily-standup",
			UserID:   name,
			EventTS:  fmt.Sprint(time.Now().Unix()),
			Date:     time.Now(),
			Question: "How do you feel?",
			Text:     response,
		}
		index := getIndex(cmd)
		if err := es.PutOne(index, doc); err != nil {
			return errors.Wrap(err, "adding response")
		}
		logrus.Infof("added %s to %s", doc.ID, index)
		return nil
	},
}

func getIndex(cmd *cobra.Command) string {
	useTest, _ := cmd.Flags().GetBool("test")
	if useTest {
		return testIndex
	}
	return elasticIndex
}
