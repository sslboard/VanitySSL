package store

import (
	"testing"
)

func TestBadgerStore(t *testing.T) {
	path := t.TempDir()
	st, err := NewBadgerStore(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer st.Close()

	if err := st.SetDomain(nil, "example.com", "cust1"); err != nil {
		t.Fatalf("set: %v", err)
	}
	id, err := st.GetDomain(nil, "example.com")
	if err != nil || id != "cust1" {
		t.Fatalf("get: %v %s", err, id)
	}
	list, err := st.ListDomains(nil)
	if err != nil || len(list) != 1 {
		t.Fatalf("list: %v", err)
	}
	if err := st.DeleteDomain(nil, "example.com"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	id, err = st.GetDomain(nil, "example.com")
	if err != nil || id != "" {
		t.Fatalf("post delete get: %v %s", err, id)
	}

	// certificate operations
	if err := st.SetCert(nil, "cert1", []byte("data")); err != nil {
		t.Fatalf("set cert: %v", err)
	}
	certData, err := st.GetCert(nil, "cert1")
	if err != nil || string(certData) != "data" {
		t.Fatalf("get cert: %v %s", err, string(certData))
	}
	if err := st.DeleteCert(nil, "cert1"); err != nil {
		t.Fatalf("delete cert: %v", err)
	}
	certData, err = st.GetCert(nil, "cert1")
	if err != nil || certData != nil {
		t.Fatalf("post delete cert: %v %v", err, certData)
	}
}
