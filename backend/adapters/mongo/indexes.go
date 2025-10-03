package mongoadapters

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureIndexes creates required indexes for reports and audit/debug collections.
func EnsureIndexes(ctx context.Context, client *mongo.Client, dbName string) error {
	db := client.Database(dbName)

	reports := db.Collection("reports")
	_, err := reports.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "kind", Value: 1}, {Key: "target_id", Value: 1}}, Options: options.Index().SetBackground(true)},
		{Keys: bson.D{{Key: "status", Value: 1}, {Key: "created_at", Value: -1}}, Options: options.Index().SetBackground(true)},
		{Keys: bson.D{{Key: "created_at", Value: -1}}, Options: options.Index().SetBackground(true)},
	})
	if err != nil {
		return err
	}

	audit := db.Collection("audit")
	_, err = audit.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "ts", Value: -1}}, Options: options.Index().SetBackground(true)},
		{Keys: bson.D{{Key: "actor_id", Value: 1}}, Options: options.Index().SetBackground(true)},
	})
	if err != nil {
		return err
	}

	// debug collection: keep logs for 30 days
	debug := db.Collection("debug")
	_, err = debug.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "ts", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(60 * 60 * 24 * 30),
	})
	if err != nil {
		return err
	}

	// optionally ensure reports collection has a validation schema (lightweight)
	// Skipped here to keep compatibility; administrators can enable validation separately.

	// allow some time for background index creation to start
	time.Sleep(200 * time.Millisecond)
	return nil
}
