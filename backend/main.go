package main

import (
	"FoodTracker/foodtracker"
	"FoodTracker/user_management"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"nilswilhelm.net/foodtracker/lib/database"
	"nilswilhelm.net/foodtracker/lib/structs"
	"os"
)

//Environment Variables
var ADDRESS string
var MongoUri string

func main() {
	initVariables()

	mongoClient, err := database.NewMongoAdapter(MongoUri)
	if err != nil {
		log.Fatal(err)
	}
	structs.MongoClient = mongoClient
	user_management.MongoClient = mongoClient

	handleRequests()

}

func handleRequests() {
	log.Println("handle requests")
	router := mux.NewRouter()
	router.HandleFunc("/", foodtracker.SimpleOk)
	router.Handle("/intake", user_management.IsAuthorized(foodtracker.IntakeHandler))
	router.Handle("/dashboard", user_management.IsAuthorized(foodtracker.DashboardDataHandler))
	router.Handle("/intake/{transactionId}", user_management.IsAuthorized(foodtracker.RemoveFromIntakeHandler))
	router.Handle("/food", user_management.IsAuthorized(foodtracker.GenericFoodHandler))
	router.Handle("/food/{id}", user_management.IsAuthorized(foodtracker.IdentifiedFoodResourceHandler))
	// TODO: Remove /foodlist/{id}
	router.HandleFunc("/foodlist/{id}", foodtracker.FoodHandler)
	router.Handle("/history", user_management.IsAuthorized(foodtracker.HistoryHandler))
	router.Handle("/transactions", user_management.IsAuthorized(foodtracker.GenericTransactionHandler))
	router.Handle("/transactions/{id}", user_management.IsAuthorized(foodtracker.IdentifiedTransactionResourceHandler))
	router.Handle("/goals", user_management.IsAuthorized(foodtracker.DailyGoalsHandler))
	router.Handle("/recalc/{date}", user_management.IsAuthorized(foodtracker.RecalcHandler))
	router.HandleFunc("/login", user_management.MakeLoginHandler())
	router.HandleFunc("/register", user_management.MakeSignUpHandler())
	router.HandleFunc("/search/{name}", foodtracker.SearchHandler)

	// Allow CORS
	handler := cors.AllowAll().Handler(router)
	log.Fatal(http.ListenAndServe(ADDRESS, handler))

}

func initVariables() {
	log.Println("init variables")
	ADDRESS = os.Getenv("ADDRESS")
	if ADDRESS == "" {
		ADDRESS = ":9000"
	}
	MongoUri = os.Getenv("MONGO_URI")
	if MongoUri == "" {
		MongoUri = "mongodb://localhost:27017"
	}
	SigningKeyString := os.Getenv("SIGNING_KEY")
	if SigningKeyString == "" {
		user_management.SigningKey = []byte("supersecret")
	} else {
		user_management.SigningKey = []byte(SigningKeyString)
	}
}
