package server

import (
	"encoding/json"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"kastill/internal/feedback"
	"kastill/internal/shop"
	"kastill/internal/stable"
	"kastill/internal/wallets"
	"kastill/web"
)

type Server struct {
	Addr  string
	T     *template.Template
	Board *stable.Board
	mu    sync.Mutex
	Orders []Order
}

type Order struct {
	ID      int           `json:"id"`
	Good    shop.Good     `json:"good"`
	Quote   stable.Quote  `json:"quote"`
	Paid    string        `json:"paid"`
	Note    string        `json:"note"`
}

type page struct {
	Title  string
	Active string
	Error  string
	Query  string
	Goods  []shop.Good
	Good   *shop.Good
	Quote  *stable.Quote
	Rate    uint64
	Orders  []Order
	Wallets []wallets.Wallet
}

func New(addr string) (*Server, error) {
	t, err := template.ParseFS(web.Templates, "templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Server{Addr: addr, T: t, Board: stable.NewBoard(10_000_000)}, nil
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	static, err := fs.Sub(web.Static, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static))))
	mux.HandleFunc("/", s.home)
	mux.HandleFunc("/shop", s.shopPage)
	mux.HandleFunc("/item/", s.item)
	mux.HandleFunc("/vision", s.vision)
	mux.HandleFunc("/honest", s.honest)
	mux.HandleFunc("/rate", s.rate)
	mux.HandleFunc("/orders", s.orders)
	mux.HandleFunc("/buy", s.buy)
	mux.HandleFunc("/wallets", func(w http.ResponseWriter, r *http.Request) {
		s.render(w, "wallets.html", page{Title: "Wallets · Kaspa Till", Active: "wallets", Wallets: wallets.All()})
	})
	mux.HandleFunc("/api/wallets", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"ok": true, "inject": []string{"kasware", "kastle"}, "ledger": "https://kasvault.io", "data": wallets.All()})
	})
	mux.HandleFunc("/safety", func(w http.ResponseWriter, r *http.Request) {
		s.render(w, "safety.html", page{Title: "Safety · Kaspa Till", Active: "safety"})
	})
	mux.HandleFunc("/idea", func(w http.ResponseWriter, r *http.Request) {
		s.render(w, "idea.html", page{Title: "Idea · Kaspa Till", Active: "idea"})
	})
	mux.HandleFunc("/explain", func(w http.ResponseWriter, r *http.Request) {
		s.render(w, "idea.html", page{Title: "Idea · Kaspa Till", Active: "idea"})
	})
	mux.HandleFunc("/why", func(w http.ResponseWriter, r *http.Request) {
		s.render(w, "why.html", page{Title: "Why · Kaspa Till", Active: "why"})
	})
	mux.HandleFunc("/234", func(w http.ResponseWriter, r *http.Request) {
		s.render(w, "framing.html", page{Title: "#234 · Kaspa Till", Active: "why"})
	})
	mux.HandleFunc("/feedback", s.feedbackPage)
	mux.HandleFunc("/api/feedback", s.apiFeedback)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{
			"ok": true, "dapp": "kastill", "layer": "kaspa-l1",
			"unit": stable.Code, "status": stable.Status,
			"l2": false, "workCredits": false, "stableLive": false,
			"sompiPerUnit": s.Board.Rate(),
		})
	})
	mux.HandleFunc("/api/goods", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"ok": true, "data": shop.Catalog, "asset": stable.Code, "status": stable.Status})
	})
	mux.HandleFunc("/api/quote", s.apiQuote)
	mux.HandleFunc("/api/order", s.apiOrder)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/static/") {
			log.Printf("%s %s", r.Method, r.URL.Path)
		}
		mux.ServeHTTP(w, r)
	})
}

func (s *Server) render(w http.ResponseWriter, name string, p page) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	p.Rate = s.Board.Rate()
	if err := s.T.ExecuteTemplate(w, name, p); err != nil {
		log.Println("template", name, err)
		http.Error(w, err.Error(), 500)
	}
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	s.render(w, "home.html", page{Title: "Kaspa Till — reserved L1 stable", Active: "home", Goods: shop.Catalog})
}

func (s *Server) shopPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "shop.html", page{Title: "Shop · Kaspa Till", Active: "shop", Goods: shop.Catalog})
}

func (s *Server) item(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/item/")
	g, ok := shop.Get(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	q, err := s.Board.Quote(g.Micro)
	p := page{Title: g.Title + " · Kaspa Till", Active: "shop", Good: &g}
	if err != nil {
		p.Error = err.Error()
	} else {
		p.Quote = &q
	}
	s.render(w, "item.html", p)
}

func (s *Server) vision(w http.ResponseWriter, r *http.Request) {
	s.render(w, "vision.html", page{Title: "Vision · Kaspa Till", Active: "vision"})
}

func (s *Server) honest(w http.ResponseWriter, r *http.Request) {
	s.render(w, "honest.html", page{Title: "Claims · Kaspa Till", Active: "honest"})
}

func (s *Server) orders(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	list := append([]Order(nil), s.Orders...)
	s.mu.Unlock()
	s.render(w, "orders.html", page{Title: "Orders · Kaspa Till", Active: "orders", Orders: list})
}

func (s *Server) rate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		n, err := strconv.ParseUint(strings.TrimSpace(r.FormValue("sompiPerUnit")), 10, 64)
		if err != nil {
			s.render(w, "rate.html", page{Title: "Rate", Active: "rate", Error: "need sompi per 1.00 kUSD"})
			return
		}
		if err := s.Board.Set(n); err != nil {
			s.render(w, "rate.html", page{Title: "Rate", Active: "rate", Error: err.Error()})
			return
		}
		http.Redirect(w, r, "/shop", http.StatusSeeOther)
		return
	}
	s.render(w, "rate.html", page{Title: "Merchant rate · Kaspa Till", Active: "rate"})
}

func (s *Server) buy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/shop", http.StatusSeeOther)
		return
	}
	g, ok := shop.Get(r.FormValue("item"))
	if !ok {
		s.render(w, "shop.html", page{Title: "Shop", Active: "shop", Goods: shop.Catalog, Error: "unknown item"})
		return
	}
	paid := strings.TrimSpace(r.FormValue("payment"))
	q, err := s.Board.Quote(g.Micro)
	if err != nil {
		s.render(w, "item.html", page{Title: g.Title, Active: "shop", Good: &g, Error: err.Error()})
		return
	}
	if paid == "" {
		s.render(w, "item.html", page{
			Title: g.Title, Active: "shop", Good: &g, Quote: &q,
			Error: "Payment required in KAS at the merchant rate. kUSD is reserved, not spendable.",
		})
		return
	}
	s.mu.Lock()
	ord := Order{ID: len(s.Orders) + 1, Good: g, Quote: q, Paid: paid, Note: "KAS settlement recorded at HTTP layer. Reserved kUSD amount is the invoice, not a transfer."}
	s.Orders = append(s.Orders, ord)
	s.mu.Unlock()
	s.render(w, "receipt.html", page{Title: "Receipt · Kaspa Till", Active: "orders", Good: &g, Quote: &q, Orders: []Order{ord}})
}

func (s *Server) apiQuote(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("item")
	g, ok := shop.Get(id)
	if !ok {
		writeJSON(w, 404, map[string]any{"ok": false, "error": "unknown item"})
		return
	}
	q, err := s.Board.Quote(g.Micro)
	if err != nil {
		writeJSON(w, 400, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true, "good": g, "quote": q})
}

func (s *Server) apiOrder(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("item")
	if id == "" && r.Method == http.MethodPost {
		body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<16))
		var req struct{ Item, Payment string }
		_ = json.Unmarshal(body, &req)
		id = req.Item
		if req.Payment != "" {
			r.Header.Set("X-Kaspa-Payment", req.Payment)
		}
	}
	g, ok := shop.Get(id)
	if !ok {
		writeJSON(w, 404, map[string]any{"ok": false, "error": "unknown item"})
		return
	}
	q, err := s.Board.Quote(g.Micro)
	if err != nil {
		writeJSON(w, 400, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	paid := strings.TrimSpace(r.Header.Get("X-Kaspa-Payment"))
	if paid == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusPaymentRequired)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "Payment Required",
			"note":  "kUSD is reserved and not live. Pay KAS now at the merchant-posted rate. Accepts kaspa-l1-stable only when that KCC-20 exists. No L2. No work credits.",
			"good":  g,
			"quote": q,
			"accepts": []map[string]any{
				{"scheme": "kaspa", "asset": "KAS", "sompi": q.SompiDueNow, "live": true},
				{"scheme": "kaspa-l1-stable", "asset": stable.Code, "micro": g.Micro, "live": false},
			},
		})
		return
	}
	s.mu.Lock()
	ord := Order{ID: len(s.Orders) + 1, Good: g, Quote: q, Paid: paid, Note: "HTTP receipt. Not an L1 kUSD transfer."}
	s.Orders = append(s.Orders, ord)
	s.mu.Unlock()
	writeJSON(w, 200, map[string]any{"ok": true, "order": ord})
}

func (s *Server) feedbackPage(w http.ResponseWriter, r *http.Request) {
	p := page{Title: "Feedback · Kaspa Till", Active: "feedback"}
	if r.Method == http.MethodPost {
		n, err := feedback.Save("kastill", r.FormValue("text"), r.FormValue("contact"))
		if err != nil {
			p.Error = err.Error()
		} else {
			p.Query = n.ID
		}
	}
	s.render(w, "feedback.html", p)
}

func (s *Server) apiFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, 405, map[string]any{"ok": false, "error": "POST"})
		return
	}
	n, err := feedback.Save("kastill", r.FormValue("text"), r.FormValue("contact"))
	if err != nil {
		body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<16))
		var req struct{ Text, Contact string }
		_ = json.Unmarshal(body, &req)
		n, err = feedback.Save("kastill", req.Text, req.Contact)
	}
	if err != nil {
		writeJSON(w, 400, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true, "id": n.ID, "dir": feedback.Dir(), "note": "This site never DMs you."})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
