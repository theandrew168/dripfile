package web

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/postgresql"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/task"
	"github.com/theandrew168/dripfile/internal/validator"
)

type registerForm struct {
	validator.Validator `form:"-"`

	Email    string `form:"Email"`
	Password string `form:"Password"`
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

	var form registerForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r)
		return
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
	var project model.Project
	var account model.Account
	err = app.store.WithTransaction(func(store *storage.Storage) error {
		// create project for the new account
		project = model.NewProject()
		err := store.Project.Create(&project)
		if err != nil {
			return err
		}

		// create the new account
		account = model.NewAccount(form.Email, string(hash), model.RoleOwner, project)
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
	session := model.NewSession(sessionHash, expiry, account)
	err = app.store.Session.Create(&session)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
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
	err = app.queue.Push(t)
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
