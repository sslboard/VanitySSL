package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/vanityssl/internal/store"
)

func TestAPI(t *testing.T) {
	st, err := store.NewBadgerStore(t.TempDir())
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	api := New(st, "")
	srv := httptest.NewServer(api.Router())
	defer srv.Close()

	c := map[string]string{"domain": "example.com", "customer_id": "c1"}
	body, _ := json.Marshal(c)
	resp, err := http.Post(srv.URL+"/customers", "application/json", bytes.NewReader(body))
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("create: %v status %d", err, resp.StatusCode)
	}

	resp, err = http.Get(srv.URL + "/customers/example.com")
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("get: %v status %d", err, resp.StatusCode)
	}
}
