package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	M_ID         primitive.ObjectID `bson:"_id,omitempty",json:"_id"`
	ID           string             `bson:"id,omitempty",json:"id"`
	Name         string             `bson:"name",json:"name"`
	Credential   string             `bson:"credential",json:"credential"`
	Transactions []Transaction      `bson:"transactions",json:"transactions"`
}

type Transaction struct {
	Amount float64 `bson:"amount",json:"amount"`
	To     string  `bson:"to",json:"to"`
	From   string  `bson:"from",json:"from"`
	Date   string  `bson:"date",json:"date"`
}
