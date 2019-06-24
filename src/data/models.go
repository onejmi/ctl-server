package data

import "time"

//User - Data to register a new user to the server
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Time     int64  `json:"time"`
}

//Session - Used to authenticate actions performed by user after logging in
type Session struct {
	Username  string `json:"username" bson:"username"`
	SessionID string `json:"session_id" bson:"session_id"`
	Created   int64  `json:"creation_time" bson:"creation_time"`
}

//Error - Format for serving errors to the caller
type Error struct {
	Message string `json:"error"`
}

//Success - Format for serving success notifications to the caller
type Success struct {
	Message string `json:"success"`
}

//Methods...

const secondsInHour = 60 * 60

//Expired - Calculates whether the session is still valid
func (session Session) Expired() bool {
	creationTimeElapsed := time.Now().Unix() - session.Created
	return creationTimeElapsed >= secondsInHour
}
