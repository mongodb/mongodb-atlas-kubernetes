package pom

import "github.com/mxschmitt/playwright-go"

const (
	confirmButton = "#confirm-action"
)

func ConfirmAction(page playwright.Page) {
	page.Click(confirmButton)
}
