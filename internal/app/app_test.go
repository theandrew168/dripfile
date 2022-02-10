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

func input(id, msg string) []chromedp.Action {
	var actions []chromedp.Action
	actions = append(actions, chromedp.WaitVisible(id))
	for _, c := range msg {
		actions = append(actions, chromedp.SendKeys(id, string(c)))
		actions = append(actions, chromedp.ActionFunc(func(context.Context) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}))
	}
	return actions
}

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
	username = test.RandomString(8)
	password = username
	email = username + "@dripfile.com"
	fmt.Println(username)

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

	timeoutCtx, cancel := context.WithTimeout(taskCtx, 20*time.Second)
	defer cancel()

	// assign chromedp context
	ctx = timeoutCtx

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
		chromedp.WaitVisible("#page-index"),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccountRegister(t *testing.T) {
	var actions []chromedp.Action
	actions = append(actions, chromedp.WaitVisible("#nav-register"))
	actions = append(actions, chromedp.Click("#nav-register"))
	actions = append(actions, input("#input-email", email)...)
	actions = append(actions, input("#input-username", username)...)
	actions = append(actions, input("#input-password", password)...)
	actions = append(actions, chromedp.WaitVisible("#submit-register"))
	actions = append(actions, chromedp.Submit("#submit-register"))
	actions = append(actions, chromedp.WaitVisible("#page-dashboard"))

	err := chromedp.Run(ctx, actions...)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccountLogout(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.WaitVisible("#nav-account"),
		chromedp.Click("#nav-account"),

		chromedp.WaitVisible("#submit-logout"),
		chromedp.Submit("#submit-logout"),

		chromedp.WaitVisible("#page-index"),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccountLogin(t *testing.T) {
	var actions []chromedp.Action
	actions = append(actions, chromedp.WaitVisible("#nav-login"))
	actions = append(actions, chromedp.Click("#nav-login"))
	actions = append(actions, input("#input-email", email)...)
	actions = append(actions, input("#input-password", password)...)
	actions = append(actions, chromedp.WaitVisible("#submit-login"))
	actions = append(actions, chromedp.Submit("#submit-login"))
	actions = append(actions, chromedp.WaitVisible("#page-dashboard"))

	err := chromedp.Run(ctx, actions...)
	if err != nil {
		t.Fatal(err)
	}
}

// TODO: verify empty
func TestLocationListEmpty(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.WaitVisible("#nav-location"),
		chromedp.Click("#nav-location"),

		chromedp.WaitVisible("#page-location-list"),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLocationCreate(t *testing.T) {
	info := test.RandomString(32)

	var actions []chromedp.Action
	actions = append(actions, chromedp.WaitVisible("#action-create"))
	actions = append(actions, chromedp.Click("#action-create"))
	actions = append(actions, input("#input-endpoint", info)...)
	actions = append(actions, input("#input-access-key-id", info)...)
	actions = append(actions, input("#input-secret-access-key", info)...)
	actions = append(actions, input("#input-bucket-name", info)...)
	actions = append(actions, chromedp.WaitVisible("#submit-create"))
	actions = append(actions, chromedp.Submit("#submit-create"))
	actions = append(actions, chromedp.WaitVisible("#page-location-read"))

	err := chromedp.Run(ctx, actions...)
	if err != nil {
		t.Fatal(err)
	}
}

// TODO: check for details
func TestLocationRead(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.WaitVisible("#page-location-read"),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLocationDelete(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.WaitVisible("#submit-delete"),
		chromedp.Submit("#submit-delete"),

		chromedp.WaitVisible("#page-location-list"),
	)
	if err != nil {
		t.Fatal(err)
	}
}

// TODO: verify empty
func TestTransferListEmpty(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.WaitVisible("#nav-transfer"),
		chromedp.Click("#nav-transfer"),

		chromedp.WaitVisible("#page-transfer-list"),
	)
	if err != nil {
		t.Fatal(err)
	}
}

// TODO: verify empty
func TestScheduleListEmpty(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.WaitVisible("#nav-schedule"),
		chromedp.Click("#nav-schedule"),

		chromedp.WaitVisible("#page-schedule-list"),
	)
	if err != nil {
		t.Fatal(err)
	}
}

// TODO: verify empty
func TestHistoryListEmpty(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.WaitVisible("#nav-history"),
		chromedp.Click("#nav-history"),

		chromedp.WaitVisible("#page-history-list"),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccountRead(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.WaitVisible("#nav-account"),
		chromedp.Click("#nav-account"),

		chromedp.WaitVisible("#page-account-read"),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccountDelete(t *testing.T) {
	err := chromedp.Run(ctx,
		chromedp.WaitVisible("#submit-delete"),
		chromedp.Submit("#submit-delete"),

		chromedp.WaitVisible("#page-index"),
	)
	if err != nil {
		t.Fatal(err)
	}
}
