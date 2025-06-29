package tests

import (
	"context"
	"log"
	"os"
	"testing"

	testapp "gitlab.com/sandstone2/fiberpoc/app/int_testing/test_app"
)

var ctx context.Context

func TestMain(m *testing.M) {
	defer testapp.CloseDbAndLogger()

	ctx = context.Background()

	// Change to the same working directory as the main app. This is so all the relative paths in the app match.
	err := os.Chdir("../../")
	if err != nil {
		log.Printf("Error: 5PGWRW - Changing working directory. Error: %v", err)
		os.Exit(1)
	}

	app, err := testapp.GetApp()
	if err != nil {
		log.Printf("Error: WCQK0D - Getting app for testing. Error: %v", err)
		os.Exit(1)
	}

	ctx = context.WithValue(ctx, "App", app)

	// Run all the tests
	exitVal := m.Run()

	// Exit with the test result code
	os.Exit(exitVal)
}

func GetContext() *context.Context {
	return &ctx
}
