package models

import (
	"errors"
	"strings"	
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name,omitempty" validate:"required,name"`
	Participants []string		`json:"participants" bson:"participants"`
	Password string             `json:"password" bson:"password,omitempty"`
}

func (r *Room) BeforeSave() error {
	hashedPassword, err := Hash(r.Password)
	if err != nil {
		return err
	}
	r.Password = string(hashedPassword)
	return nil
}

func (r *Room) Prepare() {
	r.ID = primitive.NewObjectID();
	r.Participants = []string{}
}

func (r *Room) Validate(action string) error {
	switch strings.ToLower(action) {
	case "create":
		if r.Password == "" {
			return errors.New("Required Password")
		}
		if r.Name == "" {
			return errors.New("Required Name")
		}
		return nil

	default:
		if r.Password == "" {
			return errors.New("Required Password")
		}
		if r.Name == "" {
			return errors.New("Required Name")
		}
		return nil
	}
}
