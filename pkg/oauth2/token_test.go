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
