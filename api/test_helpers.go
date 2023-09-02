package api

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelReservation_golang/db"
	"log"
	"testing"
)

const (
	testdbname = "hotel-reservation-test"
)

type testDB struct {
	Store  *db.Store
	client *mongo.Client
}

func setup(t *testing.T) *testDB {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	return &testDB{
		Store: &db.Store{
			Hotels:   hotelStore,
			Users:    db.NewMongoUserStore(client),
			Rooms:    db.NewMongoRoomStore(client, hotelStore),
			Bookings: db.NewMongoBookingStorage(client),
		},
		client: client,
	}
}

func (tdb *testDB) teardown(t *testing.T) {
	if err := tdb.client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
