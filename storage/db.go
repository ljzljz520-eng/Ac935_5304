package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"go.etcd.io/bbolt"
	"hospitaldesk/model"
)

var bucketNames = map[string][]byte{
	"employees": []byte("employees"),
	"policies":  []byte("policies"),
	"training":  []byte("training"),
	"schedules": []byte("schedules"),
	"downloads": []byte("downloads"),
	"shares":    []byte("shares"),
	"reviews":   []byte("reviews"),
}

type Store struct {
	db   *bbolt.DB
	path string
	mu   sync.RWMutex
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	db, err := bbolt.Open(filepath.Clean(path), 0600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	s := &Store{db: db, path: path}
	if err := s.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) Path() string { return s.path }

func encode(value any) ([]byte, error) { return json.Marshal(value) }

func decode(data []byte, target any) error {
	if len(data) == 0 {
		return errors.New("empty record")
	}
	return json.Unmarshal(data, target)
}

func put[T any](s *Store, bucket []byte, key string, value T) error {
	if key == "" {
		return errors.New("record key is required")
	}
	b, err := encode(value)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("database is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucket).Put([]byte(key), b) })
}

func get[T any](s *Store, bucket []byte, key string) (T, error) {
	var value T
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return value, errors.New("database is closed")
	}
	err := s.db.View(func(tx *bbolt.Tx) error {
		data := tx.Bucket(bucket).Get([]byte(key))
		if data == nil {
			return fmt.Errorf("record %s not found", key)
		}
		return decode(data, &value)
	})
	return value, err
}

func list[T any](s *Store, bucket []byte) ([]T, error) {
	items := make([]T, 0)
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, errors.New("database is closed")
	}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).ForEach(func(_, data []byte) error {
			var item T
			if err := decode(data, &item); err != nil {
				return err
			}
			items = append(items, item)
			return nil
		})
	})
	return items, err
}

func remove(s *Store, bucket []byte, key string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("database is closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucket).Delete([]byte(key)) })
}

func (s *Store) Count(bucketName string) (int, error) {
	bucket, ok := bucketNames[bucketName]
	if !ok {
		return 0, errors.New("unknown bucket")
	}
	count := 0
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return 0, errors.New("database is closed")
	}
	err := s.db.View(func(tx *bbolt.Tx) error { count = tx.Bucket(bucket).Stats().KeyN; return nil })
	return count, err
}

func (s *Store) Health() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("database is closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketNames["employees"]).ForEach(func(_, _ []byte) error { return nil })
	})
}

var _ = model.PolicyDraft
