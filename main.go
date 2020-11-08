package main

import (
	"SecuredMongoDatabase/mongo"
	"SecuredMongoDatabase/service"
	"log"
)

func main() {

	mongoConfig := mongo.NewMongoConfig("27017", "mongodb://localhost:", "usersdb")

	database := mongoConfig.InitDb()

	mongoService := service.NewService(database)

	log.Println(mongoService.GetAllUsers())

}
