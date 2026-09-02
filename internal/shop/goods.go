package shop

import "fmt"

type Good struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Blurb string `json:"blurb"`
	Micro uint64 `json:"micro"` // reserved L1 stable, 6 decimals
}

var Catalog = []Good{
	{ID: "cup", Title: "BlockDAG cup", Blurb: "Physical merch. Priced in reserved kUSD so the shelf tag does not move when KAS does.", Micro: 3_500_000},
	{ID: "pass", Title: "Month of name-desk access", Blurb: "Human subscription. Same number next month, in the reserved unit. KAS due today floats with the merchant rate.", Micro: 12_000_000},
	{ID: "lease", Title: "Subname lease (convention)", Blurb: "pay.shop.kas style lease as a commercial object. Settlement is L1 KAS until kUSD exists.", Micro: 25_000_000},
	{ID: "agent", Title: "Agent retainer (7 days)", Blurb: "An agent that may call this till. Invoice is already in the unit a future L1 stable would use.", Micro: 7_000_000},
}

func (g Good) Tag() string {
	return fmt.Sprintf("kUSD %d.%06d", g.Micro/1_000_000, g.Micro%1_000_000)
}

func Get(id string) (Good, bool) {
	for _, g := range Catalog {
		if g.ID == id {
			return g, true
		}
	}
	return Good{}, false
}
