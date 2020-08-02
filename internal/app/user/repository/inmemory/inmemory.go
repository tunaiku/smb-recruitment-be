package inmemory

import (
	"time"

	"github.com/tunaiku/mobilebanking/internal/app/domain"
	"golang.org/x/crypto/bcrypt"
)

type InMemoryUserRepository struct {
	datastore map[string]*domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	repo := new(InMemoryUserRepository)
	repo.init()
	return repo
}

func (inmem *InMemoryUserRepository) init() {
	password, err := bcrypt.GenerateFromPassword([]byte("123456"), 6)
	if err != nil {
		panic(err)
	}
	inmem.datastore = map[string]*domain.User{
		"fc55e3a8-c0fb-40c7-ab8a-9cda3fca40d4": {
			AccountReference: "10001",
			ConfiguredTransactionCredential: &domain.ConfiguredCredential{
				Pin: &domain.PinCredential{
					Pin: "111111",
				},
			},
			ID:       "fc55e3a8-c0fb-40c7-ab8a-9cda3fca40d4",
			JoinDate: time.Now(),
			Name:     "John Doe",
			Username: "john",
			Password: string(password),
		},
		"44c65528-950f-473f-ba69-00f28bc41f70": {
			AccountReference: "10002",
			ConfiguredTransactionCredential: &domain.ConfiguredCredential{
				Otp: &domain.OtpCredential{
					PhoneNumber: "081955334411",
				},
			},
			JoinDate: time.Now(),
			ID:       "44c65528-950f-473f-ba69-00f28bc41f70",
			Name:     "Jane Doe",
			Username: "jane",
			Password: string(password),
		},
	}
}

func (inmem *InMemoryUserRepository) LoadUser(id string) (*domain.User, error) {
	return inmem.datastore[id], nil
}

func (inmem *InMemoryUserRepository) LoadByUsername(username string) (*domain.User, error) {
	for _, value := range inmem.datastore {
		if value.Username == username {
			return value, nil
		}
	}
	return nil, domain.ErrUserNotFound
}
