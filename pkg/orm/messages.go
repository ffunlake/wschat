package orm

import "time"

// Message represents a chat message
type Message struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	CustomerID string    `gorm:"type:varchar(100);not null;index"`
	Message    string    `gorm:"type:text;not null"`
	Sender     string    `gorm:"type:varchar(20);not null"` // 'customer' or 'bot'
	Timestamp  time.Time `gorm:"not null"`
	Customer   Customer  `gorm:"foreignKey:CustomerID;references:CustomerID"`
}

func (Message) TableName() string {
	return "messages"
}
