package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jinzhu/gorm"
	pb "github.com/kaansari/shippy-user-service/proto/auth"
	"github.com/kaansari/shippy-user-service/user"
	_ "github.com/lib/pq" // here its needed for GORM to work
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

func main() {

	// Creates a database connection and handles
	// closing it again before exit.
	db, err := createConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	// Automatically migrates the user struct
	// into database columns/types etc. This will
	// check for changes and migrate them each time
	// this service is restarted.
	db.AutoMigrate(&pb.User{})

	repo := &user.UserRepository{db}

	tokenService := &user.TokenService{repo}

	// Set-up our gRPC server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	pb.RegisterAuthServer(s, &user.Service{repo, tokenService})
	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to  start serve: %v", err)
	}
}

func createConnection() (*gorm.DB, error) {

	// Get database details from environment variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	DBName := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	return gorm.Open(
		"postgres",
		fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=disable",
			user, password, host, DBName,
		),
	)
}
