package oauth_test

import (
	"testing"

	"github.com/kmu-kcc/buddy-backend/pkg/oauth"
)

func TestNewToken(t *testing.T) {
	token, exp, err := oauth.NewToken("20210001")
	if err != nil {
		t.Error(err)
	}

	t.Logf("token: %s\texpired_at: %d", token, exp)

	if err = token.Verify(); err != nil {
		t.Error(err)
	}

	t.Logf("ID: %s", token.ID())
}
