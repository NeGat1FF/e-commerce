package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/NeGat1FF/e-commerce/search-service/internal/models"
	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticSearchResponse struct {
	Hits struct {
		Hits []struct {
			Source models.Product `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type ElasticGetResponse struct {
	Source models.Product `json:"_source"`
}

type ElasticClient struct {
	Client    *elasticsearch.Client
	IndexName string
}

func NewElasticClient(cfg elasticsearch.Config, indexName string) (*ElasticClient, error) {
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ElasticClient{
		Client:    es,
		IndexName: indexName,
	}, nil
}

func (ec *ElasticClient) IndexProduct(ctx context.Context, product models.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}
	res, err := ec.Client.Index(ec.IndexName, bytes.NewReader(data),
		ec.Client.Index.WithContext(ctx),
		ec.Client.Index.WithDocumentID(fmt.Sprint(product.ID)))
	if err != nil || res.IsError() {
		return fmt.Errorf("error indexing document: %s", res)
	}
	return nil
}

func (ec *ElasticClient) GetProduct(ctx context.Context, id int64) (models.Product, error) {
	res, err := ec.Client.Get(ec.IndexName, fmt.Sprint(id),
		ec.Client.Get.WithContext(ctx))
	if err != nil || res.IsError() {
		return models.Product{}, fmt.Errorf("error getting document: %s", res)
	}
	defer res.Body.Close()

	var er ElasticGetResponse
	if err := json.NewDecoder(res.Body).Decode(&er); err != nil {
		return models.Product{}, err
	}
	return er.Source, nil
}

func (ec *ElasticClient) ParseQuery(query map[string]string, sortOrder string, min_price, max_price int) *bytes.Buffer {
	// Build the query
	parsedQuery := make([]map[string]interface{}, 0)
	for k, v := range query {
		parsedQuery = append(parsedQuery, map[string]interface{}{
			"match": map[string]interface{}{
				k: v,
			},
		})
	}

	if min_price > 0 {
		parsedQuery = append(parsedQuery, map[string]interface{}{
			"range": map[string]interface{}{
				"price": map[string]interface{}{
					"gte": min_price,
				},
			},
		})
	}

	if max_price > 0 {
		parsedQuery = append(parsedQuery, map[string]interface{}{
			"range": map[string]interface{}{
				"price": map[string]interface{}{
					"lte": max_price,
				},
			},
		})
	}

	q := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": parsedQuery,
			},
		},
	}

	if sortOrder != "" {
		q["sort"] = map[string]interface{}{
			"price": map[string]interface{}{
				"order": sortOrder,
			},
		}
	}

	// Convert the query to JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(q); err != nil {
		return &bytes.Buffer{}
	}

	return &buf
}

func (ec *ElasticClient) SearchProducts(ctx context.Context, query map[string]string, sortOrder string, min_price, max_price, page, limit int) ([]models.Product, error) {
	res, err := ec.Client.Search(
		ec.Client.Search.WithIndex(ec.IndexName),
		ec.Client.Search.WithBody(ec.ParseQuery(query, sortOrder, min_price, max_price)),
		ec.Client.Search.WithContext(ctx),
		ec.Client.Search.WithFrom((page-1)*limit),
		ec.Client.Search.WithSize(limit),
	)

	if err != nil || res.IsError() {
		return nil, fmt.Errorf("error searching document: %s", res)
	}
	defer res.Body.Close()

	var esResponse ElasticSearchResponse
	if err := json.NewDecoder(res.Body).Decode(&esResponse); err != nil {
		return nil, err
	}

	var products []models.Product
	for _, hit := range esResponse.Hits.Hits {
		products = append(products, hit.Source)
	}

	return products, nil
}

func (ec *ElasticClient) UpdateProduct(ctx context.Context, id int64, fields map[string]any) error {
	data, err := json.Marshal(fields)
	if err != nil {
		return err
	}

	res, err := ec.Client.Update(ec.IndexName, fmt.Sprint(id),
		bytes.NewReader([]byte(fmt.Sprintf(`{"doc": %s}`, data))),
		ec.Client.Update.WithContext(ctx))

	if err != nil || res.IsError() {
		return fmt.Errorf("error updating document: %s", res)
	}
	return nil
}

func (ec *ElasticClient) DeleteProduct(ctx context.Context, id int64) error {
	res, err := ec.Client.Delete(ec.IndexName, fmt.Sprint(id),
		ec.Client.Delete.WithContext(ctx))
	if err != nil || res.IsError() {
		return fmt.Errorf("error deleting document: %s", res)
	}
	return nil
}
