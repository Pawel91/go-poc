package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//PersonDB interface
type IPersonDB interface {
	Insert(*Person) error
	Get(int) (*Person, error)
}

type personDB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

//factory function
func NewPersonDB() (IPersonDB, error) {
	var i IPersonDB
	implementation := &personDB{}
	implementation.init()
	i = implementation
	return i, nil
}

func (db *personDB) Insert(person *Person) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	obj := bson.M{"name": person.Name, "lastname": person.LastName, "cnp": person.CNP}
	_, err := db.collection.InsertOne(ctx, obj)
	if err != nil {
		return err
	}

	return nil
}

func (db *personDB) Get(cnp int) (*Person, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	value := db.collection.FindOne(ctx, bson.M{"cnp": cnp})

	var person Person
	err := value.Decode(&person)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return &person, nil
}

func (db *personDB) init() {
	log.Print("PersonDB Init")
	defer log.Print("Finish PersonDB Init")
	var err error
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	db.client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Print("Failed to connect to mongodb")
	}

	db.collection = db.client.Database("local").Collection("test")
}
