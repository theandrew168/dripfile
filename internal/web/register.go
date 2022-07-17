package web

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/postgresql"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/task"
	"github.com/theandrew168/dripfile/internal/validator"
)

type registerForm struct {
	validator.Validator
	Email    string
	Password string
}

type registerData struct {
	Form registerForm
}

func (app *Application) handleRegister(w http.ResponseWriter, r *http.Request) {
	page := "site/auth/register.html"
	data := registerData{}
	app.render(w, r, page, data)
}

func (app *Application) handleRegisterForm(w http.ResponseWriter, r *http.Request) {
	page := "site/auth/register.html"

	form := registerForm{
		Email:    r.PostForm.Get("Email"),
		Password: r.PostForm.Get("Password"),
	}

	form.CheckRequired(form.Email, "Email")
	form.CheckRequired(form.Password, "Password")

	if !form.Valid() {
		data := registerData{
			Form: form,
		}
		app.render(w, r, page, data)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// ensure email isn't already taken
	_, err = app.store.Account.ReadByEmail(form.Email)
	if err == nil || !errors.Is(err, postgresql.ErrNotExist) {
		form.SetFieldError("Email", "An account with this email already exists")

		data := registerData{
			Form: form,
		}
		app.render(w, r, page, data)
		return
	}

	// create new project and new account within a single transaction
	var project core.Project
	var account core.Account
	err = app.store.WithTransaction(func(store *storage.Storage) error {
		// create project for the new account
		project = core.NewProject()
		err := store.Project.Create(&project)
		if err != nil {
			return err
		}

		// create the new account
		account = core.NewAccount(form.Email, string(hash), core.RoleOwner, project)
		err = store.Account.Create(&account)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		// check for TOCTOU race on account email
		if errors.Is(err, postgresql.ErrExist) {
			form.SetFieldError("Email", "An account with this email already exists")

			data := registerData{
				Form: form,
			}
			app.render(w, r, page, data)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	// generate a fresh session ID
	sessionID, err := GenerateSessionID()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	sessionHash := fmt.Sprintf("%x", sha256.Sum256([]byte(sessionID)))
	expiry := time.Now().AddDate(0, 0, 7)

	// create session model and store in the database
	session := core.NewSession(sessionHash, expiry, account)
	err = app.store.Session.Create(&session)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// send welcome email
	t, err := task.NewEmailSendTask(
		"DripFile",
		"info@dripfile.com",
		account.Email,
		account.Email,
		"Welcome to DripFile!",
		"Thanks for signing up with DripFile! I hope this adds some value.",
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// submit email task
	_, err = app.asynqClient.Enqueue(t)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// set cookie (just a session cookie after registration)
	cookie := NewSessionCookie(sessionIDCookieName, sessionID)
	http.SetCookie(w, &cookie)

	app.logger.Info("account create", map[string]string{
		"project_id": session.Account.Project.ID,
		"account_id": session.Account.ID,
	})
	app.logger.Info("account login", map[string]string{
		"project_id": session.Account.Project.ID,
		"account_id": session.Account.ID,
	})

	// redirect to dashboard
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
