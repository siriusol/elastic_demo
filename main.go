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
	confJson := readConfJson()
	sniff := false
	esConfig := config.Config{
		URL:      confJson.Url,
		Username: confJson.Username,
		Password: confJson.Password,
		Sniff:    &sniff,
	}
	client, err := elastic.NewClientFromConfig(&esConfig)
	if err != nil {
		panic(err)
	}
	info, code, err := client.Ping(confJson.Url).Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ES returned with code %d and version %s\n", code, info.Version.Number)

	esVersion, err := client.ElasticsearchVersion(confJson.Url)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ES version %s\n", esVersion)
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
