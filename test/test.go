package main

import (
	"SecuredMongoDatabase/utils"
	"context"
	"encoding/base64"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ExampleExplictencryption()
}

func ExampleExplictencryption() {
	//var localMasterKey []byte // This must be the same master key that was used to create the encryption key.

	decodeKey, err := base64.StdEncoding.DecodeString(utils.LocalMasterKey)

	kmsProviders := map[string]map[string]interface{}{
		"local": {
			"key": decodeKey,
		},
	}

	// The MongoDB namespace (db.collection) used to store the encryption data keys.
	keyVaultDBName, keyVaultCollName := "encryption", "testKeyVault"
	keyVaultNamespace := keyVaultDBName + "." + keyVaultCollName

	// The Client used to read/write application data.
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Disconnect(context.TODO()) }()

	// Get a handle to the application collection and clear existing data.
	coll := client.Database("test").Collection("coll")
	_ = coll.Drop(context.TODO())

	// Set up the key vault for this example.
	keyVaultColl := client.Database(keyVaultDBName).Collection(keyVaultCollName)
	_ = keyVaultColl.Drop(context.TODO())
	// Ensure that two data keys cannot share the same keyAltName.
	keyVaultIndex := mongo.IndexModel{
		Keys: bson.D{{"keyAltNames", 1}},
		Options: options.Index().
			SetUnique(true).
			SetPartialFilterExpression(bson.D{
				{"keyAltNames", bson.D{
					{"$exists", true},
				}},
			}),
	}
	if _, err = keyVaultColl.Indexes().CreateOne(context.TODO(), keyVaultIndex); err != nil {
		panic(err)
	}

	// Create the ClientEncryption object to use for explicit encryption/decryption. The Client passed to
	// NewClientEncryption is used to read/write to the key vault. This can be the same Client used by the main
	// application.
	clientEncryptionOpts := options.ClientEncryption().
		SetKmsProviders(kmsProviders).
		SetKeyVaultNamespace(keyVaultNamespace)
	clientEncryption, err := mongo.NewClientEncryption(client, clientEncryptionOpts)
	if err != nil {
		panic(err)
	}
	defer func() { _ = clientEncryption.Close(context.TODO()) }()

	// Create a new data key for the encrypted field.
	dataKeyOpts := options.DataKey().SetKeyAltNames([]string{"go_encryption_example"})
	dataKeyID, err := clientEncryption.CreateDataKey(context.TODO(), "local", dataKeyOpts)
	if err != nil {
		panic(err)
	}

	// Create a bson.RawValue to encrypt and encrypt it using the key that was just created.
	rawValueType, rawValueData, err := bson.MarshalValue("123456789")
	if err != nil {
		panic(err)
	}
	rawValue := bson.RawValue{Type: rawValueType, Value: rawValueData}
	encryptionOpts := options.Encrypt().
		SetAlgorithm("AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic").
		SetKeyID(dataKeyID)
	encryptedField, err := clientEncryption.Encrypt(context.TODO(), rawValue, encryptionOpts)
	if err != nil {
		panic(err)
	}

	// Insert a document with the encrypted field and then find it.
	if _, err = coll.InsertOne(context.TODO(), bson.D{{"encryptedField", encryptedField}}); err != nil {
		panic(err)
	}
	var foundDoc bson.M
	if err = coll.FindOne(context.TODO(), bson.D{}).Decode(&foundDoc); err != nil {
		panic(err)
	}

	// Decrypt the encrypted field in the found document.
	decrypted, err := clientEncryption.Decrypt(context.TODO(), foundDoc["encryptedField"].(primitive.Binary))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decrypted value: %s\n", decrypted)
}
