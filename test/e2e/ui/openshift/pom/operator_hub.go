package pom

import (
	"fmt"

	"github.com/mxschmitt/playwright-go"

	. "github.com/onsi/gomega"
)

const (
	searchLoc         = "[data-test=search-operatorhub]"
	installConfirmLoc = "[data-test-id=\"operator-install-btn\"]"
	installButtonLoc  = "[data-test=\"install-operator\"]"
	viewOperatorLoc   = "text=\"View Operator\""
)

type MarketPage struct {
	P                 playwright.Page
	CatalogSourceName string
}

func NewMarketPage(page playwright.Page) *MarketPage {
	return &MarketPage{
		page,
		"",
	}
}

func NavigateOperatorHub(page playwright.Page) *MarketPage {
	_, err := page.Goto(OperatorHubLink(), playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	Expect(err).ShouldNot(HaveOccurred(), "Could not navigate to Installed Operators page")
	return &MarketPage{
		page,
		"",
	}
}

func (m *MarketPage) Search(name string) *MarketPage {
	m.P.Type(searchLoc, name)
	return m
}

func (m *MarketPage) ChooseProviderType(catalogName string) *MarketPage {
	m.CatalogSourceName = catalogName
	m.P.Check(fmt.Sprintf("[title=\"%s\"]", m.CatalogSourceName), playwright.FrameCheckOptions{
		Timeout: playwright.Float(timeout),
	})
	return m
}

func (m *MarketPage) InstallAtlasOperator() *MarketPage {
	atlasOperatorLoc := fmt.Sprintf("[data-test=\"mongodb-atlas-kubernetes-%s-openshift-marketplace\"]", m.CatalogSourceName)
	err := m.P.Click(atlasOperatorLoc, playwright.PageClickOptions{
		Timeout: playwright.Float(timeout),
	})
	Expect(err).ShouldNot(HaveOccurred(), "Please, make sure the test-catalog is deployed")
	Expect(m.P.Click(installConfirmLoc)).ShouldNot(HaveOccurred())
	Expect(m.P.Click(installButtonLoc)).ShouldNot(HaveOccurred())
	_, err = m.P.WaitForSelector(viewOperatorLoc, playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateAttached,
		Timeout: playwright.Float(timeout),
	})
	Expect(err).ShouldNot(HaveOccurred())
	return m
}
