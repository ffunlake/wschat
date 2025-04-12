package services

import (
	"time"
	"github.com/funlake/wschat/pkg/orm"
	"github.com/funlake/wschat/pkg/database"
)

// FeedbackService handles feedback-related operations
type FeedbackService struct {
}

// CreateFeedback creates a new feedback record
func (s *FeedbackService) CreateFeedback(customerID string, rating int, comment string) (*orm.Feedback, error) {
	feedback := &orm.Feedback{
		CustomerID: customerID,
		Rating:     rating,
		Comment:    comment,
		Timestamp:  time.Now(),
	}

	if err := database.GetDB().Create(feedback).Error; err != nil {
		return nil, err
	}
	return feedback, nil
}

// GetFeedbackByCustomerID retrieves feedback by customer ID
func (s *FeedbackService) GetFeedbackByCustomerID(customerID string) ([]orm.Feedback, error) {
	var feedbacks []orm.Feedback
	if err := database.GetDB().Where("customer_id = ?", customerID).Order("timestamp desc").Find(&feedbacks).Error; err != nil {
		return nil, err
	}
	return feedbacks, nil
} 