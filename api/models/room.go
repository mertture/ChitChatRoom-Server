package models

import (
	"errors"
	"strings"	
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name,omitempty" validate:"required,name"`
	Messages []Message			`json:"messages" bson:"messages"`
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


type NewParticipantJoinedResponse struct {
    Message     string `json:"message" bson:"message"`
    Participant User `json:"participant" bson:"participant"`
}

type ParticipantLeftResponse struct {
    Message     string `json:"message" bson:"message"`
    Participant string `json:"participant" bson:"participant"`
}

type UsersResponse struct {
	Action			string	`json:"action" bson:"action"`
	Participants	[]User 	`json:"participants" bson:"participants"`
}

type MessagesResponse struct {
	Action			string	    `json:"action" bson:"action"`
	Messages	[]Message	`json:"messages" bson:"messages"`

}