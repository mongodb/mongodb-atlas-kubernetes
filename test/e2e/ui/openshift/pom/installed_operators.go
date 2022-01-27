package pom

import (
	"errors"

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
	messageBox               = "[data-test=\"msg-box-title\"]"
	// test messages
	emptyOperatorList = "empty Operator list"
)

type InstalledOperatorsPage struct {
	P playwright.Page
}

type FilteredInstallPage InstalledOperatorsPage

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

func (iop *InstalledOperatorsPage) SearchByName(smth string) *FilteredInstallPage {
	iop.P.Click(filterButtonLoc)
	iop.P.Click(filterNameTypeLoc)
	iop.P.Type(searchInputLoc, smth)
	iop.P.WaitForSelector(filterClear)
	return &FilteredInstallPage{iop.P}
}

func (iop *InstalledOperatorsPage) SearchByLabel(smth string) *FilteredInstallPage {
	iop.P.Click(filterButtonLoc)
	iop.P.Click(filterLabelTypeLoc)
	iop.P.Type(searchInputLoc, smth)
	iop.P.WaitForSelector(filterClear)
	return &FilteredInstallPage{iop.P}
}

func (fip *FilteredInstallPage) EditAOSubscription() error {
	if fip.isOperatorListExist() {
		Expect(fip.P.Click(actionOperatorMenu)).ShouldNot(HaveOccurred(), "Could not Click on "+actionOperatorMenu)
		Expect(fip.P.Click(editSubscriptionOperator)).ShouldNot(HaveOccurred(), "Could not Click on "+editSubscriptionOperator)
		return nil
	}
	return errors.New(emptyOperatorList)
}

func (fip *FilteredInstallPage) DeleteAOperator() error {
	if fip.isOperatorListExist() {
		Expect(fip.P.Click(actionOperatorMenu)).ShouldNot(HaveOccurred(), "Could not Click on "+actionOperatorMenu)
		Expect(fip.P.Click(deleteOperator)).ShouldNot(HaveOccurred(), "Could not Click on "+deleteOperator)
		Expect(ConfirmAction(fip.P)).ShouldNot(HaveOccurred(), "Could not Click on "+deleteOperator)
		return nil
	}
	return errors.New(emptyOperatorList)
}

func (fip *FilteredInstallPage) isOperatorListExist() bool {
	hasWarningMessage, err := fip.P.IsVisible(messageBox, playwright.FrameIsVisibleOptions{
		Timeout: playwright.Float(timeoutShort),
	})
	Expect(err).ShouldNot(HaveOccurred(), "Unexpected error")
	return !hasWarningMessage
}
