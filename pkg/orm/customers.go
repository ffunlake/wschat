package orm

import "time"

type Customer struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	CustomerID   string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	CustomerName string    `gorm:"type:varchar(100)"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}

func (Customer) TableName() string {
	return "customers"
}
