// Package stable is the reserved Kaspa L1 unit of account.
// It is not live. It is not an L2 bridged dollar. It is not Work Credits.
// When a capitalized issuer (or an overcollateral vault) ships a KCC-20 on L1,
// this dApp already invoices in that unit.
package stable

import (
	"fmt"
	"math"
	"sync"
)

const (
	Code       = "kUSD"
	Decimals   = 6
	MicroUnit  = 1_000_000
	SompiPerKAS = 100_000_000
	Status     = "reserved-l1-not-live"
)

type Quote struct {
	Micro          uint64  `json:"micro"`
	Display        string  `json:"display"`
	Asset          string  `json:"asset"`
	AssetStatus    string  `json:"assetStatus"`
	SompiPerUnit   uint64  `json:"sompiPerUnit"`
	SompiDueNow    uint64  `json:"sompiDueNow"`
	KASDueNow      float64 `json:"kasDueNow"`
	SchemeLive     string  `json:"schemeLive"`
	SchemeReserved string  `json:"schemeReserved"`
	Note           string  `json:"note"`
	L2             bool    `json:"l2"`
	WorkCredits    bool    `json:"workCredits"`
}

type Board struct {
	mu           sync.Mutex
	sompiPerUnit uint64 // sompi per 1.00 kUSD. Merchant-posted. Not an oracle.
}

func NewBoard(sompiPerUnit uint64) *Board {
	if sompiPerUnit == 0 {
		sompiPerUnit = 10_000_000 // demo: 0.10 KAS per reserved 1.00 — not a market
	}
	return &Board{sompiPerUnit: sompiPerUnit}
}

func (b *Board) Rate() uint64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.sompiPerUnit
}

func (b *Board) Set(sompiPerUnit uint64) error {
	if sompiPerUnit == 0 {
		return fmt.Errorf("rate must be > 0")
	}
	b.mu.Lock()
	b.sompiPerUnit = sompiPerUnit
	b.mu.Unlock()
	return nil
}

func Format(micro uint64) string {
	whole := micro / MicroUnit
	frac := micro % MicroUnit
	return fmt.Sprintf("%s %d.%06d", Code, whole, frac)
}

func (b *Board) Quote(micro uint64) (Quote, error) {
	if micro == 0 {
		return Quote{}, fmt.Errorf("amount must be > 0")
	}
	rate := b.Rate()
	if micro > math.MaxUint64/rate {
		return Quote{}, fmt.Errorf("overflow")
	}
	// sompi = micro * rate / 1e6
	sompi := (micro * rate) / MicroUnit
	return Quote{
		Micro:          micro,
		Display:        Format(micro),
		Asset:          Code,
		AssetStatus:    Status,
		SompiPerUnit:   rate,
		SompiDueNow:    sompi,
		KASDueNow:      float64(sompi) / float64(SompiPerKAS),
		SchemeLive:     "kaspa",
		SchemeReserved: "kaspa-l1-stable",
		L2:             false,
		WorkCredits:    false,
		Note:           "Unit of account is a reserved Kaspa L1 stable (kUSD). Not live. Settle today in KAS at the merchant-posted rate (not an oracle, not Circle). No L2. No gram vouchers.",
	}, nil
}
