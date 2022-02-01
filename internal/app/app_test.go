package app_test

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/chromedp/chromedp"

	"github.com/theandrew168/dripfile/internal/app"
	"github.com/theandrew168/dripfile/internal/log"
	"github.com/theandrew168/dripfile/internal/postgres"
	"github.com/theandrew168/dripfile/internal/test"
)

// chromedp docs:
// https://pkg.go.dev/github.com/chromedp/chromedp#Query

// chromedp context
var ctx context.Context

// test credentials
var (
	email    string
	username string
	password string
)

func TestMain(m *testing.M) {
	// skip UI tests if running with -short
	flag.Parse()
	if testing.Short() {
		os.Exit(0)
	}

	os.Exit(run(m))
}

func run(m *testing.M) int {
	// seed RNG
	rand.Seed(time.Now().UnixNano())

	// setup test credentials
	username = test.RandomString(16)
	password = username
	email = username + "@dripfile.com"

	// setup application deps
	cfg := test.Config()
	conn, err := postgres.Connect(cfg.DatabaseURI)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer conn.Close()

	storage := postgres.NewStorage(conn)
	logger := log.NewLogger(os.Stdout)

	// create the application
	handler := app.New(storage, logger)

	// start test server
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// start chromedp client
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// assign chromedp context
	ctx = taskCtx

	// navigate to test server
	err = chromedp.Run(ctx, chromedp.Navigate(ts.URL))
	if err != nil {
		fmt.Println(err)
		return 1
	}

	// run the tests!
	return m.Run()
}

func TestIndex(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.WaitVisible("#index", chromedp.ByID),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccountRegister(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.Click("#register", chromedp.ByID),

		chromedp.WaitVisible("#email", chromedp.ByID),
		chromedp.WaitVisible("#username", chromedp.ByID),
		chromedp.WaitVisible("#password", chromedp.ByID),

		chromedp.SendKeys("#email", email, chromedp.ByID),
		chromedp.SendKeys("#username", username, chromedp.ByID),
		chromedp.SendKeys("#password", password, chromedp.ByID),

		chromedp.WaitVisible("#submit", chromedp.ByID),
		chromedp.Submit("#submit", chromedp.ByID),
		chromedp.WaitVisible("#dashboard", chromedp.ByID),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccountLogout(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.Click("#dropdown", chromedp.ByID),

		chromedp.WaitVisible("#logout", chromedp.ByID),
		chromedp.Submit("#logout", chromedp.ByID),

		chromedp.WaitVisible("#index", chromedp.ByID),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccountLogin(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.Click("#login", chromedp.ByID),

		chromedp.WaitVisible("#email", chromedp.ByID),
		chromedp.WaitVisible("#password", chromedp.ByID),

		chromedp.SendKeys("#email", email, chromedp.ByID),
		chromedp.SendKeys("#password", password, chromedp.ByID),

		chromedp.WaitVisible("#submit", chromedp.ByID),
		chromedp.Submit("#submit", chromedp.ByID),

		chromedp.WaitVisible("#dashboard", chromedp.ByID),
	)
	if err != nil {
		t.Fatal(err)
	}
}

// TODO: check nothing listed, just center button
func TestLocationListEmpty(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.Click("#nav-location", chromedp.ByID),

		chromedp.WaitVisible("#location", chromedp.ByID),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLocationCreate(t *testing.T) {
	info := test.RandomString(32)
	err := chromedp.Run(ctx,
		chromedp.WaitVisible("#new-location", chromedp.ByID),
		chromedp.Click("#new-location", chromedp.ByID),

		chromedp.WaitVisible("#endpoint", chromedp.ByID),
		chromedp.WaitVisible("#access-key-id", chromedp.ByID),
		chromedp.WaitVisible("#secret-access-key", chromedp.ByID),
		chromedp.WaitVisible("#bucket-name", chromedp.ByID),
		chromedp.SendKeys("#endpoint", info, chromedp.ByID),
		chromedp.SendKeys("#access-key-id", info, chromedp.ByID),
		chromedp.SendKeys("#secret-access-key", info, chromedp.ByID),
		chromedp.SendKeys("#bucket-name", info, chromedp.ByID),

		chromedp.WaitVisible("#submit", chromedp.ByID),
		chromedp.Submit("#submit", chromedp.ByID),
		chromedp.WaitVisible("#location", chromedp.ByID),
	)
	if err != nil {
		t.Fatal(err)
	}
}

// TODO: check listed locations and top-right button
func TestLocationList(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.Click("#nav-location", chromedp.ByID),

		chromedp.WaitVisible("#location", chromedp.ByID),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLocationDelete(t *testing.T) {
	// TODO: click first location (location-1)?
	// TODO: click the trash icon
	// TODO: back to empty page
}

func TestAccountDelete(t *testing.T) {
}

func TestTransferListEmpty(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.Click("#nav-transfer", chromedp.ByID),

		chromedp.WaitVisible("#transfer", chromedp.ByID),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestScheduleListEmpty(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.Click("#nav-schedule", chromedp.ByID),

		chromedp.WaitVisible("#schedule", chromedp.ByID),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHistoryListEmpty(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.Click("#nav-history", chromedp.ByID),

		chromedp.WaitVisible("#history", chromedp.ByID),
	)
	if err != nil {
		t.Fatal(err)
	}
}
