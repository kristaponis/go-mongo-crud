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

	// create collection with insertOne()
	booksColl := bookStoreDB.Collection("booksColl")
	_, err = booksColl.InsertOne(ctx, bson.D{
		{Key: "title", Value: "Golang programming"},
		{Key: "author", Value: "Peter Peterson"},
		{Key: "tags", Value: bson.A{"golang", "go programming", "go mongo"}},
	})
	if err != nil {
		log.Fatal(err)
	}

	// create collection with insertMany()
	journalColl := bookStoreDB.Collection("journalColl")
	journalResult, err := journalColl.InsertMany(ctx, []interface{}{
		bson.D{
			{Key: "title", Value: "Daily Journal"},
			{Key: "description", Value: "It is a daily journal to your desk"},
			{Key: "price", Value: 3.99},
		},
		bson.D{
			{Key: "title", Value: "Go blog"},
			{Key: "description", Value: "Golang blog with useful info"},
			{Key: "price", Value: 2.49},
		},
		bson.D{
			{Key: "title", Value: "Programmer"},
			{Key: "description", Value: "All about programming"},
			{Key: "price", Value: 5.99},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(journalResult.InsertedIDs)

	// read all the documents in journal collection
	cursor, err := journalColl.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	var journals []bson.M
	if err := cursor.All(ctx, &journals); err != nil {
		log.Fatal(err)
	}
	for _, j := range journals {
		fmt.Println(j)
	}
}
