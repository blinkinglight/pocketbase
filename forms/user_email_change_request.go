package forms

import (
	"github.com/blinkinglight/pocketbase/core"
	"github.com/blinkinglight/pocketbase/mails"
	"github.com/blinkinglight/pocketbase/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// UserEmailChangeRequest defines a user email change request form.
type UserEmailChangeRequest struct {
	app  core.App
	user *models.User

	NewEmail string `form:"newEmail" json:"newEmail"`
}

// NewUserEmailChangeRequest creates a new user email change request form.
func NewUserEmailChangeRequest(app core.App, user *models.User) *UserEmailChangeRequest {
	return &UserEmailChangeRequest{
		app:  app,
		user: user,
	}
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *UserEmailChangeRequest) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.NewEmail,
			validation.Required,
			validation.Length(1, 255),
			is.Email,
			validation.By(form.checkUniqueEmail),
		),
	)
}

func (form *UserEmailChangeRequest) checkUniqueEmail(value any) error {
	v, _ := value.(string)

	if !form.app.Dao().IsUserEmailUnique(v, "") {
		return validation.NewError("validation_user_email_exists", "User email already exists.")
	}

	return nil
}

// Submit validates and sends the change email request.
func (form *UserEmailChangeRequest) Submit() error {
	if err := form.Validate(); err != nil {
		return err
	}

	return mails.SendUserChangeEmail(form.app, form.user, form.NewEmail)
}
