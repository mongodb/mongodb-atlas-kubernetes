package pom

const (
	baseURL                  = "https://console-openshift-console.apps.kubeteam-oshift.2g12.p1.openshiftapps.com"
	loginEndPoint            = "/auth/login"
	dashboardEndPoint        = "/dashboards"
	operatorHubEndPoint      = "/operatorhub"
	installedOeratorEndPoint = "/k8s/all-namespaces/operators.coreos.com~v1alpha1~ClusterServiceVersion"
	serverAPI                = "https://api.kubeteam-oshift.2g12.p1.openshiftapps.com:6443"
	tokenAuthLink            = "https://oauth-openshift.apps.kubeteam-oshift.2g12.p1.openshiftapps.com/oauth/token/request" //#nosec
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
	return tokenAuthLink
}

func ServerAPI() string {
	return serverAPI
}
