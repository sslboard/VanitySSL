package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/example/vanityssl/internal/model"
	"github.com/example/vanityssl/internal/store"
	"github.com/gorilla/mux"
)

// API wraps the HTTP handlers for managing customers.
type API struct {
	Store store.Store
	Token string
}

func New(store store.Store, token string) *API {
	return &API{Store: store, Token: token}
}

func (a *API) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.Token != "" {
			if r.Header.Get("Authorization") != "Bearer "+a.Token {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// Router returns the HTTP router for the API.
func (a *API) Router() http.Handler {
	r := mux.NewRouter()
	r.Use(a.auth)
	r.HandleFunc("/customers", a.createCustomer).Methods(http.MethodPost)
	r.HandleFunc("/customers", a.listCustomers).Methods(http.MethodGet)
	r.HandleFunc("/customers/{domain}", a.getCustomer).Methods(http.MethodGet)
	r.HandleFunc("/customers/{domain}", a.deleteCustomer).Methods(http.MethodDelete)
	return r
}

func (a *API) createCustomer(w http.ResponseWriter, r *http.Request) {
	var c model.Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if c.Domain == "" || c.CustomerID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := a.Store.SetDomain(r.Context(), strings.ToLower(c.Domain), c.CustomerID); err != nil {
		log.Printf("error setting domain: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *API) listCustomers(w http.ResponseWriter, r *http.Request) {
	m, err := a.Store.ListDomains(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	customers := make([]model.Customer, 0, len(m))
	for domain, id := range m {
		customers = append(customers, model.Customer{Domain: domain, CustomerID: id})
	}
	json.NewEncoder(w).Encode(customers)
}

func (a *API) getCustomer(w http.ResponseWriter, r *http.Request) {
	domain := mux.Vars(r)["domain"]
	id, err := a.Store.GetDomain(r.Context(), strings.ToLower(domain))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(model.Customer{Domain: domain, CustomerID: id})
}

func (a *API) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	domain := mux.Vars(r)["domain"]
	if err := a.Store.DeleteDomain(r.Context(), strings.ToLower(domain)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
