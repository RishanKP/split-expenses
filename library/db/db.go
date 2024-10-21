package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"split-expenses/library/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client        *mongo.Client
	connectionURI string
)

func init() {
	dbUser := config.DB_USER
	dbPass := config.DB_PASS
	dbCluster := config.DB_CLUSTER
	dbName := config.DB_NAME

	connectionURI = fmt.Sprintf("mongodb+srv://%s:%s@electromate.%s.mongodb.net/?retryWrites=true&w=majority&appName=%s", dbUser, dbPass, dbCluster, dbName)
}

func Connect() {
	var err error
	clientOptions := options.Client().ApplyURI(connectionURI)

	Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB!")
}

func Disconnect() {
	if err := Client.Disconnect(context.TODO()); err != nil {
		log.Fatalf("Failed to disconnect from MongoDB: %v", err)
	}
	fmt.Println("Disconnected from MongoDB")
}
