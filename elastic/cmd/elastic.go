package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/olivere/elastic"
	uuid "github.com/satori/go.uuid"

	"github.com/cheggaaa/pb"
	"github.com/seaspancode/legere/models"
	"github.com/seaspancode/services/awssvc"

	"github.com/pkg/errors"
	"github.com/seaspancode/services/elasticsvc"
	"github.com/spf13/cobra"
)

var (
	elasticURL = "https://vpc-elasti-test-4xae3alsxtpwfn66t53d46mk64.us-west-2.es.amazonaws.com"
	index      = "yeetbot"
)
var mapping = `
{
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 1
    },
    "mappings": {
        "document": {
            "properties": {
                "user":     					{"type": "text"},
                "report":  						{"type": "text"},
                "date":      					{"type": "text", "format": "dateOptionalTime"},
                "responses":  				{"type": "text"},
                "pending_response":   {"type": "boolean"},
                "responses":  				{"type": "text"}
            }
        }
    }
}
`

func init() {
	ElasticCmd.AddCommand(elasticDeleteCmd)
	ElasticCmd.AddCommand(elasticCreateCmd)
	ElasticCmd.AddCommand(elasticDeleteCreateCmd)
	ElasticCmd.AddCommand(elasticPutCmd)
	ElasticCmd.AddCommand(elasticSearchCmd)
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
	Short: "deletesw and creates an elastic index",
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

func init() {
	elasticPutCmd.AddCommand(elasticPutTweetsCmd)
	elasticPutCmd.AddCommand(elasticPutS3Cmd)
}

var elasticPutCmd = &cobra.Command{
	Use:   "put",
	Short: "put stuff in elastic search",
}

func splitFilePath(s string) []string {
	if strings.HasSuffix(s, "/") {
		s = s[:len(s)-1]
	}
	return strings.Split(s, "/")
}

var elasticPutS3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "put s3 in elastic search",
	RunE: func(cmd *cobra.Command, args []string) error {
		aws := awssvc.New(context.Background())
		if err := aws.Connect(nil); err != nil {
			return err
		}
		aws.Bucket = "ssml-docstore"
		keys, err := aws.List()
		if err != nil {
			return err
		}
		es := elasticsvc.New(context.Background())
		es.SetURL(elasticURL)
		bar := pb.StartNew(len(keys))
		for _, k := range keys {
			if !strings.Contains(k.Name, ".") {
				continue
			}
			domain := ""
			tags := []string{}

			dirs := splitFilePath(k.Name)
			if len(dirs) > 0 {
				domain = dirs[0]
			}
			for _, d := range dirs[1:] {
				if strings.Contains(d, ".") {
					break
				}
				tags = append(tags, strings.ToLower(d))
			}
			_, fileName := filepath.Split(k.Name)
			doc := &models.Document{
				ID:        uuid.Must(uuid.NewV4()),
				FileName:  fileName,
				Paths:     []string{k.Name},
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
				Source:    "manual-upload",
				Domain:    domain,
				Tags:      tags,
			}
			if err := es.PutOne("ssml-docs", doc); err != nil {
				return err
			}
			bar.Increment()
		}
		bar.Finish()

		return nil
	},
}

var elasticPutTweetsCmd = &cobra.Command{
	Use:   "tweets",
	Short: "put tweets in elastic search",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("not allowed")
	},
}

var elasticSearchCmd = &cobra.Command{
	Use: "search",
	RunE: func(cmd *cobra.Command, args []string) error {
		es := elasticsvc.New(context.Background())
		es.SetURL(elasticURL)
		query := elastic.NewPrefixQuery("path", "COMMERCIAL/CUSTOMERS/CMA CGM")
		docs := []*models.Document{}
		if _, err := es.Search("ssml-docs", query, &docs); err != nil {
			return err
		}
		fmt.Println(docs)
		return nil
	},
}
