package forms

import (
	"github.com/blinkinglight/pocketbase-mysql/core"
	"github.com/blinkinglight/pocketbase-mysql/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// UserEmailLogin defines a user email/pass login form.
type UserEmailLogin struct {
	app core.App

	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

// NewUserEmailLogin creates a new user email/pass login form.
func NewUserEmailLogin(app core.App) *UserEmailLogin {
	form := &UserEmailLogin{
		app: app,
	}

	return form
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *UserEmailLogin) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(&form.Email, validation.Required, validation.Length(1, 255), is.Email),
		validation.Field(&form.Password, validation.Required, validation.Length(1, 255)),
	)
}

// Submit validates and submits the form.
// On success returns the authorized user model.
func (form *UserEmailLogin) Submit() (*models.User, error) {
	if err := form.Validate(); err != nil {
		return nil, err
	}

	user, err := form.app.Dao().FindUserByEmail(form.Email)
	if err != nil {
		return nil, err
	}

	if !user.ValidatePassword(form.Password) {
		return nil, validation.NewError("invalid_login", "Invalid login credentials.")
	}

	return user, nil
}
