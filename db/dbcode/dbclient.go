package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db = "intouch"
var err = map[string]error{"ErrUsernameInvalid": errors.New("Username is already in use"),
	"ErrPatientNotFound":  errors.New("Patient was not found"),
	"ErrProviderNotFound": errors.New("Provider was not found")}

// IntouchClient will add methods for easy transactions with db
type IntouchClient struct {
	ProCol *mongo.Collection
	PatCol *mongo.Collection
}

// updateField takes a collection, filter, and update and updates one doc in collection
// matching the filter with the update
// returns error, nil if ok
func (ic *IntouchClient) updateField(collection *mongo.Collection, filter bson.D, update bson.D) error {
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Print(err)
		return err
	}
	return err
}

// updateFieldAll takes a collection, filter, and update and updates all docs in collection
// matching the filter with the update
// returns error, nil if ok
func (ic *IntouchClient) updateFieldAll(collection *mongo.Collection, filter bson.D, update bson.D) error {
	updateResult, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Print(err)
		return err
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n",
		updateResult.MatchedCount, updateResult.ModifiedCount)

	return err
}

// deleteEntity deletes docs matching filter from given collection
// return error, nil if ok
func (ic *IntouchClient) deleteEntity(collection *mongo.Collection, filter bson.D) error {
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Print(err)
		return err
	}

	fmt.Printf("Deleted %d documents\n", deleteResult.DeletedCount)
	return err
}

// FindProvider uses filter to find a provider matching the filter
// returns nil and err if it fails
func (ic *IntouchClient) FindProvider(filter bson.D) (*Provider, error) {
	var result Provider
	err := ic.ProCol.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &result, nil
}

// FindProviderByID helper for finding provider by id easily
func (ic *IntouchClient) FindProviderByID(id primitive.ObjectID) (*Provider, error) {
	return ic.FindProvider(bson.D{{"_id", id}})
}

// FindPatient uses filter to find a patient matching the filter
// returns nil and err if it fails
func (ic *IntouchClient) FindPatient(filter bson.D) (*Patient, error) {
	var result Patient
	err := ic.PatCol.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &result, nil
}

// FindPatientByID helper for finding patient by id easily
func (ic *IntouchClient) FindPatientByID(id primitive.ObjectID) (*Patient, error) {
	return ic.FindPatient(bson.D{{"_id", id}})
}

// IsUsernameInUse checks if username is already assigned to someone
// returns true if so, false otherwise
func (ic *IntouchClient) IsUsernameInUse(username string) bool {
	filter := bson.D{{"username", username}}
	inProviders, _ := ic.FindProvider(filter)
	inPatients, _ := ic.FindPatient(filter)
	if inProviders != nil || inPatients != nil {
		return true
	}
	return false

}

// InsertProvider inserts a provider with name, username, password and checks if username is valid
// returns nil id and error if problems arise, else id and nil
func (ic *IntouchClient) InsertProvider(name string, username string, password string) (*primitive.ObjectID, error) {
	if ic.IsUsernameInUse(username) {
		return nil, err["ErrUsernameInvalid"]
	}

	insertResult, err := ic.ProCol.InsertOne(context.TODO(), BlankUser{name, username, password})

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	if id, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		return &id, err
	}
	return nil, errors.New("Id type is not ObjectID")
}

// InsertPatient inserts a patient with name, username, password and checks if username is valid
// returns nil id and error if problems arise, else id and nil
func (ic *IntouchClient) InsertPatient(name string, username string, password string) (*primitive.ObjectID, error) {
	if ic.IsUsernameInUse(username) {
		return nil, err["ErrUsernameInvalid"]
	}

	insertResult, err := ic.PatCol.InsertOne(context.TODO(), BlankUser{name, username, password})

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	if id, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		return &id, err
	}
	return nil, errors.New("Id type is not ObjectID")
}

// DeleteProvider deletes given provider, doesn't check if exists
// will also update patients references to deleted provider
// returns error, nil if ok
func (ic *IntouchClient) DeleteProvider(id primitive.ObjectID) error {
	err := ic.deleteEntity(ic.ProCol, bson.D{{"_id", id}})
	if err != nil {
		return err
	}

	filter := bson.D{{"providers._id", id}}
	update := bson.D{{"$pull",
		bson.D{{"providers",
			bson.D{{"_id", id}}}}}}

	return ic.updateFieldAll(ic.PatCol, filter, update)
}

// DeletePatient deletes given patient, doesn't check if exists
// will also update provider references to deleted patient
// returns error, nil if ok
func (ic *IntouchClient) DeletePatient(id primitive.ObjectID) error {
	err := ic.deleteEntity(ic.PatCol, bson.D{{"_id", id}})
	if err != nil {
		return err
	}

	filter := bson.D{{"patients._id", id}}
	update := bson.D{{"$pull",
		bson.D{{"patients",
			bson.D{{"_id", id}}}}}}

	return ic.updateFieldAll(ic.ProCol, filter, update)
}

// AuthenticateProvider checks for a provider with given username and password
// if result found, returns _id else nil and error
func (ic *IntouchClient) AuthenticateProvider(username string, password string) (*primitive.ObjectID, error) {
	filter := bson.D{{"username", username}, {"password", password}}
	result, err := ic.FindProvider(filter)

	if err != nil {
		return nil, err
	}

	return &result.ID, nil
}

// AuthenticatePatient checks for a patient with given username and password
// if result found, returns _id else nil and error
func (ic *IntouchClient) AuthenticatePatient(username string, password string) (*primitive.ObjectID, error) {
	filter := bson.D{{"username", username}, {"password", password}}
	result, err := ic.FindPatient(filter)

	if result == nil {
		return nil, err
	}

	return &result.ID, nil
}

// UpdateProviderUsername checks if submitted username is in use, if not then updates provider
// and patient references to provider username.
// returns error, nil if ok
func (ic *IntouchClient) UpdateProviderUsername(id primitive.ObjectID, newusername string) error {
	if ic.IsUsernameInUse(newusername) {
		return err["ErrUsernameInvalid"]
	}

	// update provider with new username
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set",
		bson.D{{"username", newusername}}}}
	err := ic.updateField(ic.ProCol, filter, update)

	if err != nil {
		return err
	}

	// update patients referring to that provider
	// TODO: also update other collections
	filter = bson.D{{"providers._id", id}}
	update = bson.D{{"$set",
		bson.D{{"providers.$.username", newusername}}}}

	return ic.updateFieldAll(ic.PatCol, filter, update)
}

// UpdatePatientUsername checks if submitted username is in use, if not then updates patient
// and provider references to patient username.
// returns error, nil if ok
func (ic *IntouchClient) UpdatePatientUsername(id primitive.ObjectID, newusername string) error {
	if ic.IsUsernameInUse(newusername) {
		return err["ErrUsernameInvalid"]
	}

	// update patient with new username
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set",
		bson.D{{"username", newusername}}}}
	err := ic.updateField(ic.PatCol, filter, update)

	if err != nil {
		return err
	}

	// update providers referring to that patient
	// TODO: also update other collections
	filter = bson.D{{"patients._id", id}}
	update = bson.D{{"$set",
		bson.D{{"patients.$.username", newusername}}}}

	return ic.updateFieldAll(ic.ProCol, filter, update)
}

// AssociatePatient adds patient to provider's list of patients, as well as provider to
// target patient's list of providers
// returns error, nil if ok
func (ic *IntouchClient) AssociatePatient(proID primitive.ObjectID, patID primitive.ObjectID) error {
	provider, _ := ic.FindProvider(bson.D{{"_id", proID}})
	patient, _ := ic.FindPatient(bson.D{{"_id", patID}})

	if patient == nil {
		return err["ErrPatientNotFound"]
	}
	if provider == nil {
		return err["ErrProviderNotFound"]
	}

	filter := bson.D{{"_id", proID}}
	update := bson.D{{"$addToSet",
		bson.D{{"patients",
			bson.D{{"_id", patID},
				{"name", patient.Name},
				{"username", patient.Username}}}}}}

	err := ic.updateField(ic.ProCol, filter, update)

	if err != nil {
		return err
	}

	filter = bson.D{{"_id", patID}}
	update = bson.D{{"$addToSet",
		bson.D{{"providers",
			bson.D{{"_id", proID},
				{"name", provider.Name},
				{"username", provider.Username}}}}}}

	return ic.updateField(ic.PatCol, filter, update)
}

// DissociatePatient will remove patient and provider from each other's list
// returns error, nil if ok
func (ic *IntouchClient) DissociatePatient(proID primitive.ObjectID, patID primitive.ObjectID) error {
	filter := bson.D{{"_id", proID}}
	update := bson.D{{"$pull",
		bson.D{{"patients",
			bson.D{{"_id", patID}}}}}}
	err := ic.updateField(ic.ProCol, filter, update)

	if err != nil {
		return err
	}

	filter = bson.D{{"_id", patID}}
	update = bson.D{{"$pull",
		bson.D{{"providers",
			bson.D{{"_id", proID}}}}}}

	return ic.updateField(ic.PatCol, filter, update)
}

// BlankUser for inserting new document with random id rather than our own
type BlankUser struct {
	Name     string
	Username string
	Password string
}

// Provider if referenced by Patient then Salt, Password, and Patients have no values
type Provider struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string
	Username string
	Password string
	Salt     string
	Patients []Patient
}

// Patient if referenced by Provider then Salt, Password, and Patients have no values
type Patient struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string
	Username  string
	Password  string
	Salt      string
	Providers []Provider
}

// OpenConnection opens connection to mongo server and returns client
func OpenConnection() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Print(err)
		return nil
	}

	return client
}

// CreateIntouchClient uses mongo client to create collections for running queries
func CreateIntouchClient(db string, client *mongo.Client) *IntouchClient {
	providers := client.Database(db).Collection("providers")
	patients := client.Database(db).Collection("patients")

	return &IntouchClient{providers, patients}
}

func createObjID(parm string) primitive.ObjectID {
	result, _ := primitive.ObjectIDFromHex(parm)
	return result
}
