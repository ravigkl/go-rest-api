package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"fmt"
	"strings"
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	"github.com/davecgh/go-spew/spew"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"

)

//business models for application
type User struct {
	ID		int	`json: "id"`
	Email	string `json: "email"`
	Password string	`json: "password"`
}
type JWT struct {
	Token string `json: "token"`
}

type Error struct {
	Message string `json: "message"`
}

var db *sql.DB
func main()  {
	//pgUrl, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	pgUrl, err := pq.ParseURL("postgres://ubdonwkq:4-38EnF3sbuKUbVGlptbCikADaOjlwcp@rajje.db.elephantsql.com:5432/ubdonwkq")

	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(pgUrl)
	db, err = sql.Open("postgres", pgUrl)

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(db)
	//check that connection is established with Db else error
	err = db.Ping()
	if err !=nil {
		log.Fatal("DB connection failed")
	}

	router := mux.NewRouter()
	router.HandleFunc("/signup", signup).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/protected", TokenVerifyMiddleWare(protectedEndpoint)).Methods("GET")

	log.Println("Listen on port 8000 ...")
	log.Fatal(http.ListenAndServe(":8000", router))

	
}

func respondWithError(w http.ResponseWriter, status int, error Error){
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

func responseJSON(w http.ResponseWriter, data interface{}){
	json.NewEncoder(w).Encode(data)
}

func GenerateToken(user User) (string, error){
	var err error
	secret := "mysecretkey"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email" : user.Email,
		"iss" : "course",
	})

	tokenString, err := token.SignedString([]byte(secret))
	
	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil

}

func signup(w http.ResponseWriter, r *http.Request){
	var user User
	var error Error

	json.NewDecoder(r.Body).Decode(&user)

	if user.Email == ""{
		error.Message = "User name is missing"
		respondWithError(w, http.StatusBadRequest, error)
		return
	}
	
	if user.Password == ""{
		error.Message = "User Password is missing"
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte (user.Password), 10)

	if err != nil {
		log.Fatal(err)
	}
	user.Password= string(hash)
	stmt := "insert into users (email, password) values ($1, $2) RETURNING id;"
	err = db.QueryRow(stmt, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		error.Message = "Server Error"
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	responseJSON(w, user)
	spew.Dump(user)	
}

func login(w http.ResponseWriter, r *http.Request){
	var user User
	var jwt JWT
	var error Error
	json.NewDecoder(r.Body).Decode(&user)
	password := user.Password

	if user.Email == "" {
		error.Message = "Email is missing"
		respondWithError(w, http.StatusBadRequest, error)
		return
	}
	if user.Password == "" {
		error.Message = "Password is missing"
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	row := db.QueryRow("select * from users where email=$1", user.Email)
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == 	sql.ErrNoRows {
			error.Message = "The user does not exist"
			respondWithError(w, http.StatusBadRequest, error)
			return
		}else{
			log.Fatal(err)
		}

	}
	hashedPassword := user.Password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		error.Message="Invalid password"
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	token, err := GenerateToken(user)

	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	jwt.Token = token
	responseJSON(w, jwt)
}

func protectedEndpoint (w http.ResponseWriter, r *http.Request){
	fmt.Println("Protected Endpoint invoked")
	w.Write([]byte("successfully called Protected Endpoint"))
}

func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		var errorObject Error
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) == 2 {
			authToken := bearerToken[1]

			token, error := jwt.Parse(authToken, func(token *jwt.Token)(interface{}, error){
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return []byte ("mysecretkey"), nil
			})
			if error != nil {
				errorObject.Message = error.Error()
				respondWithError(w, http.StatusBadRequest, errorObject)
				return
			}
			if token.Valid {
				next.ServeHTTP(w, r)
			}else{
				errorObject.Message = error.Error()
				respondWithError(w, http.StatusBadRequest, errorObject)
				return
			}
		}else{
			errorObject.Message="Invalid toekn"
			respondWithError(w, http.StatusBadRequest, errorObject)
			return
		}
	})
}

