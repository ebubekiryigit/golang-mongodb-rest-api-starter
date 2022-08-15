package models

import (
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type Token struct {
	mgm.DefaultModel `bson:",inline"`
	User             primitive.ObjectID `json:"user" bson:"user"`
	Token            string             `json:"token" bson:"token"`
	Type             string             `json:"type" bson:"type"`
	ExpiresAt        time.Time          `json:"expires_at" bson:"expires_at"`
	Blacklisted      bool               `json:"blacklisted" bson:"blacklisted"`
}

func (model *Token) GetResponseJson() gin.H {
	return gin.H{"token": model.Token, "expires": model.ExpiresAt.Format("2006-01-02 15:04:05")}
}

func NewToken(userId primitive.ObjectID, tokenString string, tokenType string, expiresAt time.Time) *Token {
	return &Token{
		User:        userId,
		Token:       tokenString,
		Type:        tokenType,
		ExpiresAt:   expiresAt,
		Blacklisted: false,
	}
}

func (model *Token) CollectionName() string {
	return "tokens"
}

// You can override Collection functions or CRUD hooks
// https://github.com/Kamva/mgm#a-models-hooks
// https://github.com/Kamva/mgm#collections
