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
	"github.com/spf13/cobra"
)

var (
	elasticURL = ""
	index      = "yeetbot"
	mapping    = `
{
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 1
    },
    "mappings": {
        "Response": {
            "properties": {
                "user":     					{"type": "text"},
                "report":  						{"type": "text"},
                "date":      					{"type": "date", "format": "dateOptionalTime"},
                "pending_response":   {"type": "boolean"},
                "responses":  				{"type": "text"}
            }
        }
    }
}
`
)

func init() {
	elasticURL = os.Getenv("ELASTIC_URL")
	ElasticCmd.AddCommand(elasticDeleteCmd)
	ElasticCmd.AddCommand(elasticCreateCmd)
	ElasticCmd.AddCommand(elasticDeleteCreateCmd)
	ElasticCmd.AddCommand(elasticSearchCmd)
	ElasticCmd.AddCommand(elasticPutCmd)

	elasticSearchCmd.Flags().StringP("name", "n", "", "Username to search or add")
	elasticPutCmd.Flags().StringP("report", "r", "", "Report to add")
	elasticPutCmd.Flags().StringP("message", "m", "", "Message to add")
}

var ElasticCmd = &cobra.Command{
	Use:   "elastic",
	Short: "do elastic stuff",
}

var elasticDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete an elastic index",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := elasticsvc.New(context.Background())
		svc.SetURL(elasticURL)
		if err := svc.DeleteIndex(index); err != nil {
			return errors.Wrap(err, "delete index")
		}
		return nil
	},
}

var elasticCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create an elastic index",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := elasticsvc.New(context.Background())
		svc.SetURL(elasticURL)
		if err := svc.CreateIndex(index, mapping); err != nil {
			return errors.Wrap(err, "create index")
		}
		return nil
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
		if err := svc.CreateIndex(index, mapping); err != nil {
			return errors.Wrap(err, "create index")
		}
		return nil
	},
}

var elasticSearchCmd = &cobra.Command{
	Use: "search",
	RunE: func(cmd *cobra.Command, args []string) error {
		es := elasticsvc.New(context.Background())
		es.SetURL(elasticURL)
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			name = "berto"
		}
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
		if name == "" {
			name = "berto"
		}
		report, _ := cmd.Flags().GetString("report")
		if report == "" {
			report = "daily_standup"
		}
		response, _ := cmd.Flags().GetString("message")
		if response == "" {
			response = "good"
		}
		doc := &client.Response{
			ID:              uuid.NewV4(),
			User:            name,
			Report:          report,
			Date:            time.Now(),
			Responses:       []string{response},
			PendingResponse: true,
		}
		if err := es.PutOne(index, doc); err != nil {
			return errors.Wrap(err, "adding response")
		}
		fmt.Println("added")
		return nil
	},
}
