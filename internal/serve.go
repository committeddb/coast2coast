package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func Serve(addr string, routes []*Route) error {
	for _, r := range routes {
		http.HandleFunc(fmt.Sprintf("GET /%s/{id}", r.TableName), handleRoute(r))
	}

	return http.ListenAndServe(addr, nil)
}

func handleRoute(r *Route) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		reqToken := getTokenFromHeader(req)

		err := verifyToken(r.JWTVerifyURL, reqToken)
		if err != nil {
			handleError(w, err)
			return
		}

		id := req.PathValue("id")

		err = verifyIdFromToken(reqToken, id)
		if err != nil {
			handleError(w, err)
			return
		}

		v, err := r.Datastore.Query(r.Query, id)
		if err != nil {
			handleError(w, err)
			return
		}

		io.WriteString(w, fmt.Sprintf("%v", v))
	}
}

func getTokenFromHeader(req *http.Request) string {
	reqToken := req.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	return splitToken[1]
}

func handleError(w http.ResponseWriter, err error) {
	w.Write([]byte(err.Error()))
	w.WriteHeader(http.StatusInternalServerError)
}

func verifyToken(jwtVerifyURL, tokenString string) error {
	req, err := http.NewRequest(http.MethodGet, jwtVerifyURL, nil)
	if err != nil {
		return err
	}

	req.Header["Authorization"] = []string{"Bearer " + tokenString}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %v", res.StatusCode)
	}

	return nil
}

func verifyIdFromToken(tokenString, id string) error {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if claims["id"] != id {
			return fmt.Errorf("ids don't match")
		}

		return nil
	} else {
		return fmt.Errorf("claims are not ok")
	}
}
