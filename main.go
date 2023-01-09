package main

import (
	"net"
	"os"

	"github.com/HashimovH/softwareengineer-test-task-go/app"
	"github.com/HashimovH/softwareengineer-test-task-go/app/core/repository"
	"github.com/HashimovH/softwareengineer-test-task-go/app/core/service"
	driver "github.com/HashimovH/softwareengineer-test-task-go/app/driver/rpc"
	protos "github.com/HashimovH/softwareengineer-test-task-go/app/driver/rpc/protos/tickets_service"

	"github.com/hashicorp/go-hclog"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// func main() {
// 	log := hclog.Default()

// 	gs := grpc.NewServer()
// 	protos.RegisterCurrencyServer(gs)
// }

func main() {
	log := hclog.Default()

	DB := app.InitDB()

	ratingRepository := repository.NewRepository(DB)
	service_layer := service.NewService(ratingRepository)

	gs := grpc.NewServer()
	driver := driver.NewRPCAdapter(service_layer)

	protos.RegisterTicketAnalysisServiceServer(gs, driver)
	reflection.Register(gs)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Error("Unable to listen port 8080", "error", err)
		os.Exit(1)
	}
	log.Info("App started")
	gs.Serve(l)
}
