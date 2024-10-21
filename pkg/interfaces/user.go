package interfaces

import "split-expenses/pkg/models"

type UserCreationRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Contact   string `json:"contact"`
}

type UserLoginResponse struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Token     string `json:"token"`
	Contact   string `json:"contact"`
	Id        string `json:"id"`
}

func (c UserCreationRequest) AsUser() models.User {
	return models.User{
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Email:     c.Email,
		Password:  c.Password,
		Contact:   c.Contact,
	}
}
