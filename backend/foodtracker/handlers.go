package foodtracker

import (
	"FoodTracker/user_management"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"nilswilhelm.net/foodtracker/lib/constants"
	"nilswilhelm.net/foodtracker/lib/structs"
	"time"
)

const dateLayout = "2006-01-02"

// TODO: Deprecated => remove
func FoodHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getFoodHandler(w, r)
		break
	case http.MethodPost:
		postFoodHandler(w, r)
		break
	}
}

// Main Handlers
func GenericFoodHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postFoodHandler(w, r)
		break
	}
}

func IdentifiedFoodResourceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getFoodHandler(w, r)
		break
	case http.MethodDelete:
		deleteFoodHandler(w, r)
		break
	case http.MethodPut:
		updateFoodHandler(w, r)
		break
	}
}

func IntakeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("IntakeHandler called")
	switch r.Method {
	case http.MethodGet:
		getIntakeHandler(w, r)
		break
	case http.MethodPost:
		addToIntakeHandler(w, r)
		break
	}
}

// Sub handlers
func getIntakeHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	today := time.Now().Format(constants.DateFormat)
	intake, err := getIntake(today, userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println(intake)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(intake)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func addToIntakeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("addToIntakeHandler")
	userId := r.Header.Get("userId")
	println(userId)
	var transaction structs.Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		fmt.Println("Error while decoding json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = addToIntake(&transaction, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// transaction.Date = time.Now()
	transaction.UserId = userId
	log.Println(transaction)
	err = insertTransaction(transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// RemoveFromIntakeHandler removes an transaction entry and resets the intake values
func RemoveFromIntakeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionId := vars["transactionId"]
	userId := r.Header.Get("userId")
	err := rollbackTransaction(transactionId, userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = deleteTransaction(transactionId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(transactionId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// HistoryHandler gets the calorie history of the last days for a given user
func HistoryHandler(writer http.ResponseWriter, request *http.Request) {
	userId := request.Header.Get("userId")
	history, err := getHistory(userId)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	log.Printf("History: %v", history)
	err = json.NewEncoder(writer).Encode(history)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

// getFoodHandler gets food from the database
func getFoodHandler(w http.ResponseWriter, r *http.Request) {
	var food structs.Food
	vars := mux.Vars(r)
	id := vars["id"]
	params := r.URL.Query()
	if params.Get("ean") == "true" {
		food, err := getFoodByEAN(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(food)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		oid, err := primitive.ObjectIDFromHex(id)
		food, err = getFood(oid)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(food)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func deleteFoodHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := deleteFood(id)
	if err != nil {
		http.Error(w, "could not delete food", http.StatusInternalServerError)
	}
	w.WriteHeader(200)
}

// postFoodHandler writes new food into the database and search index
func postFoodHandler(w http.ResponseWriter, r *http.Request) {

	var newFood structs.Food
	println("Post Food")
	err := json.NewDecoder(r.Body).Decode(&newFood)
	if err != nil {
		fmt.Println("Error while decoding json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//log.Println(newFood)
	//if newFood.IsMeal {
	//	newFood.Nutrition = structs.Nutrition{
	//		Energy:       0,
	//		Fat:          0,
	//		Carbohydrate: 0,
	//		Protein:      0,
	//	}
	//	for _, ingredient := range newFood.Ingredients {
	//		log.Printf("ID: %s", ingredient.Food.ID.String())
	//		food, err := getFood(*ingredient.Food.ID)
	//		if err != nil {
	//			w.WriteHeader(500)
	//			return
	//		}
	//		newFood.Nutrition.Carbohydrate += food.Nutrition.Carbohydrate * ingredient.Amount / 100
	//		newFood.Nutrition.Energy += food.Nutrition.Energy * ingredient.Amount / 100
	//		newFood.Nutrition.Fat += food.Nutrition.Fat * ingredient.Amount / 100
	//		newFood.Nutrition.Protein += food.Nutrition.Protein * ingredient.Amount / 100
	//	}
	//}

	err = storeFood(newFood)
	if err != nil {
		log.Printf("Could not store food %v", newFood)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := newFood
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Could not encode %v to json", response)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func updateFoodHandler(w http.ResponseWriter, r *http.Request) {
	var foodToUpdate structs.Food
	err := json.NewDecoder(r.Body).Decode(&foodToUpdate)
	if err != nil {
		fmt.Println("Error while decoding json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(foodToUpdate)

	err = updateFood(foodToUpdate)
	if err != nil {
		log.Printf("Could not update food %v", foodToUpdate)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := foodToUpdate
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Could not encode %v to json", response)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// RecalcHandler triggers recalculation for intake of specific date
func RecalcHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	date, err := time.Parse(constants.DateFormat, vars["date"])
	if err != nil {
		log.Printf("Date %s could not be parsed", vars["date"])
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	userId := request.Header.Get("userId")

	nutrition, err := recalcIntake(userId, date)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(writer).Encode(nutrition)
	if err != nil {
		log.Printf("Could not encode %v to json", nutrition)
		writer.WriteHeader(http.StatusInternalServerError)
	}

}

// TransactionHandler gets all food entries of today for a given user
func GenericTransactionHandler(writer http.ResponseWriter, request *http.Request) {
	date := time.Now()
	dateString := request.URL.Query().Get("date")
	if dateString != "" {
		parsedDate, err := time.Parse(dateLayout, dateString)
		if err != nil {
			http.Error(writer, "could not parse date", http.StatusBadRequest)
		}
		date = parsedDate
	}
	userId := request.Header.Get("userId")
	transactions, err := getTransactions(userId, date)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	if transactions == nil {
		transactions = make([]structs.Transaction, 0)
	}
	err = json.NewEncoder(writer).Encode(transactions)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

// TransactionHandler gets all food entries of today for a given user
func IdentifiedTransactionResourceHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodDelete:
		deleteTransactionHandler(writer, request)
	}
}

func deleteTransactionHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	transactionId := vars["id"]
	err := deleteTransaction(transactionId)
	if err != nil {
		http.Error(writer, "could not delete transaction", http.StatusInternalServerError)
	}
	writer.WriteHeader(http.StatusOK)
}

// TODO deprecated => remove
func TransactionHandler(writer http.ResponseWriter, request *http.Request) {
	userId := request.Header.Get("userId")
	transactions, err := getTransactions(userId, time.Now())
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	err = json.NewEncoder(writer).Encode(transactions)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

// UsersHandler returns specific user
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	user, err := user_management.GetUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	userResponse := user_management.UserResponse{Username: user.Username}
	err = json.NewEncoder(w).Encode(userResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// SearchHandler handlers search requests and returns matching food
func SearchHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	name := vars["name"]

	foods, err := search(name)
	if err != nil {
		log.Println("Search failed")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(foods)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func DailyGoalsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("IntakeHandler called")
	switch r.Method {
	case http.MethodGet:
		getDailyGoalsHandler(w, r)
		break
	case http.MethodPost:
		postDailyGoalsHandler(w, r)
		break
	}
}

func getDailyGoalsHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	//today := time.Now().Format(DateFormat)
	dailyGoals, err := getDailyGoals(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println(dailyGoals)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(dailyGoals)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func postDailyGoalsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("addToIntakeHandler")
	userId := r.Header.Get("userId")
	println(userId)
	var dailyGoals structs.DailyGoals
	err := json.NewDecoder(r.Body).Decode(&dailyGoals)
	if err != nil {
		fmt.Println("Error while decoding json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dailyGoals.UserId = userId
	_, err = getDailyGoals(userId)
	if err != nil {
		err = insertDailyGoals(dailyGoals)
	}

	err = updateDailyGoals(dailyGoals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(dailyGoals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func SimpleOk(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
	_, err := writer.Write([]byte("Works"))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func updateDailyGoals(dailyGoals structs.DailyGoals) error {

	filter := bson.D{{"userId", dailyGoals.UserId}}
	update := bson.D{
		{"$set", bson.D{
			{"nutrition.energy", dailyGoals.Nutrition.Energy},
		}},
		{"$set", bson.D{
			{"nutrition.fat", dailyGoals.Nutrition.Fat},
		}},
		{"$set", bson.D{
			{"nutrition.carbohydrate", dailyGoals.Nutrition.Carbohydrate},
		}},
		{"$set", bson.D{
			{"nutrition.protein", dailyGoals.Nutrition.Protein},
		}},
		{"$set", bson.D{
			{"water", dailyGoals.Water},
		}},
	}

	err := structs.MongoClient.Update(filter, update, constants.DbName, constants.DAILYGOALS)
	if err != nil {
		return err
	}
	return nil
}

// DashboardDataHandler gets all food entries of today for a given user
func DashboardDataHandler(writer http.ResponseWriter, request *http.Request) {
	var date time.Time
	var err error
	dateString := request.URL.Query().Get("date")
	if dateString == "" {
		date = time.Now()
	} else {
		date, err = time.Parse(constants.DateFormat, dateString)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
	}
	userId := request.Header.Get("userId")
	transactions, err := getTransactions(userId, date)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	goals, err := getDailyGoals(userId)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	intake, err := getIntake(date.Format(constants.DateFormat), userId)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	dashboardData := structs.DashboardData{
		Intake:       intake,
		DailyGoals:   goals,
		Transactions: transactions,
	}

	err = json.NewEncoder(writer).Encode(dashboardData)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
