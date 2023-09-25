package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelReservation_golang/api"
	"hotelReservation_golang/db"
	"hotelReservation_golang/my_conf"
	"hotelReservation_golang/types"
	"log"
	"math/rand"
	"time"
)

func main() {
	if err := my_conf.LoadEnv(); err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(my_conf.DBURL))
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Database(my_conf.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		Users:    db.NewMongoUserStore(client),
		Hotels:   hotelStore,
		Rooms:    db.NewMongoRoomStore(client, hotelStore),
		Bookings: db.NewMongoBookingStorage(client),
	}
	var rooms []*types.Room
	var hotels []*types.Hotel

	user := api.AddUser(store, "Iak", "Vapvapvv", false)
	fmt.Println("IakJWT -> ", api.CreateTokenFromUser(user))
	admin := api.AddUser(store, "admin", "admin", true)
	fmt.Println("adminJWT -> ", api.CreateTokenFromUser(admin))
	hotels = append(hotels, api.AddHotel(store, "hotel1", "Moscow", 3.5, nil))
	for i := 0; i < 100; i++ {
		hotels = append(hotels, api.AddHotel(store, fmt.Sprintf("hotel%d", i+1), "Moscow", 0.5+float64(rand.Intn(5)), nil))
	}
	rooms = append(rooms, api.AddRoom(store, "large", true, 199.99, hotels[0].ID))
	rooms = append(rooms, api.AddRoom(store, "small", false, 29.99, hotels[0].ID))
	booking := api.AddBooking(store, user.ID, rooms[0].ID, 2, time.Now(), time.Now().AddDate(0, 0, 3))
	fmt.Println("booking id -> ", booking.ID)
}
