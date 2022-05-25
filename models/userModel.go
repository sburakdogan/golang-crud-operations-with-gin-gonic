package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	FullName  string             `bson:"fullName" json:"fullName" validate:"required"`
	Email     string             `json:"email" validate:"required,email"`
	Password  string             `json:"password" validate:"required,min=6"`
	CreatedOn time.Time          `bson:"createdOn" json:"createdOn"`
}
