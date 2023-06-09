package db

import (
	"fmt"
	"log"

	"github.com/itsshashank/identity-reconciliation/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserStorer interface {
	PutOrder(*types.Request) error
	GetContacts(*types.Request) (*types.Response, error)
}

type PostgresUserStore struct {
	db *gorm.DB
}

func NewPostgresUserStore(dsn string) *PostgresUserStore {
	pg, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// AutoMigrate the models
	err = pg.AutoMigrate(&types.Contact{})
	if err != nil {
		log.Fatal(err)
	}
	return &PostgresUserStore{
		db: pg,
	}
}

func (s *PostgresUserStore) PutOrder(req *types.Request) error {
	// Query the database to find the primary contact
	var existingContact types.Contact
	err := s.db.Where("(email = ? OR phone_number = ?) and link_precedence = 'primary'", req.Email, req.PhoneNumber).
		Order("created_at ASC").
		First(&existingContact).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to query the database")
	}

	// Create a new primary contact if no existing contact is found
	if existingContact.ID == 0 {
		primaryContact := types.Contact{
			Email:          req.Email,
			PhoneNumber:    req.PhoneNumber,
			LinkPrecedence: "primary",
		}
		err = s.db.Create(&primaryContact).Error
		if err != nil {
			return fmt.Errorf("failed to create contact")
		}
		return nil
	}

	// Create a new secondary contact linked to the existing primary contact
	secondaryContact := types.Contact{
		Email:          req.Email,
		PhoneNumber:    req.PhoneNumber,
		LinkedID:       existingContact.ID,
		LinkPrecedence: "secondary",
	}
	err = s.db.Create(&secondaryContact).Error
	if err != nil {
		return fmt.Errorf("failed to create contact")
	}
	return nil
}

func (s *PostgresUserStore) GetContacts(req *types.Request) (*types.Response, error) {
	// Query the database to find the primary contact
	var primaryContact types.Contact
	err := s.db.Where("(email = ? OR phone_number = ?) and link_precedence = 'primary'", req.Email, req.PhoneNumber).
		Order("created_at ASC").
		First(&primaryContact).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to query the database")
		}
		secondary := primaryContact
		primaryContact = types.Contact{}
		err = s.db.Where("email = ? OR phone_number = ?", req.Email, req.PhoneNumber).
			Order("created_at ASC").
			First(&secondary).Error
		if err != nil {
			return nil, fmt.Errorf("failed to query the database")
		}
		err = s.db.Where("id = ?", secondary.LinkedID).Find(&primaryContact).Error
		if err != nil {
			return nil, fmt.Errorf("failed to query the database")
		}
	}

	// Query the database to find secondary contacts linked to the primary contact
	var secondaryContacts []types.Contact
	err = s.db.Where("linked_id = ?", primaryContact.ID).Find(&secondaryContacts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query the database")
	}
	if len(secondaryContacts) == 0 {
		err = s.db.Where("(email = ? OR phone_number = ?) and id != ?", req.Email, req.PhoneNumber, primaryContact.ID).Find(&secondaryContacts).Error
		if err != nil {
			return nil, fmt.Errorf("failed to query the database")
		}
		for _, sc := range secondaryContacts {
			err = s.db.Model(&sc).Updates(types.Contact{LinkedID: primaryContact.ID, LinkPrecedence: "secondary"}).Error
			if err != nil {
				return nil, fmt.Errorf("failed to update the database %s", err)
			}
		}
	}

	return &types.Response{
		PrimaryContactID:    primaryContact.ID,
		Emails:              getEmails(primaryContact, secondaryContacts),
		PhoneNumbers:        getPhoneNumbers(primaryContact, secondaryContacts),
		SecondaryContactIDs: getSecondaryContactIds(secondaryContacts),
	}, nil
}

func containsString(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false || str == ""
}

func getEmails(primary types.Contact, secondary []types.Contact) []string {
	emails := make([]string, 0, len(secondary)+1)
	emails = append(emails, primary.Email)
	for _, c := range secondary {
		if !containsString(emails, c.Email) {
			emails = append(emails, c.Email)
		}
	}
	return emails
}

func getPhoneNumbers(primary types.Contact, secondary []types.Contact) []string {
	phoneNumbers := make([]string, 0, len(secondary)+1)
	phoneNumbers = append(phoneNumbers, primary.PhoneNumber)
	for _, c := range secondary {
		if !containsString(phoneNumbers, c.PhoneNumber) {
			phoneNumbers = append(phoneNumbers, c.PhoneNumber)
		}
	}
	return phoneNumbers
}

func getSecondaryContactIds(secondary []types.Contact) []int {
	ids := make([]int, len(secondary))
	for i, c := range secondary {
		ids[i] = c.ID
	}
	return ids
}
