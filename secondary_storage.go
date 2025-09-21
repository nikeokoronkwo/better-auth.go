package main

import (
	"fmt"
	"sync"
)

// Secondary Storage is used as a faster storage source for
// augmenting and syncing with the DB
//
// Uses include, but are not limited to:
// - (possible) list of usernames
// - rate limiting counters
// - session data
//
// some options use the default secondary storage of memory
type SecondaryStorage interface {
	Get(key string) (string, error)
	Set(key, value string) error
	SetTTL(key, value string, ttl int) error
	Delete(key string) error
}

// An in memory, basic implementation of {SecondaryStorage}
//
// The following does not support TTL. For that, use (TODO: Implement TTL storage)
type InMemoryStorage struct {
	data sync.Map
}

// Delete implements SecondaryStorage.
func (i* InMemoryStorage) Delete(key string) error {
	i.data.Delete(key)
	return nil
}

// Get implements SecondaryStorage.
func (i* InMemoryStorage) Get(key string) (string, error) {
	v, ok := i.data.Load(key)
	if !ok {
		return "", fmt.Errorf("Could not find the key %s in the storage", key)
	} else {
		return v.(string), nil
	}
}

// Set implements SecondaryStorage.
func (i* InMemoryStorage) Set(key string, value string) error {
	i.data.Store(key, value)
	return nil
}

func (i* InMemoryStorage) SetTTL(key string, value string, ttl int) error {
	panic("unsupported")
}

var _ SecondaryStorage = &InMemoryStorage{}
