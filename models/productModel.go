package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id           primitive.ObjectID `bson:"_id" json:"id"`
	Title        string             `json:"title" validate:"required"`
	Description  string             `json:"description" validate:"required"`
	CategoryName string             `bson:"categoryName" json:"categoryName" validate:"required"`
	CreatedOn    time.Time          `bson:"createdOn" json:"createdOn"`
	UpdatedOn    time.Time          `bson:"updatedOn" json:"updatedOn"`
}
