package db

type PaginationFilter struct {
	Limit int64
	Page  int64
}

type Store struct {
	Users    UserStorage
	Hotels   HotelStore
	Rooms    RoomStore
	Bookings BookingStorage
}
