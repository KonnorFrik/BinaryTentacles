package fake_db

import (
	"context"
	"maps"
	"slices"
	"sync"
)

type Db struct {
	idMut    sync.Mutex
	itemsMut sync.Mutex
	nextId   uint64
	items    map[uint64]any
}

func New() *Db {
	return &Db{
		items: make(map[uint64]any),
	}
}

// Create - create an object in Db.
// Returns id of created object.
// If object already exist - overwrite it.
func (d *Db) Create(
	ctx context.Context,
	obj any,
) uint64 {
	d.idMut.Lock()
	var id = d.nextId
	d.nextId++
	d.idMut.Unlock()

	d.itemsMut.Lock()
	d.items[id] = obj
	d.itemsMut.Unlock()
	return id
}

// Get - return stored object by id if exist.
func (d *Db) Get(
	ctx context.Context,
	id uint64,
) (
	any,
	bool,
) {
	d.itemsMut.Lock()
	defer d.itemsMut.Unlock()

	obj, ok := d.items[id]
	return obj, ok
}

func (d *Db) All(
	ctx context.Context,
) []any {
	d.itemsMut.Lock()
	defer d.itemsMut.Unlock()
	return slices.Collect(maps.Values(d.items))
}

type dbGetted interface {
	Get(context.Context, uint64) (any, bool)
}

func As[T any](
	ctx context.Context,
	db dbGetted,
	id uint64,
) (
	T,
	bool,
) {
	obj, exist := db.Get(ctx, id)

	if !exist {
		var empty T
		return empty, false
	}

	typedObj, ok := obj.(T)
	return typedObj, ok
}

type dbAll interface {
	All(context.Context) []any
}

// All - return all values of type T.
// Other types will be rejected.
func All[T any](
	ctx context.Context,
	db dbAll,
) []T {
	var result []T

	for _, val := range db.All(ctx) {
		typedVal, ok := val.(T)

		if ok {
			result = append(result, typedVal)
		}
	}

	return result
}
