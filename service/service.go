package service

import (
	"SecuredMongoDatabase/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type Service struct {
	db *mongo.Database
}

func NewService(db *mongo.Database) *Service {
	return &Service{
		db: db,
	}
}

func (service *Service) GetAllUsers() (users models.Users, err error) {
	log.Println("GET ALL USERS")

	if cursor, err := service.db.Collection(models.UserCollection).Find(context.TODO(), bson.M{}); err != nil {
		return nil, err
	} else {
		err = cursor.All(context.TODO(), &users)
		if err != nil {
			return nil, err
		}
	}

	return
}
