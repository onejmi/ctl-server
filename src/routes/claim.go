package routes

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"../data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var database *mongo.Database

//Claim - Request to claim light bulb
func Claim(w http.ResponseWriter, r *http.Request) {
	database = data.DatabaseClient.Database(data.DatabaseName)
	result := database.Collection("users").FindOne(context.TODO(), bson.D{{Key: "email",
		Value: strings.ToLower(r.Header.Get("email"))}})
	var user data.User
	result.Decode(&user)

	rawBody, _ := ioutil.ReadAll(r.Body)
	var body map[string]string
	json.Unmarshal(rawBody, &body)

	setOwner(&user, body["message"])

	filter := bson.D{{Key: "email", Value: strings.ToLower(r.Header.Get("email"))}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "time", Value: user.Time}}}}
	database.Collection("users").UpdateOne(context.TODO(), filter, update)

	jsonSuccess, _ := json.Marshal(data.Success{Message: "Sucessfully claimed bulb."})
	w.Write(jsonSuccess)
}

func setOwner(user *data.User, message string) {
	filter := bson.D{{Key: "name", Value: "main"}}
	update := bson.D{{Key: "$set",
		Value: bson.D{{Key: "username", Value: user.Username}, {Key: "message", Value: message}},
	}}
	upsert := true
	_, err := database.Collection("bulbs").UpdateOne(context.TODO(), filter,
		update, &options.UpdateOptions{Upsert: &upsert})
	if err != nil {
		panic(err)
	}
	user.Time = time.Now().Unix()
}
