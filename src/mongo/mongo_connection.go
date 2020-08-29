package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Host     string `env:"MONGO_HOST,required"`
	User     string `env:"MONGO_USER,required"`
	Pass     string `env:"MONGO_PASS,required"`
	Database string `env:"MONGO_DATABASE,required"`
}

// Connect connects to a mongo db instance via url and env var credentials
func (m MongoConfig) Connect(collectionName string) (*mongo.Collection, error) {

	// Set client options
	clientOptions := options.
		Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s/%s", m.User, m.Pass, m.Host, m.Database))

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// Connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	collection := client.Database(m.Database).Collection(collectionName)
	return collection, nil
}

type IMongoClient interface {
	FindOne(ctx context.Context, filter interface{}, data interface{}) error
	InsertOne(ctx context.Context, data interface{}) (interface{}, error)
}

type MongoClient struct {
	Conn *mongo.Collection
}

// FindOne wraps mongo driver method FindOne
func (m MongoClient) FindOne(ctx context.Context, filter interface{}, data interface{}) error {
	return m.Conn.FindOne(ctx, filter).Decode(data)
}

// InsertOne wraps mongo driver method InsertOne
func (m MongoClient) InsertOne(ctx context.Context, data interface{}) (interface{}, error) {
	res, err := m.Conn.InsertOne(ctx, data)
	return res.InsertedID, err
}
