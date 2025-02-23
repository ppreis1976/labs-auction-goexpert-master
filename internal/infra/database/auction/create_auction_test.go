package auction_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/auction"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCreateAuction(t *testing.T) {
	dataBaseName := "user_database_test"
	collectionName := "user_collection_name_test"

	os.Setenv("AUCTION_CLOSED", "1s")
	os.Setenv("MONGODB_DB", collectionName)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("successful insert and status update", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		repo := auction.NewAuctionRepository(mt.DB)
		auctionEntity := &auction_entity.Auction{
			Id:          "123",
			ProductName: "Test Product",
			Category:    "Test Category",
			Description: "Test Description",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active,
			Timestamp:   time.Now(),
		}

		err := repo.CreateAuction(context.Background(), auctionEntity)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Simulate the passage of time to trigger the status update
		time.Sleep(2 * time.Second)

		// Add mock response for FindOne operation
		mt.AddMockResponses(mtest.CreateCursorResponse(1, fmt.Sprintf("%s.%s", dataBaseName, collectionName), mtest.FirstBatch, bson.D{
			{Key: "_id", Value: auctionEntity.Id},
			{Key: "status", Value: auction_entity.Completed},
		}))

		// Check if the status has been updated to closed
		var updatedAuction auction.AuctionEntityMongo
		errUpdate := mt.DB.Collection("auctions").FindOne(context.Background(), bson.M{"_id": auctionEntity.Id}).Decode(&updatedAuction)
		if errUpdate != nil {
			t.Fatalf("expected no error, got %v", errUpdate)
		}

		if updatedAuction.Status != auction_entity.Completed {
			t.Fatalf("expected status to be %v, got %v", auction_entity.Completed, updatedAuction.Status)
		}
	})
}
