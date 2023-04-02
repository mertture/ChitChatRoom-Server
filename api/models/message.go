package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

)

type Message struct {
	ID       	primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Sender		Sender			   `json:"sender" bson:"sender,omitempty"`
	Content     string             `json:"content" bson:"content,omitempty"`
	CreatedAt   float64			   `json:"createdAt" bson:"createdAt,omitempty"`
}

type Sender struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email    string             `json:"email" bson:"email,omitempty" validate:"required,email"`
}

func (m *Message) Prepare() {
	m.ID = primitive.NewObjectID();
}