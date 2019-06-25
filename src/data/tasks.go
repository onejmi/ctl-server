package data

import (
	"context"

	"github.com/jasonlvhit/gocron"
)

func SetupCronJobs() {
	gocron.Every(1).Day().Do(clearSessions)
}

func clearSessions() {
	DatabaseClient.
		Database(DatabaseName).
		Collection("sessions").
		Drop(context.TODO())
}
