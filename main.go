package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"fmt"
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	"github.com/davecgh/go-spew/spew"
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
	router.HandleFunc("protected", TokenVerifyMiddleWare(protectedEndpoint)).Methods("GET")

	log.Println("Listen on port 8000 ...")
	log.Fatal(http.ListenAndServe(":8000", router))

	
}

func respondWithError(w http.ResponseWriter, status int, error Error){
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
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
	
	spew.Dump(user)	
}
func login(w http.ResponseWriter, r *http.Request){
	fmt.Println("login invoked")
	w.Write([]byte("successfully called login"))
}

func protectedEndpoint (w http.ResponseWriter, r *http.Request){
	fmt.Println("Protected Endpoint invoked")
	w.Write([]byte("successfully called Protected Endpoint"))
}

func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	fmt.Println("TokenVerifyrMiddleWare invoked")
	return nil
}
