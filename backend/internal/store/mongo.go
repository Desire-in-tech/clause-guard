package store

import (
	"context"
	"errors"
	"time"

	"clauseguard/backend/internal/analysis"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoStore struct {
	client     *mongo.Client
	collection *mongo.Collection
	ttl        time.Duration
}

func NewMongoStore(ctx context.Context, uri, database string, ttl time.Duration) (*MongoStore, error) {
	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(10).
		SetConnectTimeout(3 * time.Second).
		SetServerSelectionTimeout(3 * time.Second)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, err
	}

	collection := client.Database(database).Collection("reports")
	indexCtx, cancelIndex := context.WithTimeout(ctx, 3*time.Second)
	defer cancelIndex()
	_, err = collection.Indexes().CreateOne(indexCtx, mongo.IndexModel{
		Keys: bson.D{{Key: "created_at", Value: 1}},
		Options: options.Index().
			SetName("reports_created_at_ttl").
			SetExpireAfterSeconds(int32(ttl.Seconds())),
	})
	if err != nil {
		_ = client.Disconnect(context.Background())
		return nil, err
	}

	return &MongoStore{client: client, collection: collection, ttl: ttl}, nil
}

func (s *MongoStore) SaveReport(ctx context.Context, report *analysis.Report) error {
	id, err := newReportID()
	if err != nil {
		return err
	}
	copy := cloneReport(*report)
	copy.ID = id
	if copy.CreatedAt.IsZero() {
		copy.CreatedAt = time.Now().UTC()
	}
	opCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if _, err := s.collection.InsertOne(opCtx, copy); err != nil {
		return err
	}
	report.ID = id
	report.CreatedAt = copy.CreatedAt
	return nil
}

func (s *MongoStore) GetReport(ctx context.Context, id string) (*analysis.Report, error) {
	opCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	var report analysis.Report
	err := s.collection.FindOne(opCtx, bson.D{{Key: "_id", Value: id}}).Decode(&report)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrReportNotFound
	}
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (s *MongoStore) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}
