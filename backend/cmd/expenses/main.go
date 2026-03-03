package main

import (
	"backend/internal/config"
	"backend/internal/jwt"
	"backend/internal/service"
	"backend/internal/storage/postgres"
	expensesv1 "backend/proto/gen/expenses"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {

	cfg := config.LoadConfig()

	db, err := postgres.NewDB(cfg.Database.GetPostgresDSN())
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(jwt.CookieAuthInterceptor))

	expensesv1.RegisterExpensesServer(server, service.NewExpensesService(db))

	port := fmt.Sprintf(":%d", cfg.Expenses.Port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Expenses server running on:", port)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
