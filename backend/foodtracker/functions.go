package foodtracker

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	. "nilswilhelm.net/foodtracker/lib/constants"
	. "nilswilhelm.net/foodtracker/lib/structs"
	"time"
)

func rollbackTransaction(transactionId string, userId string) error {
	oid, err := primitive.ObjectIDFromHex(transactionId)
	filter := bson.D{{"_id", oid}}

	result := MongoClient.Get(filter, DbName, TRANSACTIONS)

	var transaction Transaction
	err = result.Decode(&transaction)
	if err != nil {
		return err
	}
	fmt.Println("TransactionID:")
	fmt.Println(transaction.ID)

	updateEnergy := -transaction.Nutrition.Energy
	updateFat := -transaction.Nutrition.Fat
	updateCarbs := -transaction.Nutrition.Carbohydrate
	updateProtein := -transaction.Nutrition.Protein
	fmt.Println("Updates")
	fmt.Println(updateEnergy)
	fmt.Println(updateFat)
	fmt.Println(updateCarbs)
	fmt.Println(updateProtein)
	fmt.Println("Updates end")

	date := transaction.Date.Format(DateFormat)
	filter = bson.D{{"date", date}, {"userid", userId}}
	update := bson.D{
		{"$inc", bson.D{
			{"nutrition.energy", updateEnergy},
		}},
		{"$inc", bson.D{
			{"nutrition.fat", updateFat},
		}},
		{"$inc", bson.D{
			{"nutrition.carbohydrate", updateCarbs},
		}},
		{"$inc", bson.D{
			{"nutrition.protein", updateProtein},
		}},
	}
	err = MongoClient.Update(filter, update, DbName, IntakeHistory)
	if err != nil {
		return err
	}

	return nil
}
func recalcIntake(userId string, date time.Time) (Nutrition, error) {
	transactions, err := getTransactions(userId, date)
	if err != nil {
		log.Printf("Could not get transactions for %s \n", date)
		return Nutrition{}, err
	}
	nut := Nutrition{}

	for _, transaction := range transactions {
		nut.Energy += transaction.Nutrition.Energy
		nut.Protein += transaction.Nutrition.Protein
		nut.Carbohydrate += transaction.Nutrition.Carbohydrate
		nut.Fat += transaction.Nutrition.Fat
	}
	err = setIntake(userId, nut, date)
	if err != nil {
		log.Printf("Could not set intake for %s \n", date)
		return Nutrition{}, err
	}
	return nut, nil
}

func setIntake(userId string, nutrition Nutrition, date time.Time) error {

	filter := bson.D{{"date", date.Format(DateFormat)}, {"userid", userId}}
	update := bson.D{
		{"$set", bson.D{
			{"nutrition.energy", nutrition.Energy},
		}},
		{"$set", bson.D{
			{"nutrition.fat", nutrition.Fat},
		}},
		{"$set", bson.D{
			{"nutrition.carbohydrate", nutrition.Carbohydrate},
		}},
		{"$set", bson.D{
			{"nutrition.protein", nutrition.Protein},
		}},
	}

	err := MongoClient.Update(filter, update, DbName, IntakeHistory)
	if err != nil {
		return err
	}
	return nil
}

func insertTransaction(transaction Transaction) error {
	_, err := MongoClient.Insert(transaction, DbName, TRANSACTIONS)
	if err != nil {
		return err
	}
	return nil
}

func deleteTransaction(ID string) error {
	oid, err := primitive.ObjectIDFromHex(ID)
	filter := bson.D{{"_id", oid}}
	err = MongoClient.Delete(filter, DbName, TRANSACTIONS)
	if err != nil {
		return err
	}
	return nil
}

func getFood(ID primitive.ObjectID) (Food, error) {
	//oid, _ := primitive.ObjectIDFromHex(ID)
	var f Food
	filter := bson.D{{"_id", ID}}
	result := MongoClient.Get(filter, DbName, FoodList)
	err := result.Decode(&f)
	if err != nil {
		return Food{}, err
	}
	fmt.Printf("Got food: %s\n", f.Name)
	return f, nil
}

func getFoodByEAN(EAN string) (Food, error) {
	var f Food
	filter := bson.D{{"ean", EAN}}
	result := MongoClient.Get(filter, DbName, FoodList)
	err := result.Decode(&f)
	if err != nil {
		return Food{}, err
	}
	fmt.Printf("Got food: %s\n", f.Name)
	return f, nil
}

func getIntake(date string, userId string) (Intake, error) {
	var n Intake
	filter := bson.D{{"date", date}, {"userid", userId}}
	result := MongoClient.Get(filter, DbName, IntakeHistory)
	err := result.Decode(&n)
	if err != nil {
		err = insertEmptyIntake(date, userId)
		if err != nil {
			return Intake{}, err
		}
		return Intake{
			UserId:    userId,
			Date:      date,
			Nutrition: Nutrition{},
		}, nil
	}
	fmt.Printf("Got nutrition for: %s\n", n.Date)
	return n, nil
}

func insertEmptyIntake(date string, userId string) error {
	println("Insert empty")
	intake := Intake{
		Date:   date,
		UserId: userId,
	}
	_, err := MongoClient.Insert(intake, DbName, IntakeHistory)
	if err != nil {
		println("Can not add new intake")
		return err
	}
	return nil
}

func storeFood(food Food) error {
	_, err := MongoClient.Insert(food, DbName, FoodList)
	if err != nil {
		return err
	}
	return nil
}

func updateFood(food Food) error {
	log.Println(food.Nutrition.Energy)
	filter := bson.D{{"id", food.ID}}
	update := bson.D{
		{"$set", bson.D{
			{"name", food.Name},
		}},
		{"$set", bson.D{
			{"ean", food.EAN},
		}},
		{"$set", bson.D{
			{"brand", food.Brand},
		}},
		{"$set", bson.D{
			{"isMeal", food.IsMeal},
		}},
		// TODO update ingredients and meta fields
		{"$set", bson.D{
			{"nutrition.energy", food.Nutrition.Energy},
		}},
		{"$set", bson.D{
			{"nutrition.fat", food.Nutrition.Fat},
		}},
		{"$set", bson.D{
			{"nutrition.carbohydrate", food.Nutrition.Carbohydrate},
		}},
		{"$set", bson.D{
			{"nutrition.protein", food.Nutrition.Protein},
		}},
	}

	return MongoClient.Update(filter, update, DbName, FoodList)
}

func deleteFood(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", objID}}
	return MongoClient.Delete(filter, DbName, FoodList)
}

func addToIntake(transaction *Transaction, userId string) error {
	log.Println("addToIntake")
	log.Println(transaction)
	var food Food
	var factor float64
	dateString := transaction.Date.Format(DateFormat)

	_, err := getIntake(dateString, userId)
	if err != nil {
		println("Can not add new intake")
		return err

	}
	oid, err := primitive.ObjectIDFromHex(transaction.FoodID)
	food, err = getFood(oid)
	if err != nil {
		println("Can not get food")
		return err
	}

	factor = transaction.Amount / 100

	transaction.Nutrition.Energy = food.Nutrition.Energy * factor
	transaction.Nutrition.Protein = food.Nutrition.Protein * factor
	transaction.Nutrition.Fat = food.Nutrition.Fat * factor
	transaction.Nutrition.Carbohydrate = food.Nutrition.Carbohydrate * factor

	transaction.FoodName = food.Name
	transaction.Food = food

	log.Printf("Try to update for date %s\n", transaction.Date)

	filter := bson.D{{"date", dateString}, {"userid", userId}}
	update := bson.D{
		{"$inc", bson.D{
			{"nutrition.energy", transaction.Nutrition.Energy},
		}},
		{"$inc", bson.D{
			{"nutrition.fat", transaction.Nutrition.Fat},
		}},
		{"$inc", bson.D{
			{"nutrition.carbohydrate", transaction.Nutrition.Carbohydrate},
		}},
		{"$inc", bson.D{
			{"nutrition.protein", transaction.Nutrition.Protein},
		}},
	}
	err = MongoClient.Update(filter, update, DbName, IntakeHistory)
	if err != nil {
		log.Println("Could not update intake")
		return err
	}

	return nil
}

func getTransactions(userId string, date time.Time) ([]Transaction, error) {
	year, month, day := date.Date()
	begin := time.Date(year, month, day, 0, 0, 0, 0, date.Location())
	end := time.Date(year, month, day, 23, 59, 59, 999, date.Location())
	filter := bson.D{
		{"date", bson.D{
			{"$gte", begin},
			{"$lte", end},
		}}, {"userid", userId},
	}
	var result []Transaction
	findOptions := options.Find()
	curser, err := MongoClient.GetMultiple(filter, DbName, TRANSACTIONS, findOptions)
	if err != nil {
		return nil, err
	}
	err = curser.All(context.Background(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getHistory(userId string) ([]Intake, error) {
	var result []Intake
	findOptions := options.Find()
	findOptions.SetLimit(7)
	// Sort by `price` field descending
	findOptions.SetSort(bson.D{{"date", -1}})
	curser, err := MongoClient.GetMultiple(bson.D{{"userid", userId}}, DbName, IntakeHistory, findOptions)
	if err != nil {
		return nil, err
	}
	err = curser.All(context.Background(), &result)
	if err != nil {
		return nil, err
	}
	log.Printf("Results: %v", result)
	return result, nil

}

func search(query string) ([]Food, error) {

	filter := bson.A{
		bson.D{{"$search", bson.D{{"text", bson.D{
			{"query", query},
			{"path", "name"},
			{"fuzzy",
				bson.D{
					{"maxEdits", 2},
					{"maxExpansions", 100},
				}},
		}}}}},
		bson.D{{"$limit", 5}},
	}

	var result []Food
	curser, err := MongoClient.Search(filter, DbName, FoodList)
	if err != nil {
		return nil, err
	}
	err = curser.All(context.Background(), &result)
	log.Println(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getDailyGoals(userId string) (DailyGoals, error) {
	var n DailyGoals
	filter := bson.D{{"userId", userId}}
	result := MongoClient.Get(filter, DbName, DAILYGOALS)
	err := result.Decode(&n)
	if err != nil {
		err = insertEmptyDailyGoals(userId)
		if err != nil {
			return DailyGoals{}, err
		}
		return DailyGoals{}, nil
	}
	return n, nil
}

func insertEmptyDailyGoals(userId string) error {
	println("Insert empty")
	goals := DailyGoals{
		UserId:    userId,
		Nutrition: Nutrition{},
		Water:     0,
	}
	_, err := MongoClient.Insert(goals, DbName, DAILYGOALS)
	if err != nil {
		println("Can not add new dailyGoals")
		return err
	}
	return nil
}

func insertDailyGoals(goals DailyGoals) error {
	println("Insert empty")

	_, err := MongoClient.Insert(goals, DbName, DAILYGOALS)
	if err != nil {
		println("Can not add new dailyGoals")
		return err
	}
	return nil
}

func getDashBoardData() {

}
