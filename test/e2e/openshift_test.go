package e2e_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/mxschmitt/playwright-go"

	actions "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/oc"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/opm"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/podman"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/ui/openshift/pom"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/ui/openshift/pom/pagereport"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	lockNamespace = "lock-test"
)

var _ = Describe("Openshift UI test", Label("openshift"), func() {
	var s utils.Secrets
	var t model.TestDataProvider
	var operatorTag, path string
	var pw *playwright.Playwright
	var browser playwright.Browser
	var page playwright.Page

	BeforeEach(func() {
		testKeys := []string{"DOCKER_REGISTRY", "DOCKER_PASSWORD", "DOCKER_USERNAME", "OPENSHIFT_USER", "OPENSHIFT_PASS", "BUNDLE_IMAGE"}
		s = prepareSecrets(testKeys)

		// check environment
		oc.Version()
		podman.Version()
		opm.Version()

		operatorTag = strings.Split(s["BUNDLE_IMAGE"], ":")[1]
		operatorTag = strings.ToLower(operatorTag)
		Expect(s["BUNDLE_IMAGE"]).ShouldNot(BeEmpty(), "Could not get a image name. Please, set up BUNDLE_IMAGE environment variable")
		Expect(operatorTag).ShouldNot(BeEmpty())

		pw, browser, page = prepareBrowser()
	})

	AfterEach(func() {
		if CurrentSpecReport().Failed() && !strings.Contains(page.URL(), "token") {
			pagereport.MakeScreenshot(page, "error")
		}
		oc.Delete(path) // we delete it all the time, because of shared space
		kubecli.DeleteResource("configmap", lockNamespace, lockNamespace) // clean lockConfig Map
		closeBrowser(pw, browser, page)
	})

	It("User can deploy Atlas Kubernetes operator from openshift", func() {
		By("user resources", func() {
			// TODO need for the next task
			t = model.NewTestDataProvider(
				"operator-in-openshift",
				model.NewEmptyAtlasKeyType().UseDefaulFullAccess(),
				[]string{"data/atlascluster_basic.yaml"},
				[]string{},
				[]model.DBUser{
					*model.NewDBUser("user1").
						WithSecretRef("dbuser-secret-u1").
						AddBuildInAdminRole(),
				},
				30000,
				[]func(*model.TestDataProvider){
					actions.DeleteFirstUser,
				},
			)
		})
		By("Login into openshift", func() {
			pom.NavigateLogin(page).With(s["OPENSHIFT_USER"], s["OPENSHIFT_PASS"]).WaitLoad()
			code := pom.NavigateTokenPage(page).With(s["OPENSHIFT_USER"], s["OPENSHIFT_PASS"]).GetCode()
			Expect(code).ShouldNot(BeEmpty())

			oc.Login(code, pom.ServerAPI())
		})

		By("Lock environment", func() {
			kubecli.CreateNamespace(lockNamespace)
			Eventually(hasLock(), "40m", "10s").Should(BeFalse()) // TODO need to look how it is working and fix timeout
			kubecli.CreateConfigMapWithLiterals(lockNamespace, lockNamespace)
		})

		By("Prepare custom catalog for openshift", func() {
			indexCatalogName := opm.AddIndex(s["BUNDLE_IMAGE"])
			podman.Login(s["DOCKER_REGISTRY"], s["DOCKER_USERNAME"], s["DOCKER_PASSWORD"])
			podman.PushIndexCatalog(indexCatalogName)

			data := utils.JSONToYAMLConvert(model.NewCatalogSource(indexCatalogName))
			path = t.Resources.GetServiceCatalogSourceFolder() + "/catalog-" + operatorTag + ".yaml" // TODO temp. need review/refactor > working with ALL resources paths
			utils.SaveToFile(path, data)

			oc.Apply(path)
		})

		By("delete installed operator, install new one", func() {
			pom.NavigateInstalledOperators(page).SearchByName("Atlas").DeleteAOperator()

			pagereport.MakeScreenshot(page, "deleted")
			pom.NavigateOperatorHub(page).ChooseProviderType(operatorTag).Search("MongoDB Atlas Operator").InstallAtlasOperator()
			pagereport.MakeScreenshot(page, "installed")
		})
		By("final screenshot, clean", func() {
			pom.NavigateInstalledOperators(page).SearchByName("Atlas").DeleteAOperator()
			pagereport.MakeScreenshot(page, "final")
			oc.Delete(path)
		})
	})
})

func prepareBrowser() (*playwright.Playwright, playwright.Browser, playwright.Page) {
	pw, err := playwright.Run()
	Expect(err).ShouldNot(HaveOccurred(), fmt.Sprintf("could not launch playwright: %v", err))
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Args:     []string{"--ignore-certificate-errors"},
	}) // , "--headless"
	Expect(err).ShouldNot(HaveOccurred(), fmt.Sprintf("could not launch Chromium: %v", err))
	page, err := browser.NewPage()
	Expect(err).ShouldNot(HaveOccurred(), fmt.Sprintf("could not open new page: %v", err))
	return pw, browser, page
}

func closeBrowser(pw *playwright.Playwright, browser playwright.Browser, page playwright.Page) {
	Expect(page.Close()).ShouldNot(HaveOccurred(), "could not close page")
	Expect(browser.Close()).ShouldNot(HaveOccurred(), "could not close browser")
	Expect(pw.Stop()).ShouldNot(HaveOccurred(), "could not stop Playwright")
}

func prepareSecrets(testKeys []string) utils.Secrets {
	s := utils.GetSecretEnvOrActrc(testKeys)
	for _, key := range testKeys {
		Expect(s).Should(HaveKeyWithValue(key, Not(BeEmpty())))
	}
	return s
}

func hasLock() func() bool {
	return func() bool {
		layout := "2006-01-02T15:04:05Z"
		if kubecli.HasConfigMap(lockNamespace, lockNamespace) {
			createTime, err := time.Parse(layout, string(kubecli.GetResourceCreationTimestamp("configmap", lockNamespace, lockNamespace)))
			Expect(err).ShouldNot(HaveOccurred())

			if time.Since(createTime).Minutes() > 15 { // TODO usually takes from 9 to 11 minutes
				kubecli.DeleteResource("configmap", lockNamespace, lockNamespace)
				return false
			}
			return true
		}
		return false
	}
}
