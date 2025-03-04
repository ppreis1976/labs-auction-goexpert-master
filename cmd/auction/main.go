package main

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/internal/infra/api/web/controller/auction_controller"
	"fullcycle-auction_go/internal/infra/api/web/controller/bid_controller"
	"fullcycle-auction_go/internal/infra/api/web/controller/user_controller"
	"fullcycle-auction_go/internal/infra/database/auction"
	"fullcycle-auction_go/internal/infra/database/bid"
	"fullcycle-auction_go/internal/infra/database/user"
	"fullcycle-auction_go/internal/usecase/auction_usecase"
	"fullcycle-auction_go/internal/usecase/bid_usecase"
	"fullcycle-auction_go/internal/usecase/user_usecase"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}

	databaseConnection, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	if err := initializeDatabase(ctx, databaseConnection); err != nil {
		log.Fatal(err.Error())
		return
	}

	router := gin.Default()

	userController, bidController, auctionsController := initDependencies(databaseConnection)

	router.GET("/auction", auctionsController.FindAuctions)
	router.GET("/auction/:auctionId", auctionsController.FindAuctionById)
	router.POST("/auction", auctionsController.CreateAuction)
	router.GET("/auction/winner/:auctionId", auctionsController.FindWinningBidByAuctionId)
	router.POST("/bid", bidController.CreateBid)
	router.GET("/bid/:auctionId", bidController.FindBidByAuctionId)
	router.GET("/user/:userId", userController.FindUserById)

	router.Run(":8080")
}

func initializeDatabase(ctx context.Context, client *mongo.Database) error {
	collections := []string{"users", "auctions", "bids"}

	for _, collectionName := range collections {
		err := client.CreateCollection(ctx, collectionName)
		if err != nil {
			return fmt.Errorf("error creating collection %s: %v", collectionName, err)
		}
		fmt.Printf("Collection %s created successfully\n", collectionName)
	}

	userCollection := client.Collection("users")
	testUser := bson.M{
		"_id":  "83b9f719-aa98-4135-88b1-3c5d6b7e5f34",
		"name": "Pericles Paulo Dos Reis",
	}

	// Check if the test user already exists
	var existingUser bson.M
	err := userCollection.FindOne(ctx, bson.M{"_id": testUser["_id"]}).Decode(&existingUser)
	if err != nil && err != mongo.ErrNoDocuments {
		return fmt.Errorf("error checking for existing test user: %v", err)
	}

	if existingUser == nil {
		_, err := userCollection.InsertOne(ctx, testUser)
		if err != nil {
			return fmt.Errorf("error inserting test user: %v", err)
		}
		fmt.Println("Test user inserted successfully")
	} else {
		fmt.Println("Test user already exists, skipping insertion")
	}

	return nil
}

func initDependencies(database *mongo.Database) (
	userController *user_controller.UserController,
	bidController *bid_controller.BidController,
	auctionController *auction_controller.AuctionController) {

	auctionRepository := auction.NewAuctionRepository(database)
	bidRepository := bid.NewBidRepository(database, auctionRepository)
	userRepository := user.NewUserRepository(database)

	userController = user_controller.NewUserController(
		user_usecase.NewUserUseCase(userRepository))
	auctionController = auction_controller.NewAuctionController(
		auction_usecase.NewAuctionUseCase(auctionRepository, bidRepository))
	bidController = bid_controller.NewBidController(bid_usecase.NewBidUseCase(bidRepository))

	return
}
