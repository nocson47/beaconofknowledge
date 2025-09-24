package mongoadapters

import (
	"context"
	"fmt"
	"time"

	"github.com/nocson47/beaconofknowledge/internal/entities"
	"github.com/nocson47/beaconofknowledge/internal/repositories"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoReportRepo struct {
	col *mongo.Collection
}

func NewMongoReportRepo(client *mongo.Client, dbName string) repositories.ReportRepository {
	col := client.Database(dbName).Collection("reports")
	return &MongoReportRepo{col: col}
}

func (m *MongoReportRepo) CreateReport(ctx context.Context, r *entities.Report) (string, error) {
	// store as BSON; use generated ObjectID but return a simple timestamp-based int id for compatibility
	doc := bson.M{
		"reporter_id": r.ReporterID,
		"kind":        r.Kind,
		"target_id":   r.TargetID,
		"reason":      r.Reason,
		"status":      r.Status,
		"created_at":  time.Now(),
	}
	res, err := m.col.InsertOne(ctx, doc)
	if err != nil {
		return "", fmt.Errorf("mongo insert: %w", err)
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}
	return "", nil
}

func (m *MongoReportRepo) GetReports(ctx context.Context, kind *string) ([]*entities.Report, error) {
	filter := bson.M{}
	if kind != nil {
		filter["kind"] = *kind
	}
	cur, err := m.col.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("mongo find: %w", err)
	}
	defer cur.Close(ctx)
	var out []*entities.Report
	for cur.Next(ctx) {
		var doc bson.M
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		rep := &entities.Report{}
		if idv, ok := doc["_id"].(primitive.ObjectID); ok {
			rep.ID = idv.Hex()
		}
		if v, ok := doc["reporter_id"].(int32); ok {
			vi := int(v)
			rep.ReporterID = &vi
		}
		if v, ok := doc["kind"].(string); ok {
			rep.Kind = v
		}
		if v, ok := doc["target_id"].(int32); ok {
			rep.TargetID = int(v)
		}
		if v, ok := doc["reason"].(string); ok {
			rep.Reason = v
		}
		if v, ok := doc["status"].(string); ok {
			rep.Status = v
		}
		if v, ok := doc["created_at"].(primitive.DateTime); ok {
			t := v.Time()
			rep.CreatedAt = t
		}
		out = append(out, rep)
	}
	return out, nil
}

func (m *MongoReportRepo) UpdateReportStatus(ctx context.Context, id string, status string, resolvedBy *int) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	upd := bson.M{"$set": bson.M{"status": status}}
	if resolvedBy != nil {
		upd["$set"].(bson.M)["resolved_by"] = *resolvedBy
		upd["$set"].(bson.M)["resolved_at"] = time.Now()
	}
	_, err = m.col.UpdateByID(ctx, oid, upd)
	if err != nil {
		return fmt.Errorf("mongo update: %w", err)
	}
	return nil
}
