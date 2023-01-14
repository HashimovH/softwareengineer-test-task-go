package main

import (
	"net"
	"os"

	"github.com/HashimovH/softwareengineer-test-task-go/app/config"
	"github.com/HashimovH/softwareengineer-test-task-go/app/core/repository"
	"github.com/HashimovH/softwareengineer-test-task-go/app/core/service"
	driver "github.com/HashimovH/softwareengineer-test-task-go/app/driver/rpc"
	protos "github.com/HashimovH/softwareengineer-test-task-go/app/driver/rpc/protos/tickets_service"

	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}

	DB := config.InitDB()

	// Register services
	scoreRepository := repository.NewRepository(DB)
	scoreService := service.NewService(scoreRepository)

	qualityRepository := repository.NewQualityRepository(DB)
	qualityService := service.NewQualityService(qualityRepository)

	driver := driver.NewRPCAdapter(scoreService, qualityService)

	gs := grpc.NewServer()

	protos.RegisterTicketAnalysisServiceServer(gs, driver)

	if isDevelopment() {
		log.Info("Starting in development mode")
		reflection.Register(gs)
	} else {
		log.Info("Starting in production mode")
	}

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Error("Unable to listen port 8080", "error", err)
		os.Exit(1)
	}
	log.Info("App started")
	if err := gs.Serve(l); err != nil {
		log.Error("error grpc serve: %v", err)
		os.Exit(1)
	}
}

func isDevelopment() bool {
	env := os.Getenv("APP_ENV")
	return env != "production"
}
