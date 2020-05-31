package tim

import (
	"context"
	"strconv"
	"time"

	"github.com/chromedp/chromedp"
)

type Credentials struct {
	Username string
	Password string
}

func showBrowser(a *chromedp.ExecAllocator) {
	chromedp.Flag("headless", false)(a)
}

type Service struct {
	ctx context.Context
}

var availableDataLabel = "#to-hero-dashboard__counters_mobile .to-hero-dashboard__counter:nth-of-type(2) .tm-traffic-counter__data .tm-traffic-counter__data-title"

func (s *Service) init(credentials Credentials) error {
	if s.ctx != nil {
		return nil
	}

	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		//showBrowser,
		chromedp.ExecPath("/usr/bin/chromium-browser"),
	)

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, _ := chromedp.NewContext(allocCtx)

	// Start the browser
	if err := chromedp.Run(ctx); err != nil {
		return err
	}
	s.ctx = ctx

	return s.login(credentials)
}

func (s *Service) login(credentials Credentials) error {
	emailLoginInput := "#caringLoginEmail"
	passwordLoginInput := "#caringLoginPwd"
	loginButton := "#button_caring_login"
	loginURL := "https://mytim.tim.it/login-authorize.html"

	if err := chromedp.Run(s.ctx,
		chromedp.Navigate(loginURL),
		chromedp.WaitVisible(emailLoginInput),
		chromedp.SendKeys(emailLoginInput, credentials.Username),
		chromedp.SendKeys(passwordLoginInput, credentials.Password),
		chromedp.Click(loginButton),
		chromedp.WaitVisible(availableDataLabel),
	); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetAvailableDataBytes(credentials Credentials) (float64, error) {
	err := s.init(credentials)

	if err != nil {
		return -1, err
	}

	var availableGBStr string
	dashboardURL := "https://mytim.tim.it/it.html"

	ctx, cancel := context.WithTimeout(s.ctx, 40*time.Second)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate(dashboardURL),
		chromedp.WaitVisible(availableDataLabel),
		chromedp.Text(availableDataLabel, &availableGBStr),
	); err != nil {
		return -1, err
	}

	availableGB, err := strconv.ParseFloat(availableGBStr, 64)

	if err != nil {
		return -1, err
	}

	return availableGB * 1024 * 1024 * 1024, nil
}
