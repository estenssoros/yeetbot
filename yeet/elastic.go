package main

import (
	"context"
	"os"

	"github.com/seaspancode/services/elasticsvc"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	elasticURL   = os.Getenv("ELASTIC_URL")
	reportIndex  = "yeetreport"
	meetingIndex = "yeetmeet"
	testIndex    = "yeetbot-test"
	indices      = []string{
		reportIndex,
		meetingIndex,
	}
	reportMapping = `
{
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 1
    },
    "mappings": {
        "report": {
            "properties": {
                "meetingID": {
                    "type": "text"
                },
                "userID": {
                    "type": "text"
                },
                "createdAt": {
                    "type": "date",
                    "format": "dateOptionalTime"
                },
                "events": {
                    "type": "object",
                    "properties": {
                        "question": {
                            "type": "text"
                        },
                        "response": {
                            "type": "text"
                        },
                        "eventTS": {
                            "type": "integer"
                        }
                    }
                },
                "done": {
                    "type": "boolean"
                }
            }
        }
    }
}
	`
	meetingMapping = `
{
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 1
    },
    "mappings": {
        "meeting": {
            "properties": {
                "name": {
                    "type": "text"
                },
                "channel": {
                    "type": "text"
                },
                "team": {
                    "type": "text"
                },
                "scheduledStart": {
                    "type": "date",
                    "format": "dateOptionalTime"
                },
                "users": {
                    "type": "object",
                    "properties": {
                        "name": {
                            "type": "text"
                        },
                        "id": {
                            "type": "text"
                        }
                    }
                },
                "questions": {
                    "type": "object",
                    "properties": {
                        "text": {
                            "type": "text"
                        },
                        "color": {
                            "type": "text"
                        },
                        "options": {
                            "type": "text"
                        }
                    }
                },
                "started": {
                    "type": "boolean"
                },
                "ended": {
                    "type": "boolean"
                }
            }
        }
    }
}
`
	indexMap = map[string]string{
		reportIndex:  reportMapping,
		meetingIndex: meetingMapping,
	}
)

func init() {
	elasticCmd.AddCommand(elasticDeleteCmd)
	elasticCmd.AddCommand(elasticCreateCmd)
	elasticCmd.AddCommand(elasticDeleteCreateCmd)
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
		for _, index := range indices {
			if err := svc.DeleteIndex(index); err != nil {
				return errors.Wrap(err, index)
			}
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
		for _, index := range indices {
			if err := svc.CreateIndex(index, indexMap[index]); err != nil {
				return errors.Wrap(err, index)
			}
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

		for _, index := range indices {
			if err := svc.DeleteIndex(index); err != nil {
				return errors.Wrap(err, index)
			}
		}

		for _, index := range indices {
			if err := svc.CreateIndex(index, indexMap[index]); err != nil {
				return errors.Wrap(err, index)
			}
		}
		return nil
	},
}
