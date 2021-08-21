package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"testing"
	"time"
)

const (
	twitterIndex = "twitter"
)

// Tweet is a structure used for serializing/deserializing data in Elasticsearch.
type Twitter struct {
	User     string                `json:"user"`
	Message  string                `json:"message"`
	Retweets int                   `json:"retweets"`
	Image    string                `json:"image,omitempty"`
	Created  time.Time             `json:"created,omitempty"`
	Tags     []string              `json:"tags,omitempty"`
	Location string                `json:"location,omitempty"`
	Suggest  *elastic.SuggestField `json:"suggest_field,omitempty"`
}

const mapping = `
{
    "settings":{
        "number_of_shards":1,
        "number_of_replicas":0
    },
    "mappings":{
        "properties":{
            "user":{
                "type":"keyword"
            },
            "message":{
                "type":"text",
                "store":true,
                "fielddata":true
            },
            "image":{
                "type":"keyword"
            },
            "created":{
                "type":"date"
            },
            "tags":{
                "type":"keyword"
            },
            "location":{
                "type":"geo_point"
            },
            "suggest_field":{
                "type":"completion"
            }
        }
    }
}`

func TestCreateIndexTwitter(t *testing.T) {
	ctx := context.Background()
	client := GetEsClient()

	exists, err := client.IndexExists(twitterIndex).Do(ctx)
	if err != nil {
		panic(err)
	}
	if !exists {
		fmt.Printf("index twitter is not exist!\n")

		createIndex, err := client.CreateIndex(twitterIndex).BodyString(mapping).Do(ctx)
		if err != nil {
			panic(err)
		}
		if createIndex.Acknowledged {
			fmt.Printf("index twitter: %s", GetLogString(createIndex))
		} else {
			fmt.Printf("index twitter is not Acknowledged")
		}
	} else {
		fmt.Printf("index twitter is exist!\n")
	}
}

func TestInsertTwitter(t *testing.T) {
	ctx := context.Background()
	client := GetEsClient()

	twitter1 := Twitter{
		User:     "Alice",
		Message:  "Take five",
		Retweets: 0,
	}
	put1, err := client.Index().Index(twitterIndex).Id("1").BodyJson(twitter1).Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("put1: %s\n", GetLogString(put1))

	twitter2 := `{"user": "Bob", "message": "It's a Raggy Waltz"}`
	put2, err := client.Index().Index(twitterIndex).Id("2").BodyString(twitter2).Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("put2: %s\n", GetLogStringPretty(put2))
}

func TestQuerySpecificTwitterId(t *testing.T) {
	ctx := context.Background()
	client := GetEsClient()

	get1, err := client.Get().Index(twitterIndex).Id("1").Do(ctx)
	if err != nil {
		panic(err)
	}
	if get1.Found {
		fmt.Printf("get twitter id 1 success:%s\n", GetLogStringPretty(get1))
	}

	// Flush to make sure the documents got written.
	resp, err := client.Flush().Index(twitterIndex).Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("flush success:%s\n", GetLogStringPretty(resp))
}

func TestTermQuery(t *testing.T) {
	ctx := context.Background()
	client := GetEsClient()

	termQuery := elastic.NewTermQuery("user", "Alice") // 严格区分大小写
	searchResult, err := client.Search().
		Index(twitterIndex).Query(termQuery).Sort("user", true).
		From(0).Size(10).Pretty(true).Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("searchResult:%s\n", GetLogStringPretty(searchResult))
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
}
