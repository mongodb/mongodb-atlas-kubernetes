package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// simple test application
// allows to check if users can use the provided connectionstring,
// and if data remains after changes are made

var (
	collection *mongo.Collection
	ctx        = context.TODO()

	// related to DB
	dbName         = "Ships"
	collectionName = "ShipSpec"
	portDefault    = "8080"
)

func main() {
	r := newRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = portDefault
	}
	fmt.Print("Using port: " + port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("GET")
	r.HandleFunc("/mongo/{key}", getKeyValue).Methods("GET")
	r.HandleFunc("/mongo/", postKeyValue).Methods("POST")
	r.HandleFunc("/mongo/{key}", deleteKeyValue).Methods("DELETE")
	return r
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "It is working")
}

func getKeyValue(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	collection, err := getMongoCollection(dbName, collectionName)
	if err != nil {
		log.Println(err)
	}
	filter := bson.M{"key": key}
	var foundShip ModelShip
	collection.FindOne(ctx, filter).Decode(&foundShip)
	res, _ := json.Marshal(&foundShip)
	fmt.Fprintf(w, string(res))
}

func postKeyValue(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Lets see\n")
	collection, err := getMongoCollection(dbName, collectionName)
	if err != nil {
		log.Println(err)
	}

	ship := getShipFromRequest(r)
	fmt.Fprintf(w, "got ship key: %s\n", ship.Key)
	res, err := collection.InsertOne(ctx, ship)
	if err != nil {
		log.Println(err)
	}

	response, _ := json.Marshal(res)
	fmt.Fprintf(w, string(response))
}

func deleteKeyValue(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	collection, err := getMongoCollection(dbName, collectionName)
	if err != nil {
		log.Println(err)
	}
	deleteResult, err := collection.DeleteMany(ctx, bson.M{"key": key})
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, fmt.Sprintf("It is working. Deleted documents: %d", deleteResult.DeletedCount))
}

func getSecret() string {
	if val, ok := os.LookupEnv("connectionStringStandardSrv"); ok {
		return val
	}
	// Before the Operator 0.6.0 is released we need to "failback" to the previous env variable name in our tests
	return os.Getenv("connectionString.standardSrv")
}

func getMongoClient() (*mongo.Client, error) {
	uri := getSecret()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Println(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println(err)
	}
	fmt.Println("Successfully connected and pinged.")
	return client, nil
}

func getMongoCollection(DbName string, CollectionName string) (*mongo.Collection, error) {
	client, err := getMongoClient()
	if err != nil {
		return nil, err
	}

	collection := client.Database(DbName).Collection(CollectionName)

	return collection, nil
}

type ModelShip struct {
	Key       string `json:"key"`
	ShipModel string `json:"shipmodel,omitempty"`
	Hp        int    `json:"hp,omitempty"`
}

func getShipFromRequest(r *http.Request) ModelShip {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var ship ModelShip
	json.Unmarshal(reqBody, &ship)
	return ship
}
