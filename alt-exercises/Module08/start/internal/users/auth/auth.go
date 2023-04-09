package auth

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/someuser/gameserver/internal/users"
)

type JwtAuthenticator struct{}

const (
	privKeyPath = "/../keys/app.rsa"     // openssl genrsa -out app.rsa keysize
	pubKeyPath  = "/../keys/app.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
	TokenName   = "x-access-token"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// read the key files
func initKeys() {
	dir, _ := os.Getwd()

	signBytes, err := ioutil.ReadFile(dir + privKeyPath)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(dir + pubKeyPath)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)
}

var authenticator *JwtAuthenticator

func GetAuthenticator() *JwtAuthenticator {
	if authenticator == nil {
		authenticator = &JwtAuthenticator{}
		initKeys()
	}
	return authenticator
}

//check in the header cookie and query params for the token with the key "x-access-token"
func (jwtAuth JwtAuthenticator) IsTokenExists(r *http.Request) (bool, string) {

}

//validate that the token for the user is a vlid token
func (jwtAuth JwtAuthenticator) IsUserTokenValid(token string) bool {
}

//returns the user from a given token
func (jwtAuth JwtAuthenticator) UserFromToken(tokenString string) (*users.User, error) {

}

// JwtVerify Middleware function
func (jwtAuth JwtAuthenticator) JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

//returns the token for the given user
func (jwtAuth JwtAuthenticator) GetTokenForUser(user *users.User) (string, error) {

}
