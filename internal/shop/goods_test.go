package shop

import "testing"

func TestCatalog(t *testing.T) {
	if len(Catalog) < 3 {
		t.Fatal("catalog")
	}
	g, ok := Get("cup")
	if !ok || g.Micro != 3_500_000 {
		t.Fatalf("%+v", g)
	}
}
