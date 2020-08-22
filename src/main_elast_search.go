package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

const EsNodeUrl = "http://172.16.54.128:9200"

func main()  {
	ctx := context.Background()

	client, err := elastic.NewClient(elastic.SetURL(EsNodeUrl))
	if err != nil {
		panic(err)
	}

	info, code, err := client.Ping(EsNodeUrl).Do(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
}
