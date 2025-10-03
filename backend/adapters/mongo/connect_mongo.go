package mongoadapters

import (
	"context"
	"fmt"
	"time"

	"github.com/nocson47/beaconofknowledge/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(cfg *config.Configuration) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%d", cfg.MongoHost, cfg.MongoPort)
	clientOpts := options.Client().ApplyURI(uri)
	if cfg.MongoUser != "" {
		clientOpts.SetAuth(options.Credential{Username: cfg.MongoUser, Password: cfg.MongoPassword})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client, nil
}
