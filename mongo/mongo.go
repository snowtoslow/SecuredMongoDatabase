package mongo

import (
	"SecuredMongoDatabase/utils"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type mongoConfig struct {
	port         string
	host         string
	dbName       string
	kmsProviders map[string]map[string]interface{}
}

func NewMongoConfig(port string, host string,
	dbName string, kmsProviders map[string]map[string]interface{}) *mongoConfig {
	return &mongoConfig{
		port:         port,
		host:         host,
		dbName:       dbName,
		kmsProviders: kmsProviders,
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

func (mongoConfig *mongoConfig) CreateDataKey() {
	kvClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoConfig.host+mongoConfig.port))
	if err != nil {
		log.Printf("Error connect to database %v: ", err)
	}

	// create key vault client and drop key vault collection
	_ = kvClient.Database("keyvault").Collection("__datakeys").Drop(context.Background())

	// create ClientEncryption
	clientEncryptOptions := options.ClientEncryption().SetKeyVaultNamespace("keyvault.__datakeys").
		SetKmsProviders(mongoConfig.kmsProviders)

	clientEncryption, err := mongo.NewClientEncryption(kvClient, clientEncryptOptions)
	if err != nil {
		log.Printf("Error during create client encryption: %v", err)
	}

	defer clientEncryption.Close(context.Background())

	// create a new data key
	dataKeyOptions := options.DataKey().SetKeyAltNames(utils.KeysArray)

	_, err = clientEncryption.CreateDataKey(context.Background(), "local", dataKeyOptions)
	if err != nil {
		log.Printf("Erro creating a new datakey: %v", err)
	}
}

func (mongoConfig *mongoConfig) CreateEncryptedClient() *mongo.Client {
	//create a client with auto encryption

	schemaMap := map[string]interface{}{
		"usersdb.users": utils.ReadJSONFile("my-magic-collection.json"),
	}

	autoEncryptionOptions := options.AutoEncryption().
		SetKeyVaultNamespace("keyvault.__datakeys").
		SetKmsProviders(mongoConfig.kmsProviders).
		SetSchemaMap(schemaMap)

	clientOptions := options.Client().ApplyURI(mongoConfig.host + mongoConfig.port).
		SetAutoEncryptionOptions(autoEncryptionOptions)

	autoEncryptionClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Printf("Error occured connecting after create autoEncClient: %v", err)
	}

	return autoEncryptionClient
}
