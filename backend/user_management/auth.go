package user_management

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var SigningKey []byte

type AuthResponse struct {
	Token     string `json:"token"`
	UserID    string `json:"userId"`
	ExpiresIn int    `json:"expiresIn"`
}

func MakeLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Login handler called")
		var loginData User

		err := json.NewDecoder(r.Body).Decode(&loginData)
		if err != nil {
			fmt.Println("Error while decoding json")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if authenticateUser(loginData.Email, loginData.Password) {
			token, err := GenerateJWT(loginData.Email)
			if err != nil {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			responseData := AuthResponse{Token: token, UserID: loginData.Email, ExpiresIn: 60 * 24 * 7}
			err = json.NewEncoder(w).Encode(responseData)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			_, err = w.Write([]byte("Login Failed"))
			if err != nil {
				log.Print("could not write response body")
			}
		}
	}
}

func MakeSignUpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Got request")
		var newUser User

		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			fmt.Println("Error while decoding json")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if user exists
		existingUser, err := GetUser(newUser.Email)
		// If not create new user
		if err != nil {
			err = AddUser(newUser)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusConflict)
			_, err := w.Write([]byte(fmt.Sprintf("User %s already exists", existingUser)))
			if err != nil {
				log.Print("could not write response body")
			}
			log.Print("User exists")
		}

	}
}

func hashPassword(password []byte) []byte {
	pwd, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Could not hash password")
	}
	return pwd
}

func verifyPassword(plainPassword []byte, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), plainPassword)
	if err != nil {
		return false
	}
	return true
}

func authenticateUser(identifier string, password string) bool {
	log.Println("Authenticate User")
	userEntry, err := GetUser(identifier)
	if err != nil {
		return false
	}
	return verifyPassword([]byte(password), userEntry.Password)
}

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")

		// Validate token
		if tokenString != "" && tokenString != "null" {
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("there was an error")
				}
				return SigningKey, nil
			})

			if err != nil {
				w.WriteHeader(401)
				return
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("Not authorized"))
			if err != nil {
				log.Print("could not write response body")
			}
		}
	})
}

// Generates JWT token for given user
func GenerateJWT(user string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = user
	claims["exp"] = time.Now().Add(time.Minute * 60 * 24 * 7).Unix()

	tokenString, err := token.SignedString(SigningKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
