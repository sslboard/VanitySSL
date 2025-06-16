package store

import "testing"

func TestCachedStore(t *testing.T) {
	backend, err := NewBadgerStore(t.TempDir())
	if err != nil {
		t.Fatalf("backend: %v", err)
	}
	defer backend.Close()

	cache, err := NewCachedStore(backend, 10)
	if err != nil {
		t.Fatalf("cache: %v", err)
	}

	if err := cache.SetDomain(nil, "a.com", "id1"); err != nil {
		t.Fatalf("set: %v", err)
	}
	if v, _ := cache.GetDomain(nil, "a.com"); v != "id1" {
		t.Fatalf("get: %s", v)
	}
	if err := cache.DeleteDomain(nil, "a.com"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if v, _ := cache.GetDomain(nil, "a.com"); v != "" {
		t.Fatalf("deleted get: %s", v)
	}
}
