package web

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/form"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/task"
)

func (app *Application) handleRegister(w http.ResponseWriter, r *http.Request) {
	page := "site/auth/register.page.html"

	data := struct {
		Form *form.Form
	}{
		Form: form.New(nil),
	}

	app.render(w, r, page, data)
}

func (app *Application) handleRegisterForm(w http.ResponseWriter, r *http.Request) {
	page := "site/auth/register.page.html"

	f := form.New(r.PostForm)
	f.Required("email", "password")

	data := struct {
		Form *form.Form
	}{
		Form: f,
	}

	if !f.Valid() {
		app.render(w, r, page, data)
		return
	}

	email := f.Get("email")
	password := f.Get("password")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// ensure email isn't already taken
	_, err = app.store.Account.ReadByEmail(email)
	if err == nil || !errors.Is(err, database.ErrNotExist) {
		f.Errors.Add("email", "An account with this email already exists")
		app.render(w, r, page, data)
		return
	}

	// create Stripe customer
	customerID, err := app.billing.CreateCustomer(email)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// create new project and new account within a single transaction
	var project core.Project
	var account core.Account
	err = app.store.WithTransaction(func(store *storage.Storage) error {
		// create project for the new account
		project = core.NewProject(customerID)
		err := store.Project.Create(&project)
		if err != nil {
			return err
		}

		// create the new account
		account = core.NewAccount(email, string(hash), core.RoleOwner, project)
		err = store.Account.Create(&account)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, database.ErrExist) {
			f.Errors.Add("email", "An account with this email already exists")
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
	_, err = app.queue.Enqueue(t)
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

	// redirect to billing setup
	http.Redirect(w, r, "/billing/setup", http.StatusSeeOther)
}
