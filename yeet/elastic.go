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
	elasticURL = os.Getenv("ELASTIC_URL")
	index      = "yeetbot"
)

var mapping = `
{
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 1
    },
    "mappings": {
        "response": {
            "properties": {
                "user":     		{"type": "text"},
                "report":  			{"type": "text"},
                "date":      		{"type": "date", "format": "dateOptionalTime"},
                "pending_response": {"type": "boolean"},
                "responses":  		{"type": "text"}
            }
        }
    }
}
`

func init() {
	elasticCmd.AddCommand(elasticDeleteCmd)
	elasticCmd.AddCommand(elasticCreateCmd)
	elasticCmd.AddCommand(elasticDeleteCreateCmd)
	elasticCmd.AddCommand(elasticSearchCmd)
	elasticCmd.AddCommand(elasticPutCmd)

	elasticSearchCmd.Flags().StringP("name", "n", "berto", "Username to search or add")
	elasticPutCmd.Flags().StringP("report", "r", "daily_standu", "Report to add")
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
		return errors.Wrap(svc.DeleteIndex(index), "delete index")
	},
}

var elasticCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create an elastic index",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := elasticsvc.New(context.Background())
		svc.SetURL(elasticURL)
		return errors.Wrap(svc.CreateIndex(index, mapping), "create index")
	},
}

var elasticDeleteCreateCmd = &cobra.Command{
	Use:   "delete-create",
	Short: "deletes and creates an elastic index",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := elasticsvc.New(context.Background())
		svc.SetURL(elasticURL)
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
		query := elastic.NewPrefixQuery("user", name)
		responses := []*client.Response{}
		if _, err := es.Search(index, query, &responses); err != nil {
			return err
		}
		for i := range responses {
			fmt.Printf("%+v", responses[i])
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
			ID:              uuid.Must(uuid.NewV4()),
			User:            name,
			Report:          report,
			Date:            time.Now(),
			Responses:       []string{response},
			PendingResponse: true,
		}
		if err := es.PutOne(index, doc); err != nil {
			return errors.Wrap(err, "adding response")
		}
		logrus.Infof("added %s", doc.ID)
		return nil
	},
}
