package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	// "fmt"
	"net/http"
	"time"
)

//-----

func signup(w http.ResponseWriter, req *http.Request) {
	var u User
	var un UserAuth

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err := decoder.Decode(&u)
	if err != nil {
		panic(err)
	}
	
	//check if username exists
	db.Where(&UserAuth{Username: u.Username}).First(&un)
	if un.Username != "" {
		w.Header().Set("Content-Type", "application/json")
		j, _ := json.Marshal("User already exists")
		w.Write(j)
		return
	}

	//check if email exists
	db.Where(&UserAuth{Email: u.Email}).First(&un)
	if un.Email != "" {
		w.Header().Set("Content-Type", "application/json")
		j, _ := json.Marshal("Email already exists")
		w.Write(j)
		return
	}
	
	//generate encrypted password
	bs, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//generate user in database
	user := UserAuth{Username: u.Username, Name: u.Name, Email: u.Email, Password: string(bs)}

	db.NewRecord(user)
	db.Create(&user)

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal("User created")
	w.Write(j)
	return
}

func login(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		// still need to implement JSON web token to see if user is already logged in
		var u User
		var un UserAuth

		decoder := json.NewDecoder(req.Body)
		defer req.Body.Close()
		err := decoder.Decode(&u)
		if err != nil {
			panic(err)
		}

		//check if username is valid
		db.Where(&UserAuth{Username: u.Username}).First(&un)
		log.Println(un.Username)

		if un.Username == "" {
			w.Header().Set("Content-Type", "application/json")
			j, _ := json.Marshal("Username or password does not match")
			w.Write(j)
			return
		}

		//compare passwords
		err = bcrypt.CompareHashAndPassword([]byte(un.Password), []byte(u.Password))

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			j, _ := json.Marshal("Username or password does not match")
			w.Write(j)
			return
		}

		//issue token upon successful login
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"Username": u.Username,
			"Time":     time.Now().Add(time.Hour * 72).Unix(),
		})

		tokenString, err := token.SignedString(mySigningKey)
		w.Header().Set("Content-Type", "application/json")
		j, _ := json.Marshal(tokenString)
		w.Write(j)
	}
}

func protected(w http.ResponseWriter, req *http.Request) {
	log.Println("Protected resource served")
}