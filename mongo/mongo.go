package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type mongoConfig struct {
	port   string
	host   string
	dbName string
}

func NewMongoConfig(port string, host string, dbName string) *mongoConfig {
	return &mongoConfig{
		port:   port,
		host:   host,
		dbName: dbName,
	}
}

func (mongoConfig *mongoConfig) InitDb() *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoConfig.host + mongoConfig.port))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	//defer client.Disconnect(ctx)

	return client.Database(mongoConfig.dbName)
}
