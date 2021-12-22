package implementation

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go_server/petMart_server/database"
	"go_server/petMart_server/dependencies/interfaces"
	"go_server/petMart_server/model"
	"os"
	"strconv"
	"time"
)

type MutationService struct {
	db *mongo.Client
}

func NewMutationService(db *mongo.Client) *MutationService{
	return &MutationService{
		db: db,
	}
}

var _ interfaces.MutationServices = &MutationService{}

func getCustomerInCustomerAcc(collection *mongo.Collection , email string , channel chan database.CustomerAccount) {
	var customer database.CustomerAccount
	filter := bson.D{{"emailadd" , email}}
	err := collection.FindOne(context.TODO() , filter).Decode(&customer)
	if err != nil {
		fmt.Println(err)
	}
	channel <- customer
}

func getCustomerInCustomerActivity(collection *mongo.Collection , email string , channel chan database.CustomerActivity) {
	var customer database.CustomerActivity
	filter := bson.D{{"email" , email}}
	err := collection.FindOne(context.TODO() , filter).Decode(&customer)
	if err != nil {
		fmt.Println(err)
	}
	channel <- customer
}

func getLastOrderID(collection *mongo.Collection , channel chan int) {
	var latestID database.OrderID
	filter := bson.D{{"name" , "counter"}}
	err := collection.FindOne(context.TODO() , filter).Decode(&latestID)
	if err != nil {
		fmt.Println(err)
	}
	channel <- latestID.LatestID
}

func (M *MutationService) RegisterNewCustomer(name string, email string , password string) error{

	customerCollection := M.db.Database(os.Getenv("DATABASE")).Collection("customerAcc")
	filter := bson.D{{"EmailAdd" , email}}

	count , err := customerCollection.CountDocuments(context.TODO() , filter)

	if err != nil {
		return err
	}

	if count == 0 {
		newCustomer := database.CustomerAccount{
			Name:        name,
			EmailAdd:    email,
			Password:    password,
			Registered:  true,
			AccessToken: "",
			TokenExpiry: 0,
		}
		_, err = customerCollection.InsertOne(context.TODO(), newCustomer)
	}

	return err
}

func (M *MutationService) LogoutCustomer(email string) string {
	customerCollection := M.db.Database(os.Getenv("DATABASE")).Collection("customerAcc")

	customerChan := make(chan database.CustomerAccount , 1)
	go getCustomerInCustomerAcc(customerCollection , email , customerChan)

	filter := bson.D{{"emailadd" , email}}
	tokenValidity := int64(0)
	update := bson.D{{"$set", bson.D{{"accesstoken", ""} , {"tokenexpiry" , tokenValidity}}}}
	opts := options.Update().SetUpsert(false)
	_ , err := customerCollection.UpdateOne(context.TODO(), filter , update , opts )
	if err != nil {
		panic(err)
	}

	customer := <-customerChan
	return customer.Name
}

func (M *MutationService) RegisterNewPets(newPet *model.PetDetails) error {
	petCollection := M.db.Database(os.Getenv("DATABASE")).Collection("customerPets")

	filter := bson.D{{"owneremail" , newPet.OwnerEmail} , {"name"  , newPet.Name} , {"species", newPet.Species}}

	count , err := petCollection.CountDocuments(context.TODO() , filter)
	if err != nil {
		return err
	}
	if count == 0 {
		newCustomerPet := database.CustomerPets{
			Name:          newPet.Name,
			Gender:        newPet.Gender,
			PetBreed:      newPet.PetBreed,
			Species:       newPet.Species,
			Weight:        newPet.Weight,
			OwnerEmail:    newPet.OwnerEmail,
			HealthConcern: newPet.HealthConcern,
			Age:           newPet.Weight,
		}
		_, err = petCollection.InsertOne(context.TODO(), newCustomerPet)
	}
	return err
}

func (M *MutationService) SaveShoppingCart(currentCart *model.CustomerCart) error {
	activityCollection := M.db.Database(os.Getenv("DATABASE")).Collection("customerActivity")

	filter := bson.D{{"email" , currentCart.Email}}
	update := bson.D{{"$set", bson.D{{"emptycart", currentCart.EmptyCart} , {"cartitems" , currentCart.Cart}}}}
	opts := options.Update().SetUpsert(true)
	_ , err := activityCollection.UpdateOne(context.TODO(), filter , update , opts )
	if err != nil {
		panic(err)
	}

	return err
}

func (M *MutationService) SaveOrder(purchasedCart *model.CustomerCart ) error {
	activityCollection := M.db.Database(os.Getenv("DATABASE")).Collection("customerActivity")
	counterCollection := M.db.Database(os.Getenv("DATABASE")).Collection("order_ID")

	counterChannel := make(chan int , 1)
	go getLastOrderID( counterCollection, counterChannel)

	customerChannel := make(chan database.CustomerActivity , 1)
	go getCustomerInCustomerActivity(activityCollection, purchasedCart.Email, customerChannel)

	itemsCount := len(purchasedCart.Cart)
	totalPrice := 0
	for _ , cart := range purchasedCart.Cart {
		totalPrice += cart.Price
	}

	date := time.Now().Format("2006.01.02 15:04:05")
	trackingID := strconv.Itoa(int(time.Now().Unix()))

	orderID := <-counterChannel

	newOrder := model.PastOrders{
		Date:       date,
		OrderID:    orderID + 1,
		Store:      "Online",
		Items:      itemsCount,
		Total:      totalPrice,
		TrackingID: trackingID,
		Status:     "Completed",
	}

	customer := <-customerChannel
	var allPastOrders []*model.PastOrders

	allPastOrders = append(allPastOrders , &newOrder)
	allPastOrders = append(allPastOrders , customer.PastOrders...)

	var EmptyCart []*model.CartItemsInput

	filter := bson.D{{"email" , purchasedCart.Email}}
	update := bson.D{{"$set", bson.D{{"emptycart", true} , {"cartitems" , EmptyCart} , {"pastorders" , allPastOrders}}}}
	opts := options.Update().SetUpsert(true)
	_ , err := activityCollection.UpdateOne(context.TODO() , filter , update , opts)
	if err != nil {
		panic(err)
		return err
	}

	filter = bson.D{{"name" , "counter"}}
	update = bson.D{{"$set", bson.D{{"latestid", orderID + 1}}}}
	_ , err2 := counterCollection.UpdateOne(context.TODO() , filter , update)

	return err2
}