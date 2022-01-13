package pom

import (
	"github.com/mxschmitt/playwright-go"
	. "github.com/onsi/gomega"
)

const (
	preLoginButton = "text=\"Cluster-Admin\""
	preLoginButton2 = "[title=\"Log in with Cluster-Admin\"]"
	//  class="pf-c-button pf-m-secondary pf-m-block" title="Log in with Cluster-Admin">Cluster-Admin</a>
	userName     = "#inputUsername"
	userPassword = "#inputPassword"
	loginLocator = "text=\"Log in\""
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
	Expect(page.Click(preLoginButton)).ShouldNot(HaveOccurred(), "Could not find 'Log in with...'")
	_, err = page.WaitForSelector(loginLocator, playwright.PageWaitForSelectorOptions{
		// State:   &"",
		// Strict:  new(bool),
		// Timeout: new(float64),
	})
	Expect(err).ShouldNot(HaveOccurred(), "Wait Could not find Login Locator")
	return &LoginPage{
		page,
	}
}

func (lp *LoginPage) With(user, password string) playwright.Page {
	Expect(lp.P.Type(userName, user)).ShouldNot(HaveOccurred(), "Could not input user name")
	Expect(lp.P.Type(userPassword, password)).ShouldNot(HaveOccurred(), "Could not input password")
	Expect(lp.P.Click(loginLocator)).ShouldNot(HaveOccurred(), "Could not LogIn")
	_, err := lp.P.WaitForNavigation(playwright.PageWaitForNavigationOptions{
		// Timeout:   new(float64),
		URL:       DashboardLink(),
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	Expect(err).ShouldNot(HaveOccurred(), "Wait dashboard page: Could not Login")
	return lp.P
}
