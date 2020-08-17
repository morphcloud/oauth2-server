package rest_v1

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/mailru/easyjson"

	"github.com/morphcloud/oauth2-server/internal/services"
	"github.com/morphcloud/oauth2-server/pkg/http_response"
)

//easyjson:json
type UserCreds struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

//easyjson:json
type LoginResponse struct {
	TokenType             string `json:"tokenType"`
	AccessToken           string `json:"accessToken"`
	AccessTokenExpiresIn  int64  `json:"accessTokenExpiresIn"`
	RefreshToken          string `json:"refreshToken"`
	RefreshTokenExpiresIn int64  `json:"refreshTokenExpiresIn"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err.Error())
	}

	var userCreds UserCreds

	if err = easyjson.Unmarshal(body, &userCreds); err != nil {
		log.Fatalln(err)
	}

	if userCreds.Login != "test" || userCreds.Password != "test" {
		log.Fatalln("Invalid Credentials.")
	}

	// TODO Get User from users collection

	JWTToken := services.NewJWTToken()
	accessTokenExpirationTime, err := strconv.ParseInt(os.Getenv("ACCESS_TOKEN_EXPIRATION_TIME"), 10, 64)
	if err != nil {
		accessTokenExpirationTime = 14400
	}

	accessTokenString, err := JWTToken.Generate("access_token", accessTokenExpirationTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	refreshTokenExpirationTime, err := strconv.ParseInt(os.Getenv("REFRESH_TOKEN_EXPIRATION_TIME"), 10, 64)
	if err != nil {
		refreshTokenExpirationTime = 7776000
	}

	refreshTokenString, err := JWTToken.Generate("refresh_token", refreshTokenExpirationTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	loginResponse := LoginResponse{
		TokenType: "Bearer",
		AccessToken: accessTokenString,
		AccessTokenExpiresIn: accessTokenExpirationTime,
		RefreshToken: refreshTokenString,
		RefreshTokenExpiresIn: refreshTokenExpirationTime,
	}

	resp := http_response.SingleJSONResponse{Data: loginResponse}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(resp, w); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
