package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tokalink/tgo/pkg/app"
	"github.com/tokalink/tgo/starter/app/handlers"
	"github.com/tokalink/tgo/starter/proto/v1/userv1/userconnect"
)

var (
	serveAddr string
	servePort string
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the TGo ConnectRPC & HTTP application server",
	Run: func(cmd *cobra.Command, args []string) {
		addr := serveAddr
		if servePort != "" {
			if strings.HasPrefix(servePort, ":") {
				addr = servePort
			} else {
				addr = fmt.Sprintf(":%s", servePort)
			}
		}

		log.Printf("[Craft] Starting TGo application on %s...\n", addr)

		application := app.New().SetAddr(addr)

		// Register default ConnectRPC services
		userHandler := handlers.NewUserHandler()
		path, handler := userconnect.NewUserServiceHandler(userHandler)
		application.Server().Register(path, handler)

		if err := application.Run(); err != nil {
			log.Fatalf("[Craft] Application error: %v", err)
		}
	},
}

func init() {
	ServeCmd.Flags().StringVarP(&serveAddr, "addr", "a", ":8080", "Server address to listen on")
	ServeCmd.Flags().StringVarP(&servePort, "port", "p", "", "Port to listen on (e.g. 8080 or 8989)")
}
