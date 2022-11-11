package dbstore

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/upadhyayap/object-store/config"
	"github.com/upadhyayap/object-store/objects"
)

type ObjectDB interface {
	// Store will store the objects in the data store. The objects will have a
	// name and kind, and the Store method should create a unique ID.
	Store(object objects.Object) error

	// GetObjectByID will retrieve the objects with the provided ID.
	GetObjectByID(id string) (objects.Object, error)

	// GetObjectByName will retrieve the objects with the given name.
	GetObjectByName(id string) (objects.Object, error)

	// ListObjects will return a list of all objects of the given kind.
	ListObjects(kind string) ([]objects.Object, error)

	// DeleteObject will delete the objects.
	DeleteObject(id string) error
}

type ObjectStore struct {
	KvStore KvStore
}

func (obs *ObjectStore) Store(object objects.Object) error {
	if err := obs.KvStore.Save(object); err != nil {
		return err
	}
	return nil
}

func (obs *ObjectStore) GetObjectByID(id string) (objects.Object, error) {
	object, err := obs.KvStore.GetByID(id)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (obs *ObjectStore) GetObjectByName(id string) (objects.Object, error) {
	object, err := obs.KvStore.GetByName(id)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (obs *ObjectStore) ListObjects(kind string) ([]objects.Object, error) {
	allObjects, err := obs.KvStore.List(kind)
	if err != nil {
		return []objects.Object{}, err
	}
	return allObjects, nil
}

func (obs *ObjectStore) DeleteObject(id string) error {
	if err := obs.KvStore.Delete(id); err != nil {
		return err
	}
	return nil
}

func (obs ObjectStore) Close() error {
	if err := obs.KvStore.Close(); err != nil {
		return err
	}

	return nil
}

type KvStore interface {
	Save(object objects.Object) error
	GetByID(id string) (objects.Object, error)
	GetByName(name string) (objects.Object, error)
	List(kind string) ([]objects.Object, error)
	Delete(id string) error
	Close() error
}

type RedisDB struct {
	Client     *redis.Client
	Ctx        context.Context
}

// If an object already exists then it will not be saved again
func (db *RedisDB) Save(object objects.Object) error {
	isKeyExists := db.Client.Exists(db.Ctx, object.GetID())
	if err := isKeyExists.Err(); err != nil {
		return err
	}

	if isKeyExists.Val() == 1 {
		return nil
	}

	values := map[string]string{"name":object.GetName(), "kind":object.GetKind()}
	if err := db.Client.HSet(db.Ctx, object.GetID(), values).Err(); err != nil {
	 	return err
	}

	if err := db.Client.LPush(db.Ctx, object.GetName(), object.GetID()).Err(); err != nil {
		return err
	}

	if err := db.Client.LPush(db.Ctx, object.GetKind(), object.GetID()).Err(); err != nil {
	 	return err
	}

	return nil
}

func (db *RedisDB) GetByID(id string) (objects.Object, error) {
	byId := db.Client.HGetAll(db.Ctx, id)
	if err := byId.Err(); err != nil {
		return nil, err
	}

	// Return nil if object is not found in the cache
	if len(byId.Val()) == 0 {
		return nil, nil
	}

	obs := &objects.Person{}
	obs.SetID(id)
	obs.SetName(byId.Val()["name"])

	return obs, nil
}

func (db *RedisDB) GetByName(name string) (objects.Object, error) {
	byName := db.Client.LRange(db.Ctx, name, 0, 0)

	if err := byName.Err(); err != nil {
		return nil, err
	}

	if len(byName.Val()) == 0 {
		return nil, nil
	}

	byId := db.Client.HGetAll(db.Ctx, byName.Val()[0])
	if err := byId.Err(); err != nil {
		return nil, err
	}

	// Return nil if object is not found in the cache
	if len(byId.Val()) == 0 {
		return nil, nil
	}

	obs := &objects.Person{}
	obs.SetID(byName.Val()[0])
	obs.SetName(byId.Val()["name"])

	return obs, nil
}

func (db *RedisDB) List(kind string) ([]objects.Object, error) {
	var allObjects = make([]objects.Object,0)
	ids := db.Client.LRange(db.Ctx, kind, 0, -1)
	if err := ids.Err(); err != nil {
		return nil, err
	}

	for _, id := range ids.Val() {
		byId := db.Client.HGetAll(db.Ctx, id)
		if err := byId.Err(); err != nil {
			return nil, err
		}
		obs := &objects.Person{}
		obs.SetID(id)
		obs.SetName(byId.Val()["name"])
		allObjects = append(allObjects, obs)
	}

	return allObjects, nil
}

func (db *RedisDB) Delete(id string) error {
	byId := db.Client.HGetAll(db.Ctx, id)
	if err := byId.Err(); err != nil {
		return err
	}

	// Return nil if object is already deleted
	if len(byId.Val()) == 0 {
		return nil
	}

	if err := db.Client.Del(db.Ctx, id).Err(); err != nil {
		return err
	}

	if err := db.Client.LRem(db.Ctx, byId.Val()["name"], 1, id).Err(); err != nil {
		return err
	}

	if err := db.Client.LRem(db.Ctx, byId.Val()["kind"], 1, id).Err(); err != nil {
		return err
	}

	return nil
}

func (db *RedisDB) Close() error {
	err := db.Client.Close()
	if err != nil {
		return err
	}

	return nil
}

func NewObjectStore(ctx context.Context, conf *config.ObjectStoreConfig) (*ObjectStore, error) {
	var obs *ObjectStore
	store, err := NewDB(ctx, conf.DBConfig)

	if err != nil {
		return &ObjectStore{}, err
	}

	obs = &ObjectStore{KvStore: store}

	return obs, nil
}

// Factory to initialize the underlying data store
func NewDB(ctx context.Context, storeConfig *config.DBStoreConfig) (KvStore, error) {
	switch storeConfig.Type {
	case "redis":
		redisDB, err := NewRedisDB(ctx, storeConfig)
		if err != nil {
			return nil, err
		}
		return redisDB, nil
	default:
		return nil, errors.New(fmt.Sprintf("%s DB is not supported", storeConfig.Type))
	}
}

func NewRedisDB(ctx context.Context, conf *config.DBStoreConfig) (*RedisDB, error) {
	rClient := redis.NewClient(&redis.Options{
		Addr:     conf.Host + ":" + conf.Port,
		Password: conf.Password,
		DB:       0, // using default db
	})

	// redis client does not gives error so I am doing an ping to check that the client is able to connect
	_, err := rClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	rdb := &RedisDB{Client: rClient, Ctx: ctx}

	return rdb, nil
}
