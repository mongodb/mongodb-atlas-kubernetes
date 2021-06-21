package pom

const (
	baseURL                  = "https://console-openshift-console.apps.openshift.mongokubernetes.com"
	loginEndPoint            = "/auth/login"
	dashboardEndPoint        = "/dashboards"
	operatorHubEndPoint      = "/operatorhub"
	installedOeratorEndPoint = "/k8s/all-namespaces/operators.coreos.com~v1alpha1~ClusterServiceVersion"
)

func LoginPageLink() string {
	return baseURL + loginEndPoint
}

func DashboardLink() string {
	return baseURL + dashboardEndPoint
}

func OperatorHubLink() string {
	return baseURL + operatorHubEndPoint
}

func InstalledOperatorLink() string {
	return baseURL + installedOeratorEndPoint
}

func TokenPageLink() string {
	return "https://oauth-openshift.apps.openshift.mongokubernetes.com/oauth/token/request"
}
