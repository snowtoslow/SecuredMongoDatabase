package utils

import (
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"log"
)

func ReadJSONFile(file string) bson.D {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("ReadFile error for %v: %v", file, err)
	}

	var fileDoc bson.D
	if err = bson.UnmarshalExtJSON(content, false, &fileDoc); err != nil {
		log.Fatalf("UnmarshalExtJSON error for file %v: %v", file, err)
	}
	return fileDoc
}
