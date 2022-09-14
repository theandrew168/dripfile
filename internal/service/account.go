package service

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/task"
)

func (s *Service) CreateAccount(email, password string) (model.Account, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.Account{}, err
	}

	// create the new account
	account := model.NewAccount(email, string(hash), model.RoleAdmin)
	err = s.store.Account.Create(&account)
	if err != nil {
		return model.Account{}, err
	}

	// send welcome email
	t := task.NewEmailSendTask(
		"DripFile",
		"info@dripfile.com",
		"",
		account.Email,
		"Welcome to DripFile!",
		"Thanks for signing up with DripFile! I hope this adds some value.",
	)
	err = s.queue.Submit(t)
	if err != nil {
		return model.Account{}, err
	}

	return account, nil
}
