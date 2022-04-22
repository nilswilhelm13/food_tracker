module FoodTracker

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/elastic/go-elasticsearch/v7 v7.9.0
	github.com/gorilla/mux v1.7.4
	github.com/rs/cors v1.7.0
	go.mongodb.org/mongo-driver v1.4.4
	golang.org/x/crypto v0.0.0-20190530122614-20be4c3c3ed5
	nilswilhelm.net/foodtracker/lib v0.0.0
)

replace nilswilhelm.net/foodtracker/lib => ../lib
