package implementation

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go_server/petMart_server/database"
	"go_server/petMart_server/dependencies/interfaces"
	"go_server/petMart_server/model"
	"google.golang.org/api/idtoken"
	"os"
	"reflect"
	"time"
)

type QueryService struct {
	db *mongo.Client
}

func NewQueryService(db *mongo.Client) *QueryService {
	return &QueryService{
		db: db,
	}
}

func decodePassword(encryptedPass string) string{
	cipherText , _ := base64.StdEncoding.DecodeString(encryptedPass)
	block , _ := aes.NewCipher([]byte(os.Getenv("SECRETKEY")))
	if len(cipherText)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, []byte(os.Getenv("IV")))
	mode.CryptBlocks(cipherText, cipherText)
	return string(cipherText)
}

var _ interfaces.QueryServices = &QueryService{}

func (Q *QueryService) CustomerPets(email string) ([]*model.PetDetails ,error) {
	var pets []*model.PetDetails

	petsCollection := Q.db.Database(os.Getenv("DATABASE")).Collection("customerPets")

	filter := bson.D{{"owneremail" , email}}

	cursor , err := petsCollection.Find(context.TODO() , filter)

	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.TODO(), &pets); err != nil {
		panic(err)
	}

	return pets, nil
}

func (Q *QueryService) EmailLogin(email string, password string) (string , string) {

	customerCollection := Q.db.Database(os.Getenv("DATABASE")).Collection("customerAcc")
	var customer database.CustomerAccount
	filter := bson.D{{"emailadd" , email}}
	err := customerCollection.FindOne(context.TODO() , filter).Decode(&customer)
	if err != nil {
		fmt.Println(err)
		return "fail" , ""
	}

	passwordDatabase := decodePassword(customer.Password)
	passwordLogin := decodePassword(password)

	if passwordDatabase == passwordLogin {
		accessToken := uuid.New()
		//Token is valid for 7 days
		tokenValidity := (time.Now().UnixNano() / 1e6) + int64(7 * 24 * 60 * 60 * 1000)
		update := bson.D{{"$set", bson.D{{"accesstoken", accessToken.String()} , {"tokenexpiry" , tokenValidity}}}}
		opts := options.Update().SetUpsert(false)
		result , err := customerCollection.UpdateOne(context.TODO(), filter , update , opts )
		if err != nil {
			panic(err)
		}
		fmt.Printf("Number of documents updated: %v\n\n", result.ModifiedCount)
		return accessToken.String() , customer.Name
	}
	return "fail", ""
}

func (Q *QueryService) GoogleLogin(accessToken string) (string , string , string) {
	customerCollection := Q.db.Database(os.Getenv("DATABASE")).Collection("customerAcc")
	payload, err := idtoken.Validate(context.Background(), accessToken, os.Getenv("GOOGLECLIENTID"))
	if err != nil {
		panic(err)
	}
	token := uuid.New()
	//Token is valid for 7 days
	tokenValidity := (time.Now().UnixNano() / 1e6) + int64(7 * 24 * 60 * 60 * 1000)
	email := fmt.Sprintf( "%v" , payload.Claims["email"])
	name  := fmt.Sprintf("%v" ,payload.Claims["name"])
	fmt.Println(name , email)
	update := bson.D{{"$set", bson.D{
		{"accesstoken", token.String()} ,
		{"tokenexpiry" , tokenValidity}  ,
		{"name" , name },
		{"emailadd" , email} ,
		{"password" , ""} ,
		{"registered" , false}}}}
	filter := bson.D{{"emailadd" , email}}
	opts := options.Update().SetUpsert(true)
	result , err := customerCollection.UpdateOne(context.TODO(), filter , update , opts )
	if err != nil {
		panic(err)
	}
	fmt.Printf("Number of documents updated: %v\n\n", result.ModifiedCount)
	return name ,  email , token.String()
}

func (Q *QueryService) GetShoppingCart(email string) ([]*model.CartItemsInput , error) {
	activityCollection := Q.db.Database(os.Getenv("DATABASE")).Collection("customerActivity")
	var result database.CustomerActivity
	filter := bson.D{{"email" , email}}
	err := activityCollection.FindOne(context.TODO() , filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}
	if result.EmptyCart {
		return nil , nil
	}
	return result.CartItems , nil
}

func (Q *QueryService) TokenVerification(token string)(*model.AuthenticatedUser , error) {
	currentTime := time.Now().UnixNano() / 1e6
	customerCollection := Q.db.Database(os.Getenv("DATABASE")).Collection("customerAcc")
	var customer database.CustomerAccount
	filter := bson.D{{"accesstoken" , token}}
	err := customerCollection.FindOne(context.TODO() , filter).Decode(&customer)
	if !reflect.ValueOf(customer).IsZero() {
		if customer.TokenExpiry > currentTime {
			authenticatedUser := model.AuthenticatedUser{
				Name:        customer.Name,
				Email:       customer.EmailAdd,
				AccessToken: customer.AccessToken,
			}
			return &authenticatedUser, nil
		}
	}
	return nil , err
}

func (Q *QueryService) GetCustomerPastOrders(email string)([]*model.PastOrders ,  error) {
	activityCollection := Q.db.Database(os.Getenv("DATABASE")).Collection("customerActivity")
	var result database.CustomerActivity
	filter := bson.D{{"email" , email}}
	err := activityCollection.FindOne(context.TODO() , filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}
	return result.PastOrders , err
}

func (Q *QueryService) GetProductsFromWarehouse(searchType int , health string, species string) ([]*model.ProductList , error) {
	warehouseCollection := Q.db.Database(os.Getenv("DATABASE")).Collection("warehouse")
	filter := bson.D{}
	if !(searchType == 0) {
		if health == "No Health Concern" {
			filter = bson.D{{"species" , species}}
		}else {
			filter = bson.D{{"species" , species} , {"healthconcern" , health}}
		}
	}
	var retrievedProducts []*model.ProductList
	cursor , err := warehouseCollection.Find(context.TODO() , filter)
	if err = cursor.All(context.TODO(), &retrievedProducts); err != nil {
		panic(err)
	}
	return retrievedProducts , err
}
