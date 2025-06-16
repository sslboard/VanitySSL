package store

import (
	"context"
	"github.com/dgraph-io/badger/v3"
)

// BadgerStore implements Store using BadgerDB.
type BadgerStore struct {
	db *badger.DB
}

// NewBadgerStore opens/creates a Badger database at the given path.
func NewBadgerStore(path string) (*BadgerStore, error) {
	opts := badger.DefaultOptions(path)
	opts = opts.WithLogger(nil)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &BadgerStore{db: db}, nil
}

func (b *BadgerStore) GetDomain(ctx context.Context, domain string) (string, error) {
	var customerID string
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("domain:" + domain))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			customerID = string(val)
			return nil
		})
	})
	if err == badger.ErrKeyNotFound {
		return "", nil
	}
	return customerID, err
}

func (b *BadgerStore) SetDomain(ctx context.Context, domain, customerID string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("domain:"+domain), []byte(customerID))
	})
}

func (b *BadgerStore) DeleteDomain(ctx context.Context, domain string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte("domain:" + domain))
	})
}

func (b *BadgerStore) ListDomains(ctx context.Context) (map[string]string, error) {
	results := make(map[string]string)
	err := b.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("domain:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			domain := string(item.Key())[len("domain:"):]
			err := item.Value(func(val []byte) error {
				results[domain] = string(val)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return results, err
}

func (b *BadgerStore) GetCert(ctx context.Context, key string) ([]byte, error) {
	var data []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("cert:" + key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			data = append([]byte(nil), val...)
			return nil
		})
	})
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	return data, err
}

func (b *BadgerStore) SetCert(ctx context.Context, key string, data []byte) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("cert:"+key), data)
	})
}

func (b *BadgerStore) DeleteCert(ctx context.Context, key string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte("cert:" + key))
	})
}

// Close closes the underlying Badger database.
func (b *BadgerStore) Close() error {
	return b.db.Close()
}
