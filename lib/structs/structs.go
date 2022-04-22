package structs

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"nilswilhelm.net/foodtracker/lib/database"
	"time"
)

var MongoClient database.MongoAdapter

type Nutrition struct {
	Energy       float64 `json:"energy,omitempty"`
	Fat          float64 `json:"fat,omitempty"`
	Carbohydrate float64 `json:"carbohydrate,omitempty"`
	Protein      float64 `json:"protein,omitempty"`
}

type Intake struct {
	Date      string    `json:"date,omitempty"`
	Nutrition Nutrition `json:"nutrition,omitempty"`
	UserId    string    `json:"userId,omitempty"`
}

// default type to store food in db
type Food struct {
	ID *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	// EAN is only set if product has one
	EAN string `json:"ean,omitempty"`
	// Name of food/meal
	Name string `json:"name,omitempty"`
	// UserId is only set for combined meal
	UserId string `json:"userId,omitempty"`
	// Macro nutrition values
	Nutrition Nutrition `json:"nutrition,omitempty"`
	// IsMeal flag if Food is a combined meal
	IsMeal bool `json:"isMeal,omitempty"`
	// Ingredients of meal
	Ingredients []Ingredient `json:"ingredients,omitempty"`
	// MetaFields for example brand
	MetaFields map[string]string `json:"meta_fields,omitempty"`
	// Brand
	Brand string `json:"brand,omitempty"`
}

type Transaction struct {
	ID        *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FoodID    string              `json:"foodId,omitempty"`
	Food      Food                `json:"food,omitempty"`
	FoodName  string              `json:"foodName,omitempty"`
	Amount    float64             `json:"amount,omitempty"`
	Nutrition Nutrition           `json:"nutrition,omitempty"`
	IsMeal    bool                `json:"isMeal,omitempty"`
	Date      time.Time           `json:"date,omitempty"`
	UserId    string              `json:"userId,omitempty"`
}

type MealRequest struct {
	// Ingredients stores id and amount of ingredients
	Ingredients []Ingredient `json:"ingredients,omitempty"`
	//Name of meal
	Name string `json:"name,omitempty"`
}

type Ingredient struct {
	Food   Food    `json:"food,omitempty"`
	Amount float64 `json:"amount,omitempty"`
}

type DailyGoals struct {
	Nutrition Nutrition `json:"nutrition"`
	UserId    string    `json:"userId" bson:"userId"`
	Water     float64   `json:"water"`
}

type DashboardData struct {
	Intake       Intake        `json:"intake"`
	DailyGoals   DailyGoals    `json:"dailyGoals"`
	Transactions []Transaction `json:"transactions"`
}
