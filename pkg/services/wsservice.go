package services

import (
	"log"
	"github.com/gorilla/websocket"
	"strings"
	"fmt"
	"strconv"
	"sync"
	"time"
)
// init services
var (
	customerService = &CustomerService{}
	messageService = &MessageService{}
	feedbackService = &FeedbackService{}
	aiService *AiService
	connCustomer sync.Map
)
type WsServeice struct {}
func (wsServeice *WsServeice) HandleMessage(conn *websocket.Conn, message string) {
	if aiService == nil {
		aiService = GetAiService()
	}
	// Split message into command and content
	parts := strings.SplitN(message, " ", 2)
	command := parts[0]
	var content string
	if len(parts) > 1 {
		content = parts[1]
	} else {
		content = command
	}

	switch command {
	default:
		_ = conn.WriteMessage(websocket.TextMessage, []byte("[error] Invalid command"))
		return
	case "/register":
		// Register new customer
		if content == "" {
			err := conn.WriteMessage(websocket.TextMessage, []byte("[error] Please provide a customer name"))
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}
			return
		}
		_, ok := connCustomer.Load(conn)
		if ok {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("[error] Please logout first"))
			return
		}
		// Check if customer already exists
		existingCustomer, err := customerService.GetCustomerByName(content)
		if err == nil && existingCustomer != nil {
			err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("[error] Customer %s is already registered", content)))
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}
			return
		}
		customer, err := customerService.CreateCustomer(newCustomerID(), content)
		if err != nil {
			err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("[error] Failed to register: %v", err)))
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}
			log.Printf("Failed to register: %v", err)
		}

		err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Successfully registered as %s with ID %s", content, customer.CustomerID)))
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
		// store the relation between conn and customer id
		connCustomer.Store(conn, customer.CustomerID)

	case "/login":
		_, ok := connCustomer.Load(conn)
		if ok {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("[error] Please logout first"))
			return
		}
		if content == "" {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("Please provide your customer name"))
		}
		customer, err := customerService.GetCustomerByName(content)
		if err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("[error] Login failed: %v", err)))
			return
		}
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Welcome back %s!", customer.CustomerName)))
		connCustomer.Store(conn, customer.CustomerID)
	case "/logout":
		connCustomer.Delete(conn)
		_ = conn.WriteMessage(websocket.TextMessage, []byte("Successfully logged out"))
		return
	case "/feedback":
		customerID, ok := connCustomer.Load(conn)
		if !ok {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("[error] Please register or login first"))
			return
		}
		if content == "" {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("Please provide both rating (1-5) and comment"))
			return
		}
		
		parts := strings.SplitN(content, " ", 2)
		if len(parts) < 2 {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("Please provide both rating (1-5) and comment"))
			return
		}
		
		rating, err := strconv.Atoi(parts[0])
		if err != nil || rating < 1 || rating > 5 {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("Rating must be a number between 1 and 5"))
			return
		}
		
		// 存储反馈到数据库
		customerIDStr, _ := customerID.(string)
		feedback, err := feedbackService.CreateFeedback(customerIDStr, rating, parts[1])
		if err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Failed to save feedback: %v", err)))
			return
		}
		
		_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Thank you for your feedback! Rating: %d", feedback.Rating)))
	
	case "/history":
		customerID, ok := connCustomer.Load(conn)
		if !ok {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("[error] Please register or login first"))
			return
		}

		// Parse limit from content
		limit := 10 // Default limit
		if content != "" {
			parsedLimit, err := strconv.Atoi(content)
			if err == nil && parsedLimit > 0 {
				limit = parsedLimit
			}
		}

		customerIDStr, _ := customerID.(string)
		history, err := messageService.GetMessageHistory(customerIDStr, limit)
		if err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Failed to get message history: %v", err)))
			return
		}

		// Build a table header
		tableHeader := "\n+------------+----------+-----------------+\n"
		tableHeader += "| Timestamp  | Sender   | Message         |\n" 
		tableHeader += "+------------+----------+-----------------+\n"
		
		// Build table rows
		var tableContent string
		for _, msg := range history {
			timestamp := msg.Timestamp.Format("15:04:05")
			// Truncate message if too long
			messageContent := msg.Message
			if len(messageContent) > 15 {
				messageContent = messageContent[:12] + "..."
			}
			tableContent += fmt.Sprintf("| %-10s | %-8s | %-15s |\n", timestamp, msg.Sender, messageContent)
		}
		
		// Add bottom border
		tableFooter := "+------------+----------+-----------------+\n"
		
		// Send the complete table
		_ = conn.WriteMessage(websocket.TextMessage, []byte(tableHeader + tableContent + tableFooter))
		return
	case "/chat":
		customerID, ok := connCustomer.Load(conn)
		if !ok {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("[error] Please register or  login first"))
			return
		}
		
		if content == "" {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("Message cannot be empty"))
			return
		}
		//Bonus 1: Sentiment Analysis → Classify feedback as positive, neutral, or negative using AI.
		go func() {
			sentimentAnalysis := aiService.SentimentAnalysis(content)
			if sentimentAnalysis != "" {
				emoji := ":|"
				if sentimentAnalysis == "Positive" {
					emoji = ":)"
				} else if sentimentAnalysis == "Negative" {
					emoji = ":("
				}
				_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Ai Sentiment Analysis results: [%s %s] %s", emoji, sentimentAnalysis, content)))
			}
		}()
		//response to client	
		_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Got your message: %s", content)))
		// customer messages storage
		customerIDStr, _ := customerID.(string)
		_, err := messageService.CreateMessage(customerIDStr, content, "customer")
		if err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Failed to save message: %v", err)))
			return
		}

		_, err = messageService.CreateMessage(customerIDStr, content, "bot")
		if err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Failed to save message: %v", err)))
			return
		}
	case "/aichat":
		_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", "Please wait a second for the answer...")))
		content = aiService.GenerateAIResponse(content)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", content)))
	}
}

func newCustomerID() string {
	// Generate client ID based on current timestamp
	return time.Now().Format("20060102150405")
}