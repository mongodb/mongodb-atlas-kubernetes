package pom

import (
	. "github.com/onsi/gomega"

	"github.com/mxschmitt/playwright-go"
)

const (
	displayCodeButton = "button"
	codeField         = "code"
)

type TokenPage struct {
	P playwright.Page
}

func NewTokenPage(page playwright.Page) *TokenPage {
	return &TokenPage{
		page,
	}
}

func NavigateTokenPage(page playwright.Page) *TokenPage {
	_, err := page.Goto(TokenPageLink(), playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	Expect(err).ShouldNot(HaveOccurred(), "Could not navigate to token page")
	return &TokenPage{
		page,
	}
}

func (e *TokenPage) GetCode() string {
	e.P.WaitForSelector(displayCodeButton)
	e.P.Click(displayCodeButton)
	code, err := e.P.InnerText(codeField)
	Expect(err).ShouldNot(HaveOccurred())
	return code
}
