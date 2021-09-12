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

package member_test

import (
	"testing"

	"github.com/kmu-kcc/buddy-backend/pkg/member"
)

func TestSignUp(t *testing.T) {
	guests := []*member.Member{
		member.New("20210001", "Test1", "Department1", "010-2021-0001", "testmail1", 1, member.Attending),
		member.New("20190002", "Test2", "Department2", "010-2019-0002", "testmail2", 2, member.Absent),
		member.New("20190003", "Test3", "Department3", "010-2019-0003", "testmail3", 3, member.Attending),
		member.New("20160004", "Test4", "Department2", "010-2016-0004", "testmail4", 4, member.Graduate),
	}

	for _, guest := range guests {
		if err := guest.SignUp(); err != nil {
			t.Error(err)
		}
	}
}

func TestSignUps(t *testing.T) {
	guests, err := member.SignUps()
	if err != nil {
		t.Error(err)
	}

	for _, guest := range guests {
		if guest.Approved {
			t.Error(member.ErrAlreadyMember)
		}
		t.Log(guest)
	}
}

func TestApprove(t *testing.T) {
	ids := []string{"20210001", "20190003"}
	if err := member.Approve(ids); err != nil {
		t.Error(err)
	}
}

func TestSignIn(t *testing.T) {
	memb := member.Member{ID: "20210001", Password: "20210001"}
	guest := member.Member{ID: "20190002", Password: "20190002"}

	if err := memb.SingIn(); err != nil {
		t.Error(err)
	}
	if err := guest.SingIn(); err == member.ErrUnderReview {
		t.Log(err)
	} else if err != nil {
		t.Error(err)
	}
}

func TestExit(t *testing.T) {
	memb := member.Member{ID: "20210001"}
	if err := memb.Exit(); err != nil {
		t.Error(err)
	}
	if err := memb.Exit(); err == member.ErrOnDelete {
		t.Log(err)
	} else if err != nil {
		t.Error(err)
	}

	memb.ID = "20190003"
	if err := memb.Exit(); err != nil {
		t.Error(err)
	}
}

func TestExits(t *testing.T) {
	if membs, err := member.Exits(); err != nil {
		t.Error(err)
	} else {
		for _, memb := range membs {
			t.Log(memb)
		}
	}
}

func TestDelete(t *testing.T) {
	if err := member.Delete([]string{"20190003"}); err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	if err := member.Approve([]string{"20190002"}); err != nil {
		t.Error(err)
	}

	memb := member.Member{ID: "20190002"}
	if err := memb.Update(map[string]interface{}{
		"attendance": member.Attending,
		"password":   "00000000"}); err != nil {
		t.Error(err)
	}
}

func TestSearch(t *testing.T) {
	if membs, err := member.Search("2021"); err != nil {
		t.Error(err)
	} else {
		for _, memb := range membs {
			t.Log(memb)
		}
	}
}

func TestActive(t *testing.T) {
	if active, err := member.Active(); err != nil {
		t.Error(err)
	} else {
		t.Logf("active: %t", active)
	}
}

func TestActivate(t *testing.T) {
	if active, err := member.Activate(true); err != nil {
		t.Error(err)
	} else {
		t.Logf("active: %t", active)
	}
}

func TestGraduates(t *testing.T) {
	members, err := member.Graduates()
	if err != nil {
		t.Error(err)
	}
	for _, memb := range members {
		t.Log(memb)
	}
}
