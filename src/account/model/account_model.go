package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	M_ID         primitive.ObjectID `bson:"_id,omitempty",json:"_id"`
	ID           string             `bson:"id,omitempty",json:"id"`
	Name         string             `bson:"name",json:"name"`
	UserCode     string             `bson:"userCode",json:"userCode"`
	Transactions []Transaction      `bson:"transactions",json:"transactions"`
}

type Transaction struct {
	Amount float64 `bson:"amount",json:"amount"`
	To     string  `bson:"to",json:"to"`
	From   string  `bson:"from",json:"from"`
}

type RegisterBody struct {
	Username string  `json:"username",form:"username"`
	Cash     float64 `json:"cash",form:"cash"`
}
