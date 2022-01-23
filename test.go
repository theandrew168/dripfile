package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
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

	err := chromedp.Run(ctx,
		chromedp.Navigate("http://localhost:5000"),
		chromedp.WaitVisible("body > footer"),

		chromedp.Click("#register", chromedp.NodeVisible),
		chromedp.WaitVisible("body > footer"),

		chromedp.SendKeys("#email", "foo@bar.com", chromedp.ByID),
		chromedp.SendKeys("#username", "foo", chromedp.ByID),
		chromedp.SendKeys("#password", "bar", chromedp.ByID),
		chromedp.Click("#submit", chromedp.NodeVisible),
		chromedp.WaitVisible("body > footer"),
	)
	if err != nil {
		log.Fatal(err)
	}
}
