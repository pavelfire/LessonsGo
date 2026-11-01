package domain
/*
import lru "github.com/hashicorp/golang-lru"

type Storage struct {
	cache *lru.Cache
}

func (s *Storage) Set(wh *warehouse.Warehouse) {
	s.cache.Put(wh.Id, *wh)
}

func (s *Storage) Get(id types.WarehouseId) *warehouse.Warehouse {
	item, ok := s.cache.Get(id)
	if ok {
		if wh, ok := item.(*warehouse.Warehouse); ok {
			return wh
		}
	}
	return nil
}
*/


import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
)
type Storage struct {
	cache *lru.Cache
}
func NewStorage(size int) (*Storage, error) {
	cache, err := lru.New(size)
	if err != nil {
		return nil, fmt.Errorf("create warehouse cache: %w", err)
	}
	return &Storage{cache: cache}, nil
}
func (s *Storage) Set(wh *warehouse.Warehouse) {
	if wh == nil {
		return
	}
	s.cache.Add(wh.Id, *wh)
}
func (s *Storage) Get(id types.WarehouseId) *warehouse.Warehouse {
	item, ok := s.cache.Get(id)
	if !ok {
		return nil
	}
	wh, ok := item.(warehouse.Warehouse)
	if !ok {
		return nil
	}
	return &wh
}


type WarehouseId int64


type Warehouse struct {
	Id   WarehouseId
	Name string
}