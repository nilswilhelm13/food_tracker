package user_management

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"nilswilhelm.net/foodtracker/lib/database"
)

var MongoClient database.MongoAdapter

const DbName = "users"

type User struct {
	Username string
	Email    string
	Password string
}

type UserResponse struct {
	Username string `json:"username"`
}

func AddUser(u User) error {
	hashedPwd := hashPassword([]byte(u.Password))
	u.Password = string(hashedPwd)

	_, err := MongoClient.Insert(u, DbName, "users")
	if err != nil {
		return err
	}
	return nil
}

func GetUser(email string) (User, error) {
	var u User
	filter := bson.D{{"email", email}}
	result := MongoClient.Get(filter, DbName, "users")
	err := result.Decode(&u)
	if err != nil {
		return User{}, err
	}
	fmt.Printf("Got user: %s\n", u.Email)
	return u, nil
}
