package mongo

import (
	"context"
	"errors"
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
func (m MongoConfig) Connect(collectionName string) (*mongo.Client, *mongo.Collection, error) {

	// Set client options
	clientOptions := options.
		Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s/%s", m.User, m.Pass, m.Host, m.Database))

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// Connect to MongoDB
	err = client.Connect(ctx)
	if err != nil {
		return nil, nil, err
	}

	log.Println("Connected to MongoDB!")
	collection := client.Database(m.Database).Collection(collectionName)
	return client, collection, nil
}

type IMongoClient interface {
	FindOne(ctx context.Context, filter interface{}, data interface{}) error
	InsertOne(ctx context.Context, data interface{}) (interface{}, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) error

	StartTransaction(ctx context.Context, callback func() error) error
}

type MongoClient struct {
	Client *mongo.Client
	Conn   *mongo.Collection
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

func (m MongoClient) UpdateOne(ctx context.Context, filter interface{}, data interface{}) error {
	res, err := m.Conn.UpdateOne(ctx, filter, data)
	if err != nil {
		return err
	}
	if res.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
		return nil
	} else if res.UpsertedCount != 0 {
		fmt.Printf("inserted a new document with ID %v\n", res.UpsertedID)
		return nil
	} else {
		return errors.New("No update occured")
	}
}

func (m MongoClient) StartTransaction(ctx context.Context, callback func() error) error {
	session, err := m.Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	if err := session.StartTransaction(); err != nil {
		return err
	}

	return mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := callback(); err != nil {
			return err
		}
		if err = session.CommitTransaction(ctx); err != nil {
			return err
		}
		return nil
	})
}
