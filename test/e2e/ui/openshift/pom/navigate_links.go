package pom

import "os"

/*
	Please, add to GitHub Secrets or to environment variables:
	OPENSHIFT_CONSOLE_URL cluster base URL
	OPENSHIFT_SERVER_API service API link
	OPENSHIFT_TOKEN_LINK page where token could be found

	Samples:
	export OPENSHIFT_CONSOLE_URL="https://some.openshiftapps.com"
	export OPENSHIFT_SERVER_API="https://api.some.openshiftapps.com:6443"
	export OPENSHIFT_TOKEN_LINK="https://oauth-some.openshiftapps.com/oauth/token/request"
*/

const (
	loginEndPoint            = "/auth/login"
	dashboardEndPoint        = "/dashboards"
	operatorHubEndPoint      = "/operatorhub"
	installedOeratorEndPoint = "/k8s/all-namespaces/operators.coreos.com~v1alpha1~ClusterServiceVersion"
)

func BaseURL() string {
	return os.Getenv("OPENSHIFT_CONSOLE_URL")
}

func LoginPageLink() string {
	return BaseURL() + loginEndPoint
}

func DashboardLink() string {
	return BaseURL() + dashboardEndPoint
}

func OperatorHubLink() string {
	return BaseURL() + operatorHubEndPoint
}

func InstalledOperatorLink() string {
	return BaseURL() + installedOeratorEndPoint
}

func TokenPageLink() string {
	return os.Getenv("OPENSHIFT_TOKEN_LINK")
}

func ServerAPI() string {
	return os.Getenv("OPENSHIFT_SERVER_API")
}
