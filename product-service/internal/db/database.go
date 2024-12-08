package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB(dbUrl string) (*mongo.Client, error) {

	return mongo.Connect(context.Background(), options.Client().ApplyURI(dbUrl).SetTimeout(10*time.Second))
}
