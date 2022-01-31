package pagereport

import (
	"fmt"

	"github.com/mxschmitt/playwright-go"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

func MakeScreenshot(page playwright.Page, name string) error {
	path := fmt.Sprintf("output/openshift/%s.png", name)
	utils.SaveToFile(path, []byte{})
	_, err := page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(path),
	})
	if err != nil {
		return err
	}
	return nil
}
