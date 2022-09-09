package service

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/task"
)

func (s *Service) CreateAccount(email, password string) (model.Account, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.Account{}, err
	}

	// create new project and new account within a single transaction
	var project model.Project
	var account model.Account
	err = s.store.WithTransaction(func(store *storage.Storage) error {
		// create project for the new account
		project = model.NewProject()
		err := store.Project.Create(&project)
		if err != nil {
			return err
		}

		// create the new account
		account = model.NewAccount(email, string(hash), model.RoleOwner, project)
		err = store.Account.Create(&account)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return model.Account{}, nil
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
