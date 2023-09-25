package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"hotelReservation_golang/my_conf"
	"hotelReservation_golang/types"
)

type BookingStorage interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, map[string]any) ([]*types.Booking, error)
	GetBookingByID(ctx context.Context, id string) (*types.Booking, error)
	UpdateBooking(context.Context, string, map[string]any) error
}

type MongoBookingStorage struct {
	client *mongo.Client
	coll   *mongo.Collection

	BookingStorage
}

func NewMongoBookingStorage(client *mongo.Client) *MongoBookingStorage {
	return &MongoBookingStorage{
		client: client,
		coll:   client.Database(my_conf.DBNAME).Collection("bookings"),
	}
}

func (s *MongoBookingStorage) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	resp, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = resp.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (s *MongoBookingStorage) GetBookings(ctx context.Context, filter map[string]any) ([]*types.Booking, error) {
	curr, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var bookings []*types.Booking
	if err := curr.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}
func (s *MongoBookingStorage) GetBookingByID(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var booking *types.Booking
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *MongoBookingStorage) UpdateBooking(ctx context.Context, id string, update map[string]any) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	m := bson.M{"$set": update}

	_, err = s.coll.UpdateByID(ctx, oid, m)
	if err != nil {
		return err
	}

	return nil
}
