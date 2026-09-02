package stable

import "testing"

func TestQuoteDemoRate(t *testing.T) {
	b := NewBoard(10_000_000) // 0.1 KAS per 1.00
	q, err := b.Quote(3_500_000)
	if err != nil {
		t.Fatal(err)
	}
	if q.Display != "kUSD 3.500000" {
		t.Fatal(q.Display)
	}
	if q.KASDueNow != 0.35 {
		t.Fatalf("kas=%v sompi=%d", q.KASDueNow, q.SompiDueNow)
	}
	if q.L2 || q.WorkCredits || q.AssetStatus != Status {
		t.Fatalf("%+v", q)
	}
}

func TestSetRate(t *testing.T) {
	b := NewBoard(1)
	if err := b.Set(0); err == nil {
		t.Fatal("expected error")
	}
	if err := b.Set(20_000_000); err != nil {
		t.Fatal(err)
	}
	if b.Rate() != 20_000_000 {
		t.Fatal(b.Rate())
	}
}
