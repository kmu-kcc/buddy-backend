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

package oauth2_test

import (
	"testing"

	"github.com/kmu-kcc/buddy-backend/pkg/oauth2"
)

func TestNewToken(t *testing.T) {
	token, exp, err := oauth2.NewToken("20210001")
	if err != nil {
		t.Error(err)
	}

	t.Logf("token: %s\nexpired_at: %d", token, exp)

	if err = token.Valid(); err != nil {
		t.Error(err)
	}

	t.Logf("ID: %s", token.ID())
}
