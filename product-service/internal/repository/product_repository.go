package repository

import (
	"context"
	"errors"

	"github.com/NeGat1FF/e-commerce/product-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrProductAlreadyExists = errors.New("product already exists")
var ErrProductNotFound = errors.New("product not found")

type MongoRepository struct {
	coll *mongo.Collection
}

func NewMongoRepository(collection *mongo.Collection) *MongoRepository {
	return &MongoRepository{
		coll: collection,
	}
}

func (r *MongoRepository) GetProductsByCategory(ctx context.Context, category string, page, limit int) ([]models.UserProduct, error) {
	filter := bson.M{"category": category}

	opts := options.Find()
	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))

	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var products []models.UserProduct
	err = cur.All(ctx, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *MongoRepository) GetProductByID(ctx context.Context, id int64) (models.UserProduct, error) {
	var product models.UserProduct

	filter := bson.M{"id": id}

	res := r.coll.FindOne(ctx, filter, nil)
	err := res.Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return product, ErrProductNotFound
		}
		return product, err
	}
	return product, nil
}

func (r *MongoRepository) CreateProduct(ctx context.Context, product models.Product) error {
	// Check if the product already exists
	_, err := r.GetProductByID(ctx, product.ID)
	if err == nil {
		return ErrProductAlreadyExists
	}
	_, err = r.coll.InsertOne(ctx, product)
	return err
}

func (r *MongoRepository) UpdateProduct(ctx context.Context, id int64, updateFields map[string]any) error {
	res, err := r.coll.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": updateFields})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (r *MongoRepository) DeleteProduct(ctx context.Context, id int64) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrProductNotFound
	}
	return err
}

func (r *MongoRepository) GetStock(ctx context.Context, id int64) (int64, error) {
	var product models.Product
	err := r.coll.FindOne(ctx, bson.M{"id": id}).Decode(&product)
	if err != nil {
		return 0, err
	}

	return product.Quantity, nil
}

func (r *MongoRepository) AddStock(ctx context.Context, id int64, quantity int64) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$inc": bson.M{"quantity": quantity}})
	return err
}

func (r *MongoRepository) ReduceStock(ctx context.Context, id int64, quantity int64) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$inc": bson.M{"quantity": -quantity}})
	return err
}
