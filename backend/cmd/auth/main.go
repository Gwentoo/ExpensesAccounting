package main

import (
	"backend/internal/config"
	"backend/internal/jwt"
	"backend/internal/service"
	redis2 "backend/internal/storage/cache"
	"backend/internal/storage/postgres"
	authv1 "backend/proto/gen/auth"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {

	cfg := config.LoadConfig()

	db, err := postgres.NewDB(cfg.Database.GetPostgresDSN())
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}

	redis := redis2.NewRedisDB(cfg.Redis.Addr, "", 0)

	ping, err := redis.Client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to ping cache: %v", err)
	}

	log.Println(ping)

	server := grpc.NewServer(grpc.UnaryInterceptor(jwt.CookieAuthInterceptor))

	authv1.RegisterAuthServer(server, service.NewAuthService(db, redis))

	port := fmt.Sprintf(":%d", cfg.Auth.Port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Auth server running on:", port)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
