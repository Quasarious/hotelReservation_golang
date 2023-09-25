package api

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelReservation_golang/db"
	"hotelReservation_golang/my_conf"
	"log"
	"testing"
)

type testDB struct {
	Store  *db.Store
	client *mongo.Client
}

func setup(t *testing.T) *testDB {
	my_conf.LoadEnv()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(my_conf.DBURL))
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
	if err := tdb.client.Database(my_conf.DBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
