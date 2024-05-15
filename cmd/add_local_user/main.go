package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

const (
	seedPath = "cmd/add_local_user/data/seed.csv"
)

type Operator struct {
	operatorID string
	email      string
	password   string
}

func main() {
	addLocal()
}

func addLocal() {
	ctx := context.Background()

	// set environment variables to use local emulator
	os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", "localhost:9099")
	csvPath := seedPath

	conf := &firebase.Config{ProjectID: "local"}
	app, err := firebase.NewApp(ctx, conf, option.WithoutAuthentication())
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	addOperatorFromCSV(ctx, app, csvPath)
}

func addOperatorFromCSV(ctx context.Context, app *firebase.App, csvPath string) {
	// make client
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Error getting Auth client: %v", err)
	}

	// read csv file
	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("Error opening CSV file: %v", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Error reading CSV records: %v", err)
	}

	// create user by each record
	for _, record := range records {
		email := record[0]
		password := record[1]
		operatorID := record[2]
		operator := Operator{operatorID, email, password}
		addOperator(ctx, authClient, operator)
	}
}

func addOperator(ctx context.Context, authClient *auth.Client, operator Operator) {
	uid := uuid.New().String()
	params := (&auth.UserToCreate{}).
		UID(uid).
		Email(operator.email).
		Password(operator.password)

	// create user
	userRecord, err := authClient.CreateUser(ctx, params)
	if err != nil {
		log.Fatalf("Error creating user for email %s: %v", operator.email, err)
	}

	// set custom claims
	customClaims := map[string]interface{}{
		"operator_id": operator.operatorID,
	}
	err = authClient.SetCustomUserClaims(ctx, userRecord.UID, customClaims)
	if err != nil {
		log.Fatalf("Error setting custom claims for user %s: %v", operator.email, err)
	}

	fmt.Printf("Successfully created user with email: %s and set custom claims. Password: %s\n", operator.email, operator.password)
}
