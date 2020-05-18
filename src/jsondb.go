package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type JsonDB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func (db *JsonDB) Init() {
	log.Print("JsonDB Init")
	defer log.Print("Finish JsonDB Init")
	var err error
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	db.client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Print("Failed to connect to mongodb")
	}

	db.collection = db.client.Database("local").Collection("test")
}

func (db *JsonDB) Test() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := db.collection.InsertOne(ctx, bson.M{"name": "Radu", "last_name": "Pavel", "cnp": 123123})
	if err != nil {
		log.Print(err)
		return
	}
	id := res.InsertedID

	log.Print(id)
}
