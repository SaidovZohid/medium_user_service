package main

import (
	"fmt"
	"log"
	"net"

	"github.com/SaidovZohid/medium_user_service/config"
	pb "github.com/SaidovZohid/medium_user_service/genproto/user_service"
	grpcPkg "github.com/SaidovZohid/medium_user_service/pkg/grpc_client"
	"github.com/SaidovZohid/medium_user_service/pkg/logger"
	"github.com/SaidovZohid/medium_user_service/service"
	"github.com/SaidovZohid/medium_user_service/storage"

	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load(".")

	psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	psqlConn, err := sqlx.Connect("postgres", psqlUrl)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	fmt.Println("Configuration: ", cfg)
	fmt.Println("Connected Succesfully!")

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})

	strg := storage.NewStoragePg(psqlConn)
	inMemory := storage.NewInMemoryStorage(rdb)
	
	grpcConn, err := grpcPkg.New(cfg)
	if err != nil {
		log.Fatalf("failed to get grpc connections: %v\n", err)
	}

	logger := logger.New()

	userService := service.NewUserService(strg, inMemory, logger)
	authService := service.NewAuthService(strg, inMemory, grpcConn, &cfg, logger)

	listen, err := net.Listen("tcp", cfg.GrpcPort) 

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, userService)
	pb.RegisterAuthServiceServer(s, authService)
	reflection.Register(s)

	log.Println("gRPC server started port in: ", cfg.GrpcPort)
	if s.Serve(listen);err != nil {
		log.Fatalf("Error while listening: %v", err)
	}
}
