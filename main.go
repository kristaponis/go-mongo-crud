package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// connect to the MongoDB server, which is running local
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// list all database names in MongoDB server (check for connection)
	dbs, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dbs)

	// create database "bookstore"
	bookStoreDB := client.Database("bookstore")

	// create one collection (simple schema)
	booksColl := bookStoreDB.Collection("booksColl")
	_, err = booksColl.InsertOne(ctx, bson.D{
		{Key: "title", Value: "Golang programming"},
		{Key: "author", Value: "Peter Peterson"},
		{Key: "tags", Value: bson.A{"golang", "go programming", "go mongo"}},
	})
	if err != nil {
		log.Fatal(err)
	}
}
