package api

import (
	"hotelReservation_golang/db"
	"hotelReservation_golang/types"
)

type HotelsSourceResp struct {
	ResourceResp
	Hotels []*types.Hotel
}

type RoomsSourceResp struct {
	ResourceResp
	Rooms []*types.Room
}

//
//type BookingsSourceResp struct {
//	ResourceResp
//	Bookings []*types.Booking
//}

type ResourceResp struct {
	Results int `json:"results"`
	Page    int `json:"page"`
}

type HotelQueryParams struct {
	db.PaginationFilter
	Rating float64
}

type RoomQueryParams struct {
	db.PaginationFilter
	Seaside bool
	Size    string
}
