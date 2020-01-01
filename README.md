#This is a REST API demo using Go language
##The external packages used for this project are
1. go get github.com/gorilla/mux
2. go get github.com/dgrijalva/jwt-go
3. go get github.com/lib/pq
4. go get golang.org/x/crypto/bcrypt
5. go get -u github.com/davecgh/go-spew/spew for json formatting detailed info
6. go get github.com/subosito/gotenv  - store configurations into .env file on root directory

##Setup the connectio with Elephant Postgres as follows:
https://www.elephantsql.com/docs/go.html

#sample login 
{
  "email": "test222@example.com",
  "password": "abcd123"
}
