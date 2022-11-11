package main

import (
	"context"
	"fmt"
	"github.com/upadhyayap/object-store/config"
	"github.com/upadhyayap/object-store/dbstore"
	"github.com/upadhyayap/object-store/util"
	"os"
)

func main() {
	ctx := context.TODO()
	defer ctx.Done()
	if err := execute(ctx); err != nil {
		fmt.Printf("Failed to run the object store due to an error, %s\n", err)
		os.Exit(1)
	}
}

func execute(ctx context.Context) error {
	conf, err := config.Load()
	if err != nil {
		return err
	}

	store, err := dbstore.NewObjectStore(ctx,conf)
	if err != nil {
		return err
	}

	defer func() {
		if err := store.Close(); err != nil {
			fmt.Println("Error in closing the connection")
		}
	}()

	persons := util.GetTestPersons()
	util.SaveTestPerson(store.KvStore, persons)

	// Get person By Id
	fmt.Println("*************************************************************************")
	fmt.Println("Trying to get the person by Id: 1")
	personById, err := store.GetObjectByID("1")
	if err != nil {
		fmt.Println("Person with id 1 is not found")
		return err
	}

	if personById == nil {
		fmt.Printf("Person with id %s is not found\n", "1")
	} else {
		fmt.Printf("Person with Id %s and name %s found\n", personById.GetID(), personById.GetName())
	}


	// Get person By name
	fmt.Println("*************************************************************************")
	fmt.Printf("Trying to get the person by name: %s\n", util.NamePrefix + "2")
	personByName, err := store.GetObjectByName(util.NamePrefix + "2")
	if err != nil {
		fmt.Printf("Person with name %s is not found\n", util.NamePrefix + "2")
		return err
	}

	if personByName == nil {
		fmt.Printf("Person with name %s is not found\n", util.NamePrefix + "2")
	} else {
		fmt.Printf("Person with Id %s and name %s found\n", personByName.GetID(), personByName.GetName())
	}

	// List Person by Kind
	fmt.Println("*************************************************************************")
	fmt.Println("Listing all the persons by kind")
	personsByKind, err := store.ListObjects(persons[0].GetKind())

	for _, person := range personsByKind {
		fmt.Printf("Person with Id %s and name %s found\n", person.GetID(), person.GetName())
	}

	fmt.Println("*************************************************************************")
	fmt.Println("Deleting person with id:5")

	// Delete object with ID:5
	err = store.DeleteObject("5")
	if err != nil {
		fmt.Println("Error in deleting")
	}

	fmt.Println("Successfully deleted the person with id:5")
	fmt.Println("Trying to get the person with id:5")

	// Try to get the deleted object
	personById, err = store.GetObjectByID("5")
	if err != nil {
		fmt.Println("Error in retrieving the person")
		return err
	}

	if personById == nil {
		fmt.Println("person with id:5 was deleted not able to find it")
	}

	fmt.Println("Done: Thanks for checking the exercise!!!")

	return nil
}
