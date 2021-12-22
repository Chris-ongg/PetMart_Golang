package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go_server/petMart_server/model"
	"math/rand"
	"os"
	"strings"
	"sync"
)


func addOrderID(client *mongo.Client){
	orderIDCollection := client.Database(os.Getenv("DATABASE")).Collection("order_ID")
	filter := bson.D{{"name" , "counter"}}
	update := bson.D{{"$set", bson.D{{"latestid", 0}}}}
	opts := options.Update().SetUpsert(true)
	_ , err := orderIDCollection.UpdateOne(context.TODO() , filter , update , opts)
	if err != nil {
		panic(err)
	}
}

func addDummyCustomer(client *mongo.Client) {
	customerCollection := client.Database(os.Getenv("DATABASE")).Collection("customerAcc")
	filter := bson.D{{"EmailAdd" , "chris10@gmail.com"}}
	count , err := customerCollection.CountDocuments(context.TODO() , filter)
	if err != nil {
		panic(err)
	}

	if count == 0 {
		newCustomer := CustomerAccount{
			Name:        "Chris1",
			EmailAdd:    "chris10@gmail.com",
			Password:    "gAUWd4hDkRAerBYxYV4tKA==",
			Registered:  true,
			AccessToken: "",
			TokenExpiry: 0,
		}
		_, err = customerCollection.InsertOne(context.TODO(), newCustomer)
	}
}

func addCatProduct(channel chan *model.ProductList) {

	for i := 1; i <17 ; i++ {
		randStock := rand.Intn(50)
		randPrice := rand.Intn(100) + 20
		healthIssue := "Diabetes"
		var name strings.Builder
		fmt.Fprintf(&name , "Cat Food %d" , i + 16)
		var imagePath strings.Builder
		fmt.Fprintf(&imagePath , "../../Images/Product/%d.jpg" , i + 16)
		if i < 9 {
			healthIssue = "Weight Management"
		}

		temp := model.ProductList{
			ItemID:        i + 16,
			Name:          name.String() ,
			HealthConcern: healthIssue,
			Species:       "Cat",
			Price:         randPrice,
			Stock:         randStock,
			ImagePath:     imagePath.String(),
		}
		name.Reset()
		imagePath.Reset()
		channel <- &temp
	}
}

func addDogProduct(channel chan *model.ProductList) {
	for i := 1; i <17 ; i++ {
		randStock := rand.Intn(50)
		randPrice := rand.Intn(100) + 20
		healthIssue := "Diabetes"
		var name strings.Builder
		fmt.Fprintf(&name , "Dog Food %d" , i)
		var imagePath strings.Builder
		fmt.Fprintf(&imagePath , "../../Images/Product/%d.jpg" , i)
		if i < 9 {
			healthIssue = "Weight Management"
		}

		temp := model.ProductList{
			ItemID:        i ,
			Name:          name.String() ,
			HealthConcern: healthIssue,
			Species:       "Dog",
			Price:         randPrice,
			Stock:         randStock,
			ImagePath:     imagePath.String(),
		}
		name.Reset()
		imagePath.Reset()
		channel <- &temp
	}
}

func addProductsToWarehouse(client *mongo.Client) {
	warehouseCollection := client.Database(os.Getenv("DATABASE")).Collection("warehouse")
	docsCount , _ := warehouseCollection.CountDocuments(context.TODO(), bson.D{})

	if docsCount == 32 {
		return
	}

	productChannel := make(chan *model.ProductList)
	var wg sync.WaitGroup
	go func(channel chan *model.ProductList) {
		wg.Add(1)
		defer wg.Done()
		addDogProduct(channel)
	}(productChannel)
	go func(channel chan *model.ProductList) {
		wg.Add(1)
		defer wg.Done()
		addCatProduct(channel)
	}(productChannel)

	wg.Wait()

	var warehouse []interface{}
	for i := 0 ; i < 32 ; i++ {
		product := <- productChannel
		warehouse = append(warehouse, product)
	}
	close(productChannel)


	warehouseCollection.InsertMany(context.TODO() , warehouse)
}

func InitialiseDB(client *mongo.Client) {
	addProductsToWarehouse(client)
	addOrderID(client)
	addDummyCustomer(client)
}
