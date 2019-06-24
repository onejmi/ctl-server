package util

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/mongo/options"

	"../../data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//IsAuthenticated - Ensures user has logged in and has a valid SessionID
func IsAuthenticated(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Session") == "" {
			invalidSession(w)
		} else {
			sessionID := strings.Trim(r.Header.Get("Session"), "\"")
			allowPartialResults := false
			var database *mongo.Database = data.DatabaseClient.Database(data.DatabaseName)
			var session data.Session
			err := database.Collection("sessions").FindOne(context.TODO(),
				bson.D{{Key: "session_id", Value: sessionID}},
				&options.FindOneOptions{AllowPartialResults: &allowPartialResults}).Decode(&session)
			if err != nil {
				if strings.Contains(err.Error(), "no documents") {
					invalidSession(w)
				} else {
					panic(err)
				}
			} else {
				if session.Expired() {
					deleteSession(database, session)
					invalidSession(w)
				} else {
					r.Header.Set("username", session.Username)
					handler.ServeHTTP(w, r)
				}
			}
		}
	})
}

func extractRequestFields(r *http.Request, result *map[string]string) {
	rawJSONAuth, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(rawJSONAuth, &result)
}

func deleteSession(database *mongo.Database, session data.Session) {
	filter := bson.D{{Key: "session_id", Value: session.SessionID}}
	_, err := database.Collection("sessions").DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
}

func invalidSession(w http.ResponseWriter) {
	jsonError, _ := json.Marshal(data.Error{Message: "Invalid Session."})
	w.Write(jsonError)
}
