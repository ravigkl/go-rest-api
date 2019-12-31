package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"fmt"
)

func main()  {

	router := mux.NewRouter()
	router.HandleFunc("/signup", signup).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("protected", TokenVerifyMiddleWare(protectedEndpoint)).Methods("GET")

	log.Println("Listen on port 8000 ...")
	log.Fatal(http.ListenAndServe(":8000", router))

	
}

func signup(w http.ResponseWriter, r *http.Request){
	fmt.Println("signup invoked")
	w.Write([]byte("successfully called signup"))
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
