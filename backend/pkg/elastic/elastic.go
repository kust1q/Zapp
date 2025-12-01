package elastic

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticClient(addresses []string) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating the client: %w", err)
	}

	res, err := es.Info()
	if err != nil {
		return nil, fmt.Errorf("error getting info response: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error: %s", res.String())
	}

	return es, nil
}
