package main

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var url = "mongodb://localhost:27017"

var db *mongo.Database
var once sync.Once

func GetClient() *mongo.Database {
	once.Do( func (){
		clientOpt := options.Client().ApplyURI(url)
		mongo, err := mongo.Connect(context.TODO(), clientOpt)

		if err != nil {
			log.Fatal(err)
		}

		err = mongo.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}
		db = mongo.Database("taskmanager")
	})
	return db
}
