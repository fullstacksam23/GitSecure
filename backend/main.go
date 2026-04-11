package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/fullstacksam23/GitSecure/internal/api"
	"github.com/fullstacksam23/GitSecure/internal/db"
	"github.com/fullstacksam23/GitSecure/internal/worker"
	"github.com/joho/godotenv"
)

//TODO: REFACTOR THE MODELS FOLDER -- GET RID OF MODELS FOLDER AND PUT EACH MODEL IN THEIR PACKAGE -- SHARED MODELS IN CORE

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel() // cancels the current context at the end of the current function

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	//Init DB
	url := os.Getenv("SUPABASE_URL")
	service_role_key := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	err = db.InitSupabase(url, service_role_key)
	if err != nil {
		log.Fatal("Failed to connect to Supabase:", err)
	}

	mode := os.Getenv("MODE")

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("GITHUB_TOKEN not set")
	}

	if mode == "worker" {

		fmt.Println("Starting worker...")
		worker.StartWorker(ctx, githubToken)
		return
	}

	fmt.Println("Starting API server")
	app := api.New()
	err = app.Start(ctx)
	if err != nil {
		fmt.Println("server failed to start")
	}
}
