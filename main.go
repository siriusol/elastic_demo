package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/olivere/elastic/v7"
	"github.com/olivere/elastic/v7/config"
)

func main() {
	ctx := context.Background()
	client := GetEsClient()
	confJson := readConfJson()
	info, code, err := client.Ping(confJson.Url).Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ES returned with code %d and version %s\n", code, info.Version.Number)
	fmt.Printf("ES info: %s\n", GetLogString(info))

	esVersion, err := client.ElasticsearchVersion(confJson.Url)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ES version %s\n", esVersion)

	client.IndexGet()
}

func initClientFromConfig(confJson EsConfJson) (*elastic.Client, error) {
	sniff := false
	esConfig := config.Config{
		URL:      confJson.Url,
		Username: confJson.Username,
		Password: confJson.Password,
		Sniff:    &sniff,
	}
	return elastic.NewClientFromConfig(&esConfig)
}

func readConfJson() EsConfJson {
	bytes, err := ioutil.ReadFile("./conf.json")
	if err != nil {
		panic(err)
	}
	var confJson EsConfJson
	err = json.Unmarshal(bytes, &confJson)
	if err != nil {
		panic(err)
	}
	return confJson
}

func GetEsClient() *elastic.Client {
	confJson := readConfJson()
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(confJson.Url),
		elastic.SetBasicAuth(confJson.Username, confJson.Password),
		elastic.SetSniff(false),
	}
	client, err := elastic.NewClient(opts...)
	if err != nil {
		panic(err)
	}

	return client
}
