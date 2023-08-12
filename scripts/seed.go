package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelReservation_golang/db"
	"hotelReservation_golang/types"
	"log"
)

var (
	ctx        = context.Background()
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStorage
)

func seedHotel(name string, location string, rating float64) {
	hotel := types.Hotel{
		ID:       primitive.NewObjectID(),
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{Size: "small", Price: 99.9, HotelID: hotel.ID},
		{Size: "normal", Price: 199.9, HotelID: hotel.ID},
		{Size: "kingsize", Price: 259.99, HotelID: hotel.ID},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(insertedRoom)

	}
	fmt.Println(insertedHotel)
}

func initDB() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}

func seedUser(firstName, lastName, email string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  "123456!",
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()
	seedHotel("KremlinHotel", "Moscow", 3.2)
	seedHotel("BigBenHotel", "London", 4.3)
	seedHotel("Ratatoir", "Paris", 4.5)
	seedUser("Ivan", "Ivanov", "ivan@mail.ru")
}
