package cmd

import (
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"bufio"
	"os"
	"strings"
	"github.com/funlake/wschat/config"
	"time"
	"github.com/spf13/viper"
	"log"
)

func init() {
	clientCommand.PersistentFlags().StringVarP(&configFile, "config", "c", "./config/chat.yaml", "Start client with config file")
}
var clientCommand = &cobra.Command{
	Use:   "client",
	Short: "Create a new chat session",
	Run: func(cmd *cobra.Command, args []string) {
		// 配置日志输出格式，不显示时间、文件信息等
		log.SetFlags(0)
		log.SetOutput(os.Stdout)
		
		//load config
		err := config.LoadConfig(configFile)
		if err != nil {
			log.Printf("Error loading config: %v", err)
			return
		}
		
		// 创建调试日志记录器，用于记录详细信息到文件
		debugLog, err := os.OpenFile("client_debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			// 配置详细日志到文件
			debugLogger := log.New(debugLog, "", log.LstdFlags|log.Lshortfile)
			defer debugLog.Close()
			
			// 输出到控制台和文件
			debugLogger.Println("Client started")
		}
		
		// Connect to WebSocket server
		websocketURL := viper.GetString("client.websocket.url")
		readBufferSize := viper.GetInt("client.websocket.read_buffer_size")
		writeBufferSize := viper.GetInt("client.websocket.write_buffer_size")
		
		// 创建WebSocket拨号器
		dialer := websocket.Dialer{
			ReadBufferSize:  readBufferSize,
			WriteBufferSize: writeBufferSize,
		}
		
		// 设置连接超时
		dialer.HandshakeTimeout = time.Duration(10) * time.Second
		
		log.Printf("Connecting to %s...", websocketURL)
		conn, _, err := dialer.Dial(websocketURL, nil)
		if err != nil {
			log.Printf("Failed to connect to WebSocket server: %v", err)
			return
		}
		defer conn.Close()	
		log.Println("Connected to server successfully!")
		
		// Display available commands to the user
		displayHelp()

		// Start a goroutine to read messages from server
		go func() {
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Printf("Error reading from server: %v", err)
					os.Exit(1)
					return
				}
				log.Printf("\n> Bot response: %s", message)
			}
		}()

		// Main loop for sending messages
		for {
			// receive message from command line
			reader := bufio.NewReader(os.Stdin)
			// read string until newline
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error reading input: %v", err)
				return
			}
			input = strings.TrimSpace(input)
			// Check if input starts with valid command prefix

			err = conn.WriteMessage(websocket.TextMessage, []byte(input))
			if err != nil {
				log.Printf("Error sending message: %v", err)
				return
			}
		}
	},
}
// Helper function to display available commands
func displayHelp() {
	log.Println("\n+------------------------ Help Menu ------------------------+")
	log.Println("Available commands:")
	log.Println("/register <name>  - Register as a new user")
	log.Println("/login <name>     - Login with existing name") 
	log.Println("/logout           - Logout from the chat server")
	log.Println("/history [limit]         - View history chat messages")
	log.Println("/chat <text>   - Send a message")
	log.Println("/aichat <text>          - Communicate with AI")
	log.Println("/feedback <rating> <comment> - Send feedback")
	log.Println("+--------------------------------------------------------------+")
}