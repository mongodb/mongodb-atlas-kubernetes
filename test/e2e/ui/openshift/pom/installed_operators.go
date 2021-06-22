package pom

import (
	"github.com/mxschmitt/playwright-go"
	. "github.com/onsi/gomega"
)

const (
	searchInputLoc     = "[data-test-id=item-filter]"
	filterButtonLoc    = "[data-test-id=dropdown-button]"
	filterNameTypeLoc  = "#NAME-link"
	filterLabelTypeLoc = "#LABEL-link"
	// operator table group
	// operatorsList            = "ReactVirtualized__VirtualGrid__innerScrollContainer"
	actionOperatorMenu       = "[data-test-id=kebab-button]"
	editSubscriptionOperator = "[data-test-action=\"Edit Subscription\"]"
	deleteOperator           = "[data-test-action=\"Uninstall Operator\"]"
)

type InstalledOperatorsPage struct {
	P playwright.Page
}

func NewInstalledOperators(page playwright.Page) *InstalledOperatorsPage {
	return &InstalledOperatorsPage{
		page,
	}
}

func NavigateInstalledOperators(page playwright.Page) *InstalledOperatorsPage {
	_, err := page.Goto(InstalledOperatorLink(), playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	Expect(err).ShouldNot(HaveOccurred(), "Could not navigate to Installed Operators page")
	return &InstalledOperatorsPage{
		page,
	}
}

func (iop *InstalledOperatorsPage) SearchByName(smth string) *InstalledOperatorsPage {
	iop.P.Click(filterButtonLoc)
	iop.P.Click(filterNameTypeLoc)
	iop.P.Type(searchInputLoc, smth)
	return iop
}

func (iop *InstalledOperatorsPage) SearchByLabel(smth string) *InstalledOperatorsPage {
	iop.P.Click(filterButtonLoc)
	iop.P.Click(filterLabelTypeLoc)
	iop.P.Type(searchInputLoc, smth)
	return iop
}

func (iop *InstalledOperatorsPage) EditAOSubscription() {
	// TODO doesnt work w/o search
	iop.P.Click(actionOperatorMenu)
	iop.P.Click(editSubscriptionOperator)
}

func (iop *InstalledOperatorsPage) DeleteAOperator() {
	// TODO doesnt work w/o search
	iop.P.Click(actionOperatorMenu)
	iop.P.Click(deleteOperator)
	ConfirmAction(iop.P)
}
