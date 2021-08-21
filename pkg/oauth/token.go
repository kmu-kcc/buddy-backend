// Package oauth provides OAuth 2.0 verification.
package oauth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kmu-kcc/buddy-backend/config"
	"github.com/kmu-kcc/buddy-backend/pkg/member"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Token represents an access token.
type Token string

var (
	tokens = make(map[Token]map[string]interface{})

	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

// NewToken generates an access token.
func NewToken(id string) (Token, int64, error) {
	atClaims := make(jwt.MapClaims)
	exp := time.Now().Add(6 * time.Hour).Unix()

	atClaims["authorized"] = true
	atClaims["id"] = id
	atClaims["expire"] = exp

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", 0, err
	}
	tokens[Token(token)] = make(map[string]interface{})
	tokens[Token(token)]["id"] = id
	tokens[Token(token)]["exp"] = exp

	time.AfterFunc(6*time.Hour, func() { delete(tokens, Token(token)) })

	return Token(token), exp, err
}

// Verify reports whether t is valid or not.
func (t Token) Verify() error {
	if meta, ok := tokens[t]; !ok {
		return ErrInvalidToken
	} else if time.Unix(meta["exp"].(int64), 0).Before(time.Now()) {
		delete(tokens, t)
		return ErrExpiredToken
	}
	return nil
}

// ID returns the id of t.
func (t Token) ID() string {
	meta, ok := tokens[t]
	if !ok {
		return ""
	}
	return meta["id"].(string)
}

// Role returns the role corresponding to t.
func (t Token) Role() (map[int]bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return nil, err
	}

	memb := new(member.Member)

	if err = client.Database("club").
		Collection("members").
		FindOne(
			ctx,
			bson.D{bson.E{Key: "id", Value: t.ID()}}).
		Decode(memb); err != nil {
		return nil, err
	}

	return memb.Role, client.Disconnect(ctx)
}
