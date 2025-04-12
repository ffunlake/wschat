package cmd
import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/funlake/wschat/pkg/services"
	"github.com/funlake/wschat/pkg/database"
	"time"
	"encoding/json"
	"github.com/funlake/wschat/config"
	"github.com/spf13/viper"
	"log"
)

func init() {
	serverCommand.PersistentFlags().StringVarP(&configFile, "config", "c", "./config/chat.yaml", "Start client with config file")
}
var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "Create a new chat server",
	Run: func(cmd *cobra.Command, args []string) {
		// 配置日志格式
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		
		//load config
		err := config.LoadConfig(configFile)
		if err != nil {
			log.Printf("Error loading config: %v", err)
			return
		}
		log.Println("Starting WebSocket server...")
		// service init
		wsService := &services.WsServeice{}
		customerService := &services.CustomerService{}
		messageService := &services.MessageService{}
		
		defer database.CloseDB()
		// create websocket upgrader
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // allow all origins
			},
		}
		// define websocket handler
		http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Printf("Error upgrading connection: %v", err)
				return
			}
			defer conn.Close()

			log.Println("New client connected")
			
			for {
				// receive message from client
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Printf("Error reading message: %v", err)
					break
				}

				msg := string(message)
				wsService.HandleMessage(conn, msg)
			}
			
			log.Println("Connection closed")
		})
		// Handle GET requests for message history
		http.HandleFunc("/message/list", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Get customer name from query parameter
			customerName := r.URL.Query().Get("customer_name")
			if customerName == "" {
				http.Error(w, "customer_name parameter is required", http.StatusBadRequest)
				return
			}

			// Get customer by name
			customer, err := customerService.GetCustomerByName(customerName)
			if err != nil {
				http.Error(w, fmt.Sprintf("Customer not found: %v", err), http.StatusNotFound)
				return
			}

			// Get message history for customer
			messages, err := messageService.GetMessagesByCustomerID(customer.CustomerID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to get messages: %v", err), http.StatusInternalServerError)
				return
			}

			// Convert messages to JSON
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(messages); err != nil {
				http.Error(w, fmt.Sprintf("Failed to encode messages: %v", err), http.StatusInternalServerError)
				return
			}
		})
		// Get server configuration from viper
		host := viper.GetString("server.host")
		port := viper.GetInt("server.port")
		readTimeout := viper.GetInt("server.read_timeout")
		writeTimeout := viper.GetInt("server.write_timeout") 
		idleTimeout := viper.GetInt("server.idle_timeout")

		// Configure server timeouts
		server := &http.Server{
			Addr: fmt.Sprintf("%s:%d", host, port),
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
			IdleTimeout: time.Duration(idleTimeout) * time.Second,
			Handler: nil,
		}

		// Start HTTP server with configuration
		log.Printf("Http server is running on %s:%d", host, port)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Error: %v", err)
		}
		return
	},
}