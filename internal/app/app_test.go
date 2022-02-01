package app_test

import (
	"context"
	"flag"
	"fmt"
	"io"
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
	logger := log.NewLogger(io.Discard)

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

func TestAccountRegister(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.WaitVisible("#register", chromedp.ByID),
		chromedp.Click("#register", chromedp.ByID),

		chromedp.WaitVisible("#submit", chromedp.ByID),
		chromedp.SendKeys("#email", email, chromedp.ByID),
		chromedp.SendKeys("#username", username, chromedp.ByID),
		chromedp.SendKeys("#password", password, chromedp.ByID),
		chromedp.Submit("#submit", chromedp.ByID),

		chromedp.WaitVisible("#dashboard", chromedp.ByID),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccountLogout(t *testing.T) {
}

func TestAccountLogin(t *testing.T) {
}

func TestLocationListEmpty(t *testing.T) {
}

func TestLocationCreate(t *testing.T) {
}

func TestLocationList(t *testing.T) {
}

func TestLocationDelete(t *testing.T) {
}

func TestAccountDelete(t *testing.T) {
}
