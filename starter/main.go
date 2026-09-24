package main

import (
	"log"

	"github.com/tgo-framework/tgo/pkg/app"
	"github.com/tgo-framework/tgo/pkg/config"
	"github.com/tgo-framework/tgo/starter/app/handlers"
	"github.com/tgo-framework/tgo/starter/proto/v1/userv1/userconnect"
)

func main() {
	// 1. Load Configuration
	_ = config.MustLoad()

	// 2. Initialize App Container
	application := app.New().SetAddr(":8080")

	// 3. Register ConnectRPC Services
	userHandler := handlers.NewUserHandler()
	path, handler := userconnect.NewUserServiceHandler(userHandler)
	application.Server().Register(path, handler)

	// 4. Run the application (blocks until exit)
	if err := application.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
