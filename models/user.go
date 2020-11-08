package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ObjId     primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Id        int                `bson:"id," json:"id"`
	Name      string             `bson:"first_name" json:"first_name"`
	Surname   string             `bson:"last_name" json:"last_name"`
	Idnp      string             `bson:"idnp" json:"idnp"`
	Email     string             `bson:"email" json:"email"`
	IpAddress string             `bson:"ip_address" json:"ip_address"`
}

type Users []User

var UserCollection = "users"
