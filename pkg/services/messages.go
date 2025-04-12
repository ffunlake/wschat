package services

import (
	"time"
	"github.com/funlake/wschat/pkg/orm"
	"github.com/funlake/wschat/pkg/database"
)

// MessageService handles message-related operations
type MessageService struct {
}

// CreateMessage creates a new message record
func (s *MessageService) CreateMessage(customerID string, content string, sender string) (*orm.Message, error) {
	message := &orm.Message{
		CustomerID: customerID,
		Message:    content,
		Sender:     sender, // "customer" or "bot"
		Timestamp:  time.Now(),
	}

	if err := database.GetDB().Create(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

// GetMessagesByCustomerID retrieves messages by customer ID
func (s *MessageService) GetMessagesByCustomerID(customerID string) ([]orm.Message, error) {
	var messages []orm.Message
	if err := database.GetDB().Where("customer_id = ?", customerID).Order("timestamp desc").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// GetMessageHistory retrieves message history by customer ID
func (s *MessageService) GetMessageHistory(customerID string,limit int) ([]orm.Message, error) {
	var messages []orm.Message
	if err := database.GetDB().Where("customer_id = ?", customerID).Order("timestamp desc").Limit(limit).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
