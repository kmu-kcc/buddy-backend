// Copyright 2021 KMU KCC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package oauth2 provides OAuth 2.0 verification.
package oauth2

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
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
)

// NewToken generates an access token.
func NewToken(id string) (Token, int64, error) {
	atClaims := make(jwt.MapClaims)
	exp := time.Now().Add(6 * time.Hour).Unix()

	atClaims["authorized"] = true
	atClaims["id"] = id
	atClaims["expire"] = exp

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(config.AccessSecret))
	if err != nil {
		return "", 0, err
	}
	tokens[Token(token)] = make(map[string]interface{})
	tokens[Token(token)]["id"] = id
	tokens[Token(token)]["exp"] = exp

	time.AfterFunc(6*time.Hour, func() { delete(tokens, Token(token)) })

	return Token(token), exp, err
}

// Valid reports whether t is valid or not.
func (t Token) Valid() error {
	if _, exists := tokens[t]; !exists {
		return ErrInvalidToken
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
func (t Token) Role() (member.Role, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return member.Role{}, err
	}
	defer client.Disconnect(ctx)

	memb := new(member.Member)

	if err = client.Database("club").Collection("members").FindOne(ctx, bson.D{bson.E{Key: "id", Value: t.ID()}}).Decode(memb); err != nil {
		return member.Role{}, err
	}
	return memb.Role, nil
}
