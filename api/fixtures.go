package api

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelReservation_golang/db"
	"hotelReservation_golang/types"
	"log"
	"time"
)

func AddBooking(store *db.Store, uid, rid primitive.ObjectID, np int, fd, td time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:     uid,
		RoomID:     rid,
		NumPersons: np,
		FromDate:   fd,
		TillDate:   td,
	}

	insertedBooking, err := store.Bookings.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}

	return insertedBooking
}

func AddRoom(store *db.Store, size string, seaside bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: seaside,
		Price:   price,
		HotelID: hotelID,
	}

	insertedRoom, err := store.Rooms.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func AddHotel(store *db.Store, name, loc string, rating float64, rooms []primitive.ObjectID) *types.Hotel {
	var roomIDS = rooms
	if rooms == nil {
		roomIDS = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomIDS,
		Rating:   rating,
	}
	insertedHotel, err := store.Hotels.InsertHotel(context.Background(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddUser(store *db.Store, fn, ln string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.ru", fn, ln),
		FirstName: fn,
		LastName:  ln,
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin
	insertedUser, err := store.Users.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s -> %s\n", user.Email, CreateTokenFromUser(user))

	return insertedUser
}
