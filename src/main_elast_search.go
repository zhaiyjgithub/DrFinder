package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

const ES_NODE_URL = "http://172.16.54.128:9200"

func main()  {
	ctx := context.Background()

	client, err := elastic.NewClient(elastic.SetURL(ES_NODE_URL))
	if err != nil {
		panic(err)
	}

	info, code, err := client.Ping(ES_NODE_URL).Do(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
}
