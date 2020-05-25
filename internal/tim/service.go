package tim

import (
	"context"
	"strconv"

	"github.com/chromedp/chromedp"
)

type Credentials struct {
	Username string
	Password string
}

func noHeadless(a *chromedp.ExecAllocator) {
	chromedp.Flag("headless", false)(a)
}

func GetAvailableDataBytes(credentials Credentials) (int, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:])

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var availableGBStr string

	emailLoginInput := "#caringLoginEmail"
	passwordLoginInput := "#caringLoginPwd"
	loginButton := "#button_caring_login"
	availableDataLabel := "#to-hero-dashboard__counters_mobile .to-hero-dashboard__counter:nth-of-type(2) .tm-traffic-counter__data .tm-traffic-counter__data-title"

	loginURL := "https://mytim.tim.it/login-authorize.html"
	//dashboardURL := "https://mytim.tim.it/it.html"

	if err := chromedp.Run(ctx,
		chromedp.Navigate(loginURL),
		chromedp.WaitVisible(emailLoginInput),
		chromedp.SendKeys(emailLoginInput, credentials.Username),
		chromedp.SendKeys(passwordLoginInput, credentials.Password),
		chromedp.Click(loginButton),
		chromedp.WaitVisible(availableDataLabel),
		chromedp.Text(availableDataLabel, &availableGBStr),
	); err != nil {
		return -1, err
	}

	availableGB, err := strconv.Atoi(availableGBStr)

	if err != nil {
		return -1, err
	}

	return availableGB * 1024 * 1024 * 1024, nil
}
