package fee_test

import (
	"testing"

	"github.com/kmu-kcc/buddy-backend/pkg/fee"
)

func TestDones(t *testing.T) {
	if members, err := fee.Dones(2021, 1); err != nil {
		t.Error(err)
	} else {
		t.Log(members)
	}
}

func TestYets(t *testing.T) {

	if members, err := fee.Yets(2021, 1); err != nil {
		t.Error(err)
	} else {
		t.Log(members)
	}
}

func TestAll(t *testing.T) {
	if logs, err := fee.All(2021, 1); err != nil {
		t.Error(err)
	} else {
		t.Log(logs)
	}
}
