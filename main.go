package main

import (
	"SecuredMongoDatabase/mongo"
	"SecuredMongoDatabase/utils"
	"context"
	"encoding/base64"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

func main() {

	myMagicKmsProvider := make(map[string]map[string]interface{})

	decodeKey, err := base64.StdEncoding.DecodeString(utils.LocalMasterKey)

	if err != nil {
		log.Printf("Error decoding localMatesKey: %v", err)
	}

	myMagicKmsProvider = map[string]map[string]interface{}{
		"local": {"key": decodeKey},
	}

	mongoConfig := mongo.NewMongoConfig("27017", "mongodb://localhost:", "usersdb", myMagicKmsProvider)

	mongoConfig.CreateDataKey()

	client := mongoConfig.CreateEncryptedClient()

	if err = client.Disconnect(context.Background()); err != nil {
		log.Printf("Error disconecting client: %v", err)
	}

	collection := client.Database("usersdb").Collection("users")

	res, err := collection.FindOne(context.Background(), bson.D{}).DecodeBytes()
	if err != nil {
		log.Fatalf("FindOne error: %v", err)
	}
	fmt.Println(res)

}
