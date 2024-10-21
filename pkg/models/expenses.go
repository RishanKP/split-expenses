package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	SPLIT_TYPE_PERCENTAGE = "PERCENTAGE"
	SPLIT_TYPE_EQUAL      = "EQUAL"
	SPLIT_TYPE_EXACT      = "EXACT"
)

type Participant struct {
	UserID     string  `bson:"user_id" json:"user_id"`
	Amount     float64 `bson:"amount,omitempty" json:"amount,omitempty"`
	Percentage float64 `bson:"percentage,omitempty" json:"percentage,omitempty"`
}

type Expense struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Amount       float64            `bson:"amount" json:"amount"`
	Description  string             `bson:"description" json:"description"`
	SplitType    string             `bson:"splitType" json:"splitType"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	Participants []Participant      `bson:"participants" json:"participants"`
}
