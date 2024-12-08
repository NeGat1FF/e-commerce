package repository_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/NeGat1FF/e-commerce/product-service/internal/models"
	"github.com/NeGat1FF/e-commerce/product-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

// func setupTestDB(t *testing.T) (*mongo.Collection, func()) {
// 	// Set up MongoDB client
// 	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
// 	require.NoError(t, err)

// 	// Get the test collection
// 	collection := client.Database("testdb").Collection("products")

// 	// Clean up function to disconnect and drop the test database
// 	cleanup := func() {
// 		err := client.Database("testdb").Drop(context.Background())
// 		require.NoError(t, err)
// 		client.Disconnect(context.Background())
// 	}

// 	return collection, cleanup
// }

var collection *mongo.Collection

func TestMain(m *testing.M) {
	mongoContainer, err := mongodb.Run(context.Background(), "mongo", mongodb.WithUsername("root"), mongodb.WithPassword("password"))
	if err != nil {
		panic(err)
	}

	str, err := mongoContainer.Endpoint(context.Background(), "")
	if err != nil {
		panic(err)
	}

	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb://root:password@%s", str))
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
	}

	collection = client.Database("testdb").Collection("products")

	code := m.Run()

	err = mongoContainer.Terminate(context.Background())
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestMongoRepository_GetProductsByCategory(t *testing.T) {
	// Clean up the collection
	collection.DeleteMany(context.Background(), bson.M{})

	products := []models.Product{
		{ID: 1, Name: "Product 1", Category: "Category 1"},
		{ID: 2, Name: "Product 2", Category: "Category 2"},
		{ID: 3, Name: "Product 3", Category: "Category 1"},
	}

	var interfaceProducts []interface{}
	for _, p := range products {
		interfaceProducts = append(interfaceProducts, p)
	}

	_, err := collection.InsertMany(context.Background(), interfaceProducts)
	require.NoError(t, err)

	repo := repository.NewMongoRepository(collection)

	page := 1
	limit := 10

	result, err := repo.GetProductsByCategory(context.Background(), "Category 1", page, limit)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	expectedProducts := []models.UserProduct{
		{ID: 1, Name: "Product 1", Category: "Category 1"},
		{ID: 3, Name: "Product 3", Category: "Category 1"},
	}
	assert.Equal(t, expectedProducts, result)

	result, err = repo.GetProductsByCategory(context.Background(), "Category 2", page, limit)
	require.NoError(t, err)
	assert.Len(t, result, 1)

	expectedProducts = []models.UserProduct{
		{ID: 2, Name: "Product 2", Category: "Category 2"},
	}
	assert.Equal(t, expectedProducts, result)

	result, err = repo.GetProductsByCategory(context.Background(), "Category 3", page, limit)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestMongoRepository_GetProductByID(t *testing.T) {
	// Clean up the collection
	collection.DeleteMany(context.Background(), bson.M{})

	products := []models.UserProduct{
		{ID: 1, Name: "Product 1", Category: "Category 1"},
		{ID: 2, Name: "Product 2", Category: "Category 2"},
		{ID: 3, Name: "Product 3", Category: "Category 1"},
	}

	var interfaceProducts []interface{}
	for _, p := range products {
		interfaceProducts = append(interfaceProducts, p)
	}

	_, err := collection.InsertMany(context.Background(), interfaceProducts)
	require.NoError(t, err)

	repo := repository.NewMongoRepository(collection)

	// Test getting a product by ID
	product, err := repo.GetProductByID(context.Background(), 2)
	require.NoError(t, err)
	assert.Equal(t, products[1], product)

	// Test getting a non-existent product
	_, err = repo.GetProductByID(context.Background(), 4)
	require.Error(t, err)
}

func TestMongoRepository_CreateProduct(t *testing.T) {
	// Clean up the collection
	collection.DeleteMany(context.Background(), bson.M{})

	repo := repository.NewMongoRepository(collection)

	// Create a product
	product := models.Product{ID: 1, Name: "Product 1", Category: "Category 1"}
	err := repo.CreateProduct(context.Background(), product)
	require.NoError(t, err)

	// Verify the product
	var result models.Product
	err = collection.FindOne(context.Background(), bson.M{"id": 1}).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, product, result)
}

func TestMongoRepository_UpdateProduct(t *testing.T) {
	// Clean up the collection
	collection.DeleteMany(context.Background(), bson.M{})

	repo := repository.NewMongoRepository(collection)

	product := models.Product{ID: 1, Name: "Product 1", Category: "Category 1"}

	// Insert a test product
	_, err := collection.InsertOne(context.Background(), product)
	require.NoError(t, err)

	// Update the product
	updateFields := map[string]any{"name": "Product 2", "category": "Category 2"}
	err = repo.UpdateProduct(context.Background(), 1, updateFields)
	require.NoError(t, err)

	// Verify the product
	var result models.Product
	err = collection.FindOne(context.Background(), bson.M{"id": 1}).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, result.Name, "Product 2")
	assert.Equal(t, result.Category, "Category 2")
}

func TestMongoRepository_DeleteProduct(t *testing.T) {
	// Clean up the collection
	collection.DeleteMany(context.Background(), bson.M{})

	repo := repository.NewMongoRepository(collection)

	product := models.Product{ID: 1, Name: "Product 1", Category: "Category 1"}

	// Insert a test product
	_, err := collection.InsertOne(context.Background(), product)
	require.NoError(t, err)

	// Delete the product
	err = repo.DeleteProduct(context.Background(), 1)
	require.NoError(t, err)

	// Verify the product is deleted
	count, err := collection.CountDocuments(context.Background(), bson.M{"id": 1})
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestMongoRepository_AddStock(t *testing.T) {
	// Clean up the collection
	collection.DeleteMany(context.Background(), bson.M{})

	repo := repository.NewMongoRepository(collection)

	// Insert a test product
	_, err := collection.InsertOne(context.Background(), bson.M{"id": 1, "quantity": 10})
	require.NoError(t, err)

	// Add stock
	err = repo.AddStock(context.Background(), 1, 5)
	require.NoError(t, err)

	// Verify the stock quantity
	var result bson.M
	err = collection.FindOne(context.Background(), bson.M{"id": 1}).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, int64(15), result["quantity"])
}

func TestMongoRepository_ReduceStock(t *testing.T) {
	// Clean up the collection
	collection.DeleteMany(context.Background(), bson.M{})

	repo := repository.NewMongoRepository(collection)

	// Insert a test product
	_, err := collection.InsertOne(context.Background(), bson.M{"id": 1, "quantity": 10})
	require.NoError(t, err)

	// Reduce stock
	err = repo.ReduceStock(context.Background(), 1, 5)
	require.NoError(t, err)

	// Verify the stock quantity
	var result bson.M
	err = collection.FindOne(context.Background(), bson.M{"id": 1}).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, int64(5), result["quantity"])
}
