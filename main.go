package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelReservation_golang/api"
	"hotelReservation_golang/db"
	"hotelReservation_golang/middleware"
	"hotelReservation_golang/my_conf"
	"log"
	"os"
)

// Configuration
// 1.MongoDB endpoint
// 2. ListenAddress of HTTP server
// 3. JWT secret
// 4.MongoDBName

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	if err := my_conf.LoadEnv(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(my_conf.DBNAME)
	fmt.Println(my_conf.DBURL)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(my_conf.DBURL))
	if err != nil {
		log.Fatal(err)
	}

	// handlers initialization
	var (
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStorage(client)
		store        = &db.Store{
			Users:    userStore,
			Hotels:   hotelStore,
			Rooms:    roomStore,
			Bookings: bookingStore,
		}
		hotelHandler   = api.NewHotelHandler(store)
		userHandler    = api.NewUserHandler(store)
		authHandler    = api.NewAuthHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)

		app   = fiber.New(config)
		auth  = app.Group("/api/")
		apiv1 = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		admin = apiv1.Group("/admin", middleware.AdminAuth)
	)

	// authentication
	auth.Post("/auth", authHandler.HandleAuthenticate)

	//Versioned API handlers
	// user handlers
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	// hotel handlers
	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotels/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotels/:id/rooms", hotelHandler.HandleGetRooms)
	apiv1.Delete("/hotels/:id", hotelHandler.HandleDeleteHotel)

	// rooms handlers
	apiv1.Post("/rooms/:id/book", roomHandler.HandleBookRoom)
	apiv1.Get("/rooms", roomHandler.HandleGetRooms)

	// TODO: handle cancel booking

	// bookings handlers
	apiv1.Get("/bookings/:id", bookingHandler.HandleGetBooking)

	// admin handlers
	admin.Get("/bookings", bookingHandler.HandleGetBookings)

	apiv1.Put("/bookings/:id", bookingHandler.HandleCancelBookings)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	app.Listen(listenAddr)
}
