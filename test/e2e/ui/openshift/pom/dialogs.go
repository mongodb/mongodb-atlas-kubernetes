package pom

import "github.com/mxschmitt/playwright-go"

const (
	confirmButton = "#confirm-action"
)

func ConfirmAction(page playwright.Page) error {
	err := page.Click(confirmButton)
	page.WaitForSelector(confirmButton, playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateDetached,
		Timeout: playwright.Float(60000),
	})
	return err
}
