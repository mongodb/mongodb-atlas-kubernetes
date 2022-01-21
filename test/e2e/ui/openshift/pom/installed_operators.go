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
	filterClear        = "text=\"Clear all filters\""
	// operator table group
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
	iop.P.WaitForSelector(filterClear)
	return iop
}

func (iop *InstalledOperatorsPage) SearchByLabel(smth string) *InstalledOperatorsPage {
	iop.P.Click(filterButtonLoc)
	iop.P.Click(filterLabelTypeLoc)
	iop.P.Type(searchInputLoc, smth)
	iop.P.WaitForSelector(filterClear)
	return iop
}

func (iop *InstalledOperatorsPage) EditAOSubscription() {
	// TODO doesnt work w/o search
	Expect(iop.P.Click(actionOperatorMenu)).ShouldNot(HaveOccurred(), "Could not Click on "+actionOperatorMenu)
	Expect(iop.P.Click(editSubscriptionOperator)).ShouldNot(HaveOccurred(), "Could not Click on "+editSubscriptionOperator)
}

func (iop *InstalledOperatorsPage) DeleteAOperator() {
	// TODO doesnt work w/o search
	Expect(iop.P.Click(actionOperatorMenu)).ShouldNot(HaveOccurred(), "Could not Click on "+actionOperatorMenu)
	Expect(iop.P.Click(deleteOperator)).ShouldNot(HaveOccurred(), "Could not Click on "+deleteOperator)
	Expect(ConfirmAction(iop.P)).ShouldNot(HaveOccurred(), "Could not Click on "+deleteOperator)
}
