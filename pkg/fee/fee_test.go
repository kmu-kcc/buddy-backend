package fee_test

import (
	"context"
	"testing"
	"time"

	"github.com/kmu-kcc/buddy-backend/config"
	"github.com/kmu-kcc/buddy-backend/pkg/fee"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreate(t *testing.T) {
	fees := []*fee.Fee{
		fee.New(2022, 1, 0, 30000),
		fee.New(2022, 2, 0, 40000),
	}

	for _, fee := range fees {
		if err := fee.Create(); err != nil {
			t.Error(err)
		}
	}
}

func TestAmount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	log1 := fee.NewLog("abc", "회비 납부", 20000, 0)
	log2 := fee.NewLog("abc", "회비 납부", 30000, 0)

	testFee := new(fee.Fee)
	testFee.Year = 2022
	testFee.Semester = 2
	testFee.Amount = 30000
	testFee.Logs = []primitive.ObjectID{log1.ID, log2.ID}

	if _, err := client.Database("club").Collection("logs").InsertMany(ctx, []interface{}{log1, log2}); err != nil {
		t.Error(err)
	}
	if _, err := client.Database("club").Collection("fees").InsertOne(ctx, testFee); err != nil {
		t.Error(err)
	}

	sum, err := fee.Amount(2022, 2, "abc")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(sum)
	}
}

func TestPayers(t *testing.T) {
	f := fee.Fee{Year: 2021, Semester: 1}
	if members, err := f.Payers(); err != nil {
		t.Error(err)
	} else {
		t.Log(members)
	}
}

func TestDeptors1(t *testing.T) {
	f := fee.Fee{Year: 2021, Semester: 1}
	if members, a, err := f.Deptors(); err != nil {
		t.Error(err)
	} else {
		t.Log(members, a)
	}
}

func TestSearch(t *testing.T) {
	f := fee.Fee{Year: 2021, Semester: 1}
	if a, logs, b, err := f.Search(); err != nil {
		t.Error(err)
	} else {
		t.Log(a, logs, b)
	}
}

func TestPay(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		t.Fatal(err)
	}

	testLog := fee.NewLog("20181681", "회비 납부", 0, 0)
	testLog2 := fee.NewLog("20181682", "회비 납부", 0, 0)

	if err := fee.Pay(2021, 4, []string{testLog.MemberID, testLog2.MemberID}, []int{10000, 1000}); err != nil {
		t.Fatal(err)
	}

	if err = client.Disconnect(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestDeposit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))

	if err != nil {
		t.Fatal(err)
	}

	if err := fee.Deposit(2021, 4, 100, "test"); err != nil {
		t.Fatal(err)
	}

	if err = client.Disconnect(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestExempt(t *testing.T) {
	if err := fee.New(2021, 1, 100000, 15000).Exempt("20210001"); err != nil {
		t.Fatal(err)
	}
}
