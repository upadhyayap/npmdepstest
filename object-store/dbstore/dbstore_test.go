package dbstore

import (
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/upadhyayap/object-store/objects"
	"math/rand"
	"sort"
	"strconv"
	"testing"
)

var redisServer *miniredis.Miniredis
var rdb *RedisDB
var store *ObjectStore

func Test_Store(t *testing.T) {
	setup()
	defer teardown()

	testPerson := &objects.Person{Name: "Anand", ID: "1"}

	err := store.Store(testPerson)

	if err != nil {
		t.Error("Failed to save object")
	}
}

func TestRedisDB_GetByID(t *testing.T) {
	setup()
	defer teardown()

	persons := GetTestPersons()
	SaveTestPerson(store.KvStore, persons)

	id := rand.Intn(10)
	personId := strconv.Itoa(id)

	testPerson, err := rdb.GetByID(personId)

	if err != nil {
		t.Error("Failed to get the object by id")
	}

	assert.NotNil(t, testPerson, "Unable to find the object")
	assert.Equal(t, "test" + personId, testPerson.GetName(), "found person with wrong name")
	assert.Equal(t, personId, testPerson.GetID(), "found person with wrong id")
}

func TestRedisDB_GetByName(t *testing.T) {
	setup()
	defer teardown()

	persons := GetTestPersons()
	SaveTestPerson(store.KvStore, persons)

	id := rand.Intn(10)
	personId := strconv.Itoa(id)

	testPerson, err := rdb.GetByName("test" + personId)

	if err != nil {
		t.Error("Failed to get the object by name")
	}

	assert.NotNil(t, testPerson, "Unable to find the object")
	assert.Equal(t, "test" + personId, testPerson.GetName(), "found person with wrong name")
	assert.Equal(t, personId, testPerson.GetID(), "found person with wrong id")
}

func Test_ListByKind(t *testing.T) {
	setup()
	defer teardown()

	persons := GetTestPersons()
	SaveTestPerson(store.KvStore, persons)

	personsFromCache, err := rdb.List("*objects.Person")
	if err != nil {
		t.Error("Failed to list objects by kind")
	}

	sort.Slice(personsFromCache[:], func(i, j int) bool {
		firstId, _ := strconv.Atoi(personsFromCache[i].GetID())
		secondId, _ := strconv.Atoi(personsFromCache[j].GetID())
		return  firstId < secondId
	})

	for i := 0; i < len(persons); i++ {
		testPerson, err := rdb.GetByID(personsFromCache[i].GetID())

		if err != nil {
			t.Error("Failed to get the object by id")
		}

		assert.NotNil(t, testPerson, "Unable to find the object")
		assert.Equal(t, "test" + strconv.Itoa(i), testPerson.GetName(), "found person with wrong name")
		assert.Equal(t, strconv.Itoa(i), testPerson.GetID(), "found person with wrong id")
	}

}

func TestRedisDB_Delete(t *testing.T) {
	setup()
	defer teardown()

	persons := GetTestPersons()
	SaveTestPerson(store.KvStore, persons)

	err := rdb.Delete("1")

	if err != nil {
		t.Error("Failed to delete object")
	}

	testPerson, err := rdb.GetByID("1")
	if err != nil {
		t.Error("Failed to get the object by id")
	}

	assert.Nil(t, testPerson, "Found the deleted object")
}

func mockRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()

	if err != nil {
		return nil
	}

	return s
}

func setup() {
	redisServer = mockRedis()
	if redisServer == nil {
		return
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
	rdb = &RedisDB{Client: redisClient, Ctx: context.Background()}
	store = &ObjectStore{KvStore: rdb}
}

func teardown() {
	redisServer.Close()
}

func SaveTestPerson(store KvStore, persons []*objects.Person) {
	for _, person := range persons {
		err := store.Save(person)
		if err != nil {
			fmt.Println("Error in saving object")
			return
		}
	}
}

func GetTestPersons()[]*objects.Person {
	persons := make([]*objects.Person,0)
	for i := 0; i <= 10 ; i++ {
		testPerson := &objects.Person{Name: "test" + strconv.Itoa(i), ID: strconv.Itoa(i)}
		persons = append(persons, testPerson)
	}

	return persons
}