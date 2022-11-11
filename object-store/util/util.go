package util

import (
	"fmt"
	"github.com/upadhyayap/object-store/dbstore"
	"github.com/upadhyayap/object-store/objects"
	"strconv"
)

const NamePrefix = "anand"

func SaveTestPerson(store dbstore.KvStore, persons []*objects.Person) {
	for _, person := range persons {
		err := store.Save(person)
		fmt.Printf("Storing person with id %s and name %s\n", person.GetID(), person.GetName())
		if err != nil {
			fmt.Println("Error in saving object")
			return
		}
	}
}

func GetTestPersons()[]*objects.Person {
	persons := make([]*objects.Person,0)
	for i := 0; i <= 10 ; i++ {
		testPerson := &objects.Person{Name: NamePrefix + strconv.Itoa(i), ID: strconv.Itoa(i)}
		persons = append(persons, testPerson)
	}

	return persons
}