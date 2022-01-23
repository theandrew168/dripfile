package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/chromedp/chromedp"
)

func RandomString(n int) string {
	valid := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

	buf := make([]byte, n)
	for i := range buf {
		buf[i] = valid[rand.Intn(len(valid))]
	}

	return string(buf)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel := context.WithTimeout(taskCtx, 15*time.Second)
	defer cancel()

	email := RandomString(16) + "@example.com"
	fmt.Printf("register: %s\n", email)

	err := chromedp.Run(ctx,
		chromedp.Navigate("http://localhost:5000"),
		chromedp.WaitVisible("#register", chromedp.ByID),
		chromedp.Click("#register", chromedp.ByID),

		chromedp.WaitVisible("#submit", chromedp.ByID),
		chromedp.SendKeys("#email", email, chromedp.ByID),
		chromedp.SendKeys("#username", "foo", chromedp.ByID),
		chromedp.SendKeys("#password", "bar", chromedp.ByID),
		chromedp.Submit("#submit", chromedp.ByID),

		chromedp.WaitVisible("#dashboard", chromedp.ByID),
		chromedp.ActionFunc(func(context.Context) error {
			time.Sleep(time.Second)
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
}
