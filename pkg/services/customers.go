package services

import (
	"time"
	"github.com/funlake/wschat/pkg/orm"
	"github.com/funlake/wschat/pkg/database"
)

// CustomerService handles customer-related operations
type CustomerService struct {
	
}

// CreateCustomer creates a new customer record
func (s *CustomerService) CreateCustomer(customerID string, name string) (*orm.Customer, error) {
	customer := &orm.Customer{
		CustomerID:   customerID,
		CustomerName: name,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := database.GetDB().Create(customer).Error; err != nil {
		return nil, err
	}
	return customer, nil
}

// GetCustomerByID retrieves a customer by their ID
func (s *CustomerService) GetCustomerByID(customerID string) (*orm.Customer, error) {
	var customer orm.Customer
	if err := database.GetDB().Where("customer_id = ?", customerID).First(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}
// GetCustomerByName retrieves a customer by their name
func (s *CustomerService) GetCustomerByName(name string) (*orm.Customer, error) {
	var customer orm.Customer
	if err := database.GetDB().Where("customer_name = ?", name).First(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

// UpdateCustomer updates customer information
func (s *CustomerService) UpdateCustomer(customerID string, name string) (*orm.Customer, error) {
	customer, err := s.GetCustomerByID(customerID)
	if err != nil {
		return nil, err
	}

	customer.CustomerName = name
	customer.UpdatedAt = time.Now()

	if err := database.GetDB().Save(customer).Error; err != nil {
		return nil, err
	}
	return customer, nil
}

// DeleteCustomer removes a customer record
func (s *CustomerService) DeleteCustomer(customerID string) error {
	return database.GetDB().Where("customer_id = ?", customerID).Delete(&orm.Customer{}).Error
}

// ListCustomers returns all customers
func (s *CustomerService) ListCustomers() ([]orm.Customer, error) {
	var customers []orm.Customer
	if err := database.GetDB().Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}
