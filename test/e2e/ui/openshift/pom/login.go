package pom

import (
	"github.com/mxschmitt/playwright-go"
	. "github.com/onsi/gomega"
)

const (
	preLoginButton = "text=\"Deployment-Admin\""
	userName       = "#inputUsername"
	userPassword   = "#inputPassword"
	loginLocator   = "text=\"Log in\""
)

type LoginPage struct {
	P playwright.Page
}

func NewLogin(page playwright.Page) *LoginPage {
	return &LoginPage{
		page,
	}
}

func NavigateLogin(page playwright.Page) *LoginPage {
	_, err := page.Goto(LoginPageLink(), playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	Expect(err).ShouldNot(HaveOccurred(), "Could not navigate to login page")
	return &LoginPage{
		page,
	}
}

func (lp *LoginPage) With(user, password string) *LoginPage {
	Expect(lp.P.Click(preLoginButton)).ShouldNot(HaveOccurred(), "Could not find 'Log in with...'")
	_, err := lp.P.WaitForSelector(loginLocator, playwright.PageWaitForSelectorOptions{})
	Expect(err).ShouldNot(HaveOccurred(), "Wait Could not find Login Locator")

	Expect(lp.P.Type(userName, user)).ShouldNot(HaveOccurred(), "Could not input user name")
	Expect(lp.P.Type(userPassword, password)).ShouldNot(HaveOccurred(), "Could not input password")
	Expect(lp.P.Click(loginLocator)).ShouldNot(HaveOccurred(), "Could not LogIn")
	return lp
}

func (lp *LoginPage) WaitLoad() *LoginPage {
	_, err := lp.P.WaitForNavigation(playwright.PageWaitForNavigationOptions{
		URL: DashboardLink(),
	})
	Expect(err).ShouldNot(HaveOccurred(), "Wait dashboard page: Could not Login")
	return lp
}
