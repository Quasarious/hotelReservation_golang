package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelReservation_golang/api"
	"hotelReservation_golang/db"
	"hotelReservation_golang/db/fixtures"
	"hotelReservation_golang/types"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		Users:    db.NewMongoUserStore(client),
		Hotels:   hotelStore,
		Rooms:    db.NewMongoRoomStore(client, hotelStore),
		Bookings: db.NewMongoBookingStorage(client),
	}
	rooms := []*types.Room{}

	user := fixtures.AddUser(store, "Iak", "Vapvapvv", false)
	fmt.Println("IakJWT -> ", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("adminJWT -> ", api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(store, "hotel1", "Moscow", 3.5, nil)
	rooms = append(rooms, fixtures.AddRoom(store, "large", true, 199.99, hotel.ID))
	rooms = append(rooms, fixtures.AddRoom(store, "small", false, 29.99, hotel.ID))
	booking := fixtures.AddBooking(store, user.ID, rooms[0].ID, 2, time.Now(), time.Now().AddDate(0, 0, 3))
	fmt.Println("booking id -> ", booking.ID)
}
