package app_test

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"math/rand"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/chromedp/chromedp"

	"github.com/theandrew168/dripfile/internal/app"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/log"
	"github.com/theandrew168/dripfile/internal/payment"
	"github.com/theandrew168/dripfile/internal/postgres"
	"github.com/theandrew168/dripfile/internal/random"
	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/task"
	"github.com/theandrew168/dripfile/internal/test"
)

// chromedp docs:
// https://pkg.go.dev/github.com/chromedp/chromedp#Query

// minio container info
const (
	s3Endpoint        = "localhost:9000"
	s3BucketNameFoo   = "foo"
	s3BucketNameBar   = "bar"
	s3AccessKeyID     = "AKIAIOSFODNN7EXAMPLE"
	s3SecretAccessKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
)

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
	username = random.String(8)
	password = username
	email = username + "@dripfile.com"
	fmt.Println(username)

	// TODO: setup test bucket in S3

	logger := log.NewLogger(os.Stdout)
	cfg := test.Config()

	secretKeyBytes, err := hex.DecodeString(cfg.SecretKey)
	if err != nil {
		logger.Error(err)
		return 1
	}

	// create secret.Box
	var secretKey [32]byte
	copy(secretKey[:], secretKeyBytes)
	box := secret.NewBox(secretKey)

	// open a database connection pool
	pool, err := postgres.ConnectPool(cfg.DatabaseURI)
	if err != nil {
		logger.Error(err)
		return 1
	}
	defer pool.Close()

	storage := database.NewPostgresStorage(pool)
	queue := task.NewPostgresQueue(pool)

	// init the billing interface
	var billing payment.Billing
	if cfg.StripePublicKey != "" && cfg.StripeSecretKey != "" {
		billing = payment.NewStripeBilling(cfg.StripePublicKey, cfg.StripeSecretKey)
	} else {
		billing = payment.NewLogBilling(logger)
	}

	// create the application
	handler := app.New(cfg, box, storage, queue, billing, logger)

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
		logger.Error(err)
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
	actions = append(actions, chromedp.WaitVisible("#page-register"))

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

func TestAccountLoginInvalidEmail(t *testing.T) {
	// random email that won't be valid
	email := random.String(8) + "@dripfile.com"

	var actions []chromedp.Action
	actions = append(actions, chromedp.WaitVisible("#nav-login"))
	actions = append(actions, chromedp.Click("#nav-login"))
	actions = append(actions, chromedp.WaitVisible("#page-login"))

	actions = append(actions, input("#input-email", email)...)
	actions = append(actions, input("#input-password", password)...)
	actions = append(actions, chromedp.WaitVisible("#submit-login"))
	actions = append(actions, chromedp.Submit("#submit-login"))

	actions = append(actions, chromedp.WaitVisible("#error-email"))

	err := chromedp.Run(ctx, actions...)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccountLoginInvalidPassword(t *testing.T) {
	// random password that won't be valid
	password := random.String(8)

	var actions []chromedp.Action
	actions = append(actions, chromedp.WaitVisible("#nav-login"))
	actions = append(actions, chromedp.Click("#nav-login"))
	actions = append(actions, chromedp.WaitVisible("#page-login"))

	actions = append(actions, input("#input-email", email)...)
	actions = append(actions, input("#input-password", password)...)
	actions = append(actions, chromedp.WaitVisible("#submit-login"))
	actions = append(actions, chromedp.Submit("#submit-login"))

	actions = append(actions, chromedp.WaitVisible("#error-password"))

	err := chromedp.Run(ctx, actions...)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccountLogin(t *testing.T) {
	var actions []chromedp.Action
	actions = append(actions, chromedp.WaitVisible("#nav-login"))
	actions = append(actions, chromedp.Click("#nav-login"))
	actions = append(actions, chromedp.WaitVisible("#page-login"))

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
	var actions []chromedp.Action
	actions = append(actions, chromedp.WaitVisible("#action-create"))
	actions = append(actions, chromedp.Click("#action-create"))
	actions = append(actions, chromedp.WaitVisible("#page-location-create"))

	actions = append(actions, input("#input-endpoint", s3Endpoint)...)
	actions = append(actions, input("#input-bucket-name", s3BucketNameFoo)...)
	actions = append(actions, input("#input-access-key-id", s3AccessKeyID)...)
	actions = append(actions, input("#input-secret-access-key", s3SecretAccessKey)...)
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
