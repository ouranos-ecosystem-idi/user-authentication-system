package firebase

import (
	"authenticator-backend/extension/logger"
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// NewClient
// Summary: This is the function which creates the firebase client.
// input: projectID(string): project ID
// input: emulatorHost(string): emulator host
// output: (*auth.Client) auth client
// output: (error) error object
func NewClient(projectID string, emulatorHost string) (*auth.Client, error) {
	ctx := context.Background()

	var app *firebase.App
	var err error
	if emulatorHost != "" {
		conf := &firebase.Config{ProjectID: "local"}
		app, err = firebase.NewApp(ctx, conf, option.WithoutAuthentication())
		if err != nil {
			logger.Set(nil).Errorf(err.Error())

			return nil, err
		}
	} else {
		config := &firebase.Config{ProjectID: projectID}
		app, err = firebase.NewApp(ctx, config)
		if err != nil {
			logger.Set(nil).Errorf(err.Error())

			return nil, err
		}
	}

	client, err := app.Auth(ctx)
	if err != nil {
		logger.Set(nil).Errorf(err.Error())

		return nil, err
	}

	return client, nil
}
