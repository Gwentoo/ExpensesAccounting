package main

import (
	"backend/internal/config"
	authv1 "backend/proto/gen/auth"
	expensesv1 "backend/proto/gen/expenses"
	"context"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, X-Requested-With, X-Grpc-Web")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Expose-Headers", "Grpc-Metadata, Grpc-Status, Grpc-Message")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	ctx := context.Background()
	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			md := metadata.MD{}

			auth := req.Header.Get("Authorization")

			if auth == "" {
				auth = req.Header.Get("authorization")
			}

			if auth != "" {
				md.Set("authorization", auth)
			}

			return md
		}),
	)

	cfg := config.LoadConfig()

	authConn, err := grpc.NewClient(
		cfg.Auth.Host+":"+strconv.Itoa(cfg.Auth.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	err = authv1.RegisterAuthHandler(ctx, mux, authConn)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	expensesConn, err := grpc.NewClient(
		cfg.Expenses.Host+":"+strconv.Itoa(cfg.Expenses.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	err = expensesv1.RegisterExpensesHandler(ctx, mux, expensesConn)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	handler := allowCORS(mux)

	gwServer := &http.Server{
		Addr:              ":" + cfg.GatewayPort,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Println("Serving gRPC-Gateway on http://localhost:8080")
	log.Fatalln(gwServer.ListenAndServe())
}
