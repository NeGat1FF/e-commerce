package elasticsearch_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	es "github.com/NeGat1FF/e-commerce/search-service/internal/elasticsearch"
	"github.com/NeGat1FF/e-commerce/search-service/internal/models"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tc "github.com/testcontainers/testcontainers-go/modules/elasticsearch"
)

var ec *es.ElasticClient

func TestMain(m *testing.M) {

	elasticsearchContainer, err := tc.Run(context.Background(), "docker.elastic.co/elasticsearch/elasticsearch:8.16.0", tc.WithPassword("password"))
	if err != nil {
		panic(err)
	}

	ec, err = es.NewElasticClient(elasticsearch.Config{
		Addresses: []string{elasticsearchContainer.Settings.Address},
		Username:  elasticsearchContainer.Settings.Username,
		Password:  elasticsearchContainer.Settings.Password,
		CACert:    elasticsearchContainer.Settings.CACert,
	}, "test")
	if err != nil {
		panic(err)
	}

	// setup()
	code := m.Run()
	// teardown()
	os.Exit(code)
}

func TestElasticClient_AddProduct(t *testing.T) {
	// Clear the index
	ec.Client.Indices.Delete([]string{"test"})

	// Create the index
	ec.Client.Indices.Create("test")

	product := models.Product{
		ID:          1,
		Name:        "Product 1",
		Price:       100,
		Description: "Description",
		Attributes: map[string]any{
			"key1": "value1",
			"key2": "value2",
		},
	}

	err := ec.IndexProduct(context.Background(), product)
	require.NoError(t, err)

	_, err = ec.Client.Indices.Refresh(ec.Client.Indices.Refresh.WithIndex("test"))
	require.NoError(t, err)

	res, err := ec.Client.Search(
		ec.Client.Search.WithIndex("test"),
		ec.Client.Search.WithQuery("Product 1"),
	)
	require.NoError(t, err)
	assert.Equal(t, false, res.IsError())

	var esResponse es.ElasticSearchResponse
	err = json.NewDecoder(res.Body).Decode(&esResponse)

	require.NoError(t, err)
	assert.Equal(t, 1, len(esResponse.Hits.Hits))
}

func TestElasticClient_GetProduct(t *testing.T) {
	// Clear the index
	ec.Client.Indices.Delete([]string{"test"})

	// Create the index
	ec.Client.Indices.Create("test")

	product := models.Product{
		ID:          1,
		Name:        "Product 1",
		Price:       100,
		Description: "Description",
		Attributes: map[string]any{
			"key1": "value1",
			"key2": "value2",
		},
	}

	data, err := json.Marshal(product)
	require.NoError(t, err)

	res, err := ec.Client.Index("test", bytes.NewReader(data), ec.Client.Index.WithDocumentID(fmt.Sprint(product.ID)))
	require.NoError(t, err)
	assert.Equal(t, false, res.IsError())

	_, err = ec.Client.Indices.Refresh(ec.Client.Indices.Refresh.WithIndex("test"))
	require.NoError(t, err)

	product, err = ec.GetProduct(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, "Product 1", product.Name)
}

func TestElasticClient_SearchProducts(t *testing.T) {
	// Clear the index
	ec.Client.Indices.Delete([]string{"test"})

	// Create the index
	ec.Client.Indices.Create("test")

	products := []models.Product{
		{
			ID:          1,
			Name:        "Product 2",
			Price:       100,
			Description: "Description",
			Attributes: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			ID:          2,
			Name:        "Product 1",
			Price:       200,
			Description: "Description",
			Attributes: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	data, err := json.Marshal(products[0])
	require.NoError(t, err)

	res, err := ec.Client.Index("test", bytes.NewReader(data))
	require.NoError(t, err)
	assert.Equal(t, false, res.IsError())

	data, err = json.Marshal(products[1])
	require.NoError(t, err)

	res, err = ec.Client.Index("test", bytes.NewReader(data))
	require.NoError(t, err)
	assert.Equal(t, false, res.IsError())

	_, err = ec.Client.Indices.Refresh(ec.Client.Indices.Refresh.WithIndex("test"))
	require.NoError(t, err)

	query := map[string]string{
		"name.keyword": "Product 2",
	}

	products, err = ec.SearchProducts(context.Background(), query, "", 0, 0, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, 1, len(products))
	assert.Equal(t, "Product 2", products[0].Name)
}

func TestElasticClient_UpdateProduct(t *testing.T) {
	// Clear the index
	ec.Client.Indices.Delete([]string{"test"})

	// Create the index
	ec.Client.Indices.Create("test")

	product := models.Product{
		ID:          1,
		Name:        "Product 1",
		Price:       100,
		Description: "Description",
		Attributes: map[string]any{
			"key1": "value1",
			"key2": "value2",
		},
	}

	err := ec.IndexProduct(context.Background(), product)
	require.NoError(t, err)

	_, err = ec.Client.Indices.Refresh(ec.Client.Indices.Refresh.WithIndex("test"))
	require.NoError(t, err)

	updateFields := map[string]any{
		"price": 250,
		"name":  "Product 2",
	}
	err = ec.UpdateProduct(context.Background(), 1, updateFields)
	require.NoError(t, err)

	_, err = ec.Client.Indices.Refresh(ec.Client.Indices.Refresh.WithIndex("test"))
	require.NoError(t, err)

	res, err := ec.Client.Search(
		ec.Client.Search.WithIndex("test"),
		ec.Client.Search.WithQuery("Product 2"),
	)
	require.NoError(t, err)
	assert.Equal(t, false, res.IsError())

	var esResponse es.ElasticSearchResponse
	err = json.NewDecoder(res.Body).Decode(&esResponse)
	require.NoError(t, err)
	assert.Equal(t, 1, len(esResponse.Hits.Hits))
	assert.Equal(t, "Product 2", esResponse.Hits.Hits[0].Source.Name)
	assert.Equal(t, float64(250), esResponse.Hits.Hits[0].Source.Price)
}
func TestElasticClient_DeleteProduct(t *testing.T) {
	product := models.Product{
		ID:          1,
		Name:        "Product 1",
		Price:       100,
		Description: "Description",
		Attributes: map[string]any{
			"key1": "value1",
			"key2": "value2",
		},
	}

	err := ec.IndexProduct(context.Background(), product)
	require.NoError(t, err)

	_, err = ec.Client.Indices.Refresh(ec.Client.Indices.Refresh.WithIndex("test"))
	require.NoError(t, err)

	err = ec.DeleteProduct(context.Background(), 1)
	require.NoError(t, err)

	_, err = ec.Client.Indices.Refresh(ec.Client.Indices.Refresh.WithIndex("test"))
	require.NoError(t, err)

	res, err := ec.Client.Get("test", fmt.Sprint(product.ID))
	require.NoError(t, err)
	assert.Equal(t, true, res.IsError())
}
