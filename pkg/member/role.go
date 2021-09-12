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

// Package member provides access to the club member of the Buddy System.
package member

// Role represents the member role.
type Role struct {
	Master             bool `json:"-" bson:"master"`
	MemberManagement   bool `json:"member_management" bson:"member_management"`
	ActivityManagement bool `json:"activity_management" bson:"activity_management"`
	FeeManagement      bool `json:"fee_management" bson:"fee_management"`
}

// NewRole returns a new role without any authorities.
func NewRole() *Role { return &Role{} }
