package orm

import "time"

// Feedback represents customer feedback
type Feedback struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	CustomerID string    `gorm:"type:varchar(100);not null;index"`
	Rating     int       `gorm:"not null"` // 1-5 rating
	Comment    string    `gorm:"type:text"`
	Timestamp  time.Time `gorm:"not null"`
	Customer   Customer  `gorm:"foreignKey:CustomerID;references:CustomerID"`
}

func (Feedback) TableName() string {
	return "feedback"
}
