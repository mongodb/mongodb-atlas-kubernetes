package pom

import (
	"github.com/mxschmitt/playwright-go"

	. "github.com/onsi/gomega"
)

const (
	timeout           = 300000
	searchLoc         = "[data-test=search-operatorhub]"
	atlasOperatorLoc  = "[data-test=\"mongodb-atlas-kubernetes-community-operators-openshift-marketplace\"]"
	installConfirmLoc = "[data-test-id=\"operator-install-btn\"]"
	installButtonLoc  = "[data-test=\"install-operator\"]"
	succesIcon        = "[data-test=\"success-icon\"]"
)

type MarketPage struct {
	P playwright.Page
}

func NewMarketPage(page playwright.Page) *MarketPage {
	return &MarketPage{
		page,
	}
}

func NavigateOperatorHub(page playwright.Page) *MarketPage {
	_, err := page.Goto(OperatorHubLink(), playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	Expect(err).ShouldNot(HaveOccurred(), "Could not navigate to Installed Operators page")
	return &MarketPage{
		page,
	}
}

func (m *MarketPage) Search(name string) *MarketPage {
	m.P.Type(searchLoc, name)
	return m
}

func (m *MarketPage) ChooseProviderType(providerLocator string) *MarketPage {
	m.P.Check(providerLocator)
	return m
}

func (m *MarketPage) InstallAtlasOperator() *MarketPage {
	m.P.Click(atlasOperatorLoc)
	m.P.Click(installConfirmLoc)
	m.P.Click(installButtonLoc)
	t := new(float64)
	*t = timeout
	_, err := m.P.WaitForSelector(succesIcon, playwright.PageWaitForSelectorOptions{
		Timeout: t,
	})
	Expect(err).ShouldNot(HaveOccurred())
	return m
}
