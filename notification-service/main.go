package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	. "nilswilhelm.net/foodtracker/lib/constants"
	"nilswilhelm.net/foodtracker/lib/database"
	. "nilswilhelm.net/foodtracker/lib/structs"
	"time"
)

var mongoClient database.MongoAdapter

const MongoUri = "mongodb+srv://nils:v0W0rHH4gJHay61N@cluster0.ddwq5.mongodb.net/foodtracker?retryWrites=true&w=majority"

func main() {
	client, err := database.NewMongoAdapter(MongoUri)
	if err != nil {
		log.Fatal(err)
	}
	mongoClient = client

	goals, err := getDailyGoals()
	if err != nil {
		panic(err)
	}

	intake, err := getIntake()
	if err != nil {
		panic(err)
	}

	userData := NewMapper(goals, intake).DoMap()

	for _, data := range userData {

		handleNotification(data)
	}

}

func handleNotification(data DailyGoalsAndIntake) {
	if !goalsReached(data, .5) {
		log.Println("Not reached yet")
	}
}

func getDailyGoals() ([]DailyGoals, error) {
	filter := bson.D{{}}

	var userData []DailyGoals
	findOptions := options.Find()
	cursor, err := mongoClient.GetMultiple(filter, DbName, DAILYGOALS, findOptions)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &userData)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func getIntake() ([]Intake, error) {

	today := time.Now().Format(DateFormat)
	filter := bson.D{{"date", today}}

	var userData []Intake
	findOptions := options.Find()
	cursor, err := mongoClient.GetMultiple(filter, DbName, IntakeHistory, findOptions)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &userData)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func goalsReached(data DailyGoalsAndIntake, factor float64) bool {
	intake := data.intake.Nutrition
	goals := data.goals.Nutrition

	return intake.Energy >= goals.Energy*factor &&
		intake.Carbohydrate >= goals.Carbohydrate*factor &&
		intake.Protein >= goals.Protein*factor &&
		intake.Fat >= goals.Fat*factor
}
