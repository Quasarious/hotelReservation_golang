package db

const (
	DBNAME     = "hotel-reservation"
	TestDBNAME = "hotel-reservation-test"
	DBURI      = "mongodb://localhost:27017"
)

type Store struct {
	Users  UserStorage
	Hotels HotelStore
	Rooms  RoomStore
}
