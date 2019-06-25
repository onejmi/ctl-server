package routes

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"

	"../data"
	"go.mongodb.org/mongo-driver/bson"
)

//LoginRequest - structure which stores data coming from a login request to the server
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Login - Request to login a user onto the application
func Login(w http.ResponseWriter, r *http.Request) {
	raw, _ := ioutil.ReadAll(r.Body)
	var loginReq LoginRequest
	json.Unmarshal(raw, &loginReq)
	session, authenticated := authenticate(loginReq)
	if authenticated {
		sessionJSON, err := json.Marshal(session)
		if err != nil {
			panic(err)
		}
		w.Write(sessionJSON)
	} else {
		jsonError, _ := json.Marshal(data.Error{Code: 17, Message: "Invalid login details"})
		w.Write(jsonError)
	}
}

func authenticate(loginReq LoginRequest) (session data.Session, authenticated bool) {
	db := data.DatabaseClient.Database(data.DatabaseName)
	users := db.Collection("users")
	//check if a document exists with username + password combo
	count, err := users.CountDocuments(context.TODO(), bson.D{{Key: "email", Value: strings.ToLower(loginReq.Email)},
		{Key: "password", Value: loginReq.Password}})
	if err != nil {
		panic(err)
	}
	if count > 0 {
		authenticated = true
		rawSessionID, _ := uuid.New()
		sessionIDBytes := [16]byte(rawSessionID)
		sessionID := hex.EncodeToString(sessionIDBytes[:])
		sessions := db.Collection("sessions")
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "session_id", Value: sessionID},
				{Key: "creation_time", Value: time.Now().Unix()},
			}},
		}
		upsert := true
		_, err := sessions.UpdateOne(context.TODO(),
			bson.D{{Key: "username", Value: strings.ToLower(loginReq.Email)}},
			update, &options.UpdateOptions{Upsert: &upsert})
		if err != nil {
			panic(err)
		}
		session = data.Session{Username: loginReq.Email, SessionID: sessionID, Created: time.Now().Unix()}
	} else {
		authenticated = false
	}
	return session, authenticated
}
