package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client = nil
var ic *IntouchClient = nil

var providerA = Provider{primitive.NilObjectID, "Hello World", "HW", "123456", "salt", []Patient{}}

func TestInsertProvider(t *testing.T) {
	id, err := ic.InsertProvider(providerA.Name, providerA.Username, providerA.Password)

	if err != nil {
		t.Errorf("Insert failed with %s", err.Error())
	}
	if !isExistProvider(bson.D{{"_id", *id}}) {
		t.Errorf("No provider with id %s", *id)
	}

	fmt.Println("Inserted a single document: ", id)
	ic.DeleteProvider(*id)
}

func TestAuthenticateProvider(t *testing.T) {
	providerB := Provider{primitive.NilObjectID, "Brian Lim", "blim", "catsarecute", "salt", []Patient{}}
	provID, _ := ic.InsertProvider(providerB.Name, providerB.Username, providerB.Password)
	id, err := ic.AuthenticateProvider(providerB.Username, providerB.Password)

	if err != nil {
		t.Errorf("Authentication failed with %s", err.Error())
	} else if provID != nil && *provID != *id {
		t.Errorf("Authentication does not match %s, got %s", *provID, *id)
	}

	ic.DeleteProvider(*provID)
}

func TestDeleteProvider(t *testing.T) {
	id, _ := ic.InsertProvider(providerA.Name, providerA.Username, providerA.Password)
	err := ic.DeleteProvider(*id)
	if err != nil {
		t.Errorf("Delete failed with %s", err.Error())
	}

	if isExistProvider(bson.D{{"_id", *id}}) {
		t.Errorf("Delete failed to delete %s", *id)
	}

}

func setup() {
	client = OpenConnection()
	ic = CreateIntouchClient("test", client)
	providerA = Provider{primitive.NilObjectID, "Hello World", "HW", "123456", "salt", []Patient{}}
}

func teardown() {
	client.Disconnect(context.TODO())
}

func isExistProvider(filter bson.D) bool {
	prov, err := ic.FindProvider(filter)
	if err != nil {
		fmt.Printf(err.Error())
	}
	if prov != nil {
		return true
	}
	return false
}

func TestMain(m *testing.M) {
	setup()
	run := m.Run()
	teardown()
	os.Exit(run)
}
