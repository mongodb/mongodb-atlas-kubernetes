package int

import (
	"fmt"
	"net/http"
	"path"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/fakerest"
)

func NewFakeDeploymentAndProjectInAtlasServer() *fakerest.Server {
	return fakerest.NewCombinedServer(
		newExistingProjectServer().Server,
		newExistingDeploymentServer("AWS").Server,
	)
}

type existingProjectServer struct {
	*fakerest.Server
	projectDeleted bool
}

func newExistingProjectServer() *existingProjectServer {
	eps := existingProjectServer{projectDeleted: false}
	eps.Server = fakerest.NewServer(fakerest.Script{
		fakerest.ReplyToGet("/api/atlas/v1.0/groups/byName/",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				return fakerest.OKJSON(rsp, fakeProjectJSON(path.Base(req.URL.Path), 1))
			}),
		fakerest.ReplyToGet("/api/atlas/v1.0/groups/64a6c94435fc02228a1d6aa1/userSecurity",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				return fakerest.OKJSON(rsp, fakeUserSecurityJSON())
			}),
		fakerest.ReplyToGet("/api/atlas/v1.0/groups/64a6c94435fc02228a1d6aa1/accessList/0.0.0.0",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				return fakerest.OKJSON(rsp, `{"STATUS":"ACTIVE"}`)
			}),
		fakerest.ReplyToGet("/api/atlas/v1.0/groups/64a6c94435fc02228a1d6aa1/accessList",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				return fakerest.OKJSON(rsp, fakeAccessListJSON())
			}),
		fakerest.ReplyToPost("/api/atlas/v1.0/groups/64a6c94435fc02228a1d6aa1/accessList",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				return fakerest.StatusJSON(rsp, http.StatusCreated, fakeAccessListJSON())
			}),
		fakerest.ReplyToGet("/api/atlas/v1.0/groups/64a6c94435fc02228a1d6aa1/privateEndpoint/",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				return fakerest.EmptyJSONArray(rsp)
			}),
		fakerest.ReplyToGet("/api/atlas/v1.0/groups/64a6c94435fc02228a1d6aa1/customDBRoles/roles",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				return fakerest.EmptyJSONArray(rsp)
			}),
		fakerest.ReplyToGet("/api/atlas/v1.0/groups/64a6c94435fc02228a1d6aa1/customDBRoles/roles",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				return fakerest.EmptyJSONArray(rsp)
			}),
		fakerest.ReplyToGet("/api/atlas/v1.0/groups/64a6c94435fc02228a1d6aa1/peers",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				return fakerest.OKJSON(rsp, fakeAtlasEmptyListJSON())
			}),
		fakerest.ReplyToGet("/api/atlas/v1.0/groups/64a6c94435fc02228a1d6aa1/containers",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				return fakerest.OKJSON(rsp, fakeAtlasEmptyListJSON())
			}),
		fakerest.ReplyToGet("/api/atlas/v1.0/groups/64a6c94435fc02228a1d6aa1/integrations",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				return fakerest.OKJSON(rsp, fakeAtlasEmptyListJSON())
			}),
		fakerest.ReplyToGet("/api/atlas/v1.0/groups/64a6c94435fc02228a1d6aa1/teams",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				return fakerest.OKJSON(rsp, fakeAtlasEmptyListJSON())
			}),
		fakerest.ReplyToGet("/api/atlas/v1.0/groups/64a6c94435fc02228a1d6aa1/encryptionAtRest",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				return fakerest.OKJSON(rsp, fakeEncryptionStatusJSON())
			}),
		fakerest.ReplyToGet("/api/atlas/v1.0/groups/64a6c94435fc02228a1d6aa1/auditLog",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				return fakerest.OKJSON(rsp, fakeAuditLogStatusJSON())
			}),
		fakerest.ReplyToPatch("/api/atlas/v1.0/groups/64a6c94435fc02228a1d6aa1/auditLog",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				return fakerest.StatusJSON(rsp, http.StatusAccepted, fakeAuditLogStatusJSON())
			}),
		fakerest.ReplyToGet("/api/atlas/v1.0/groups/64a6c94435fc02228a1d6aa1/settings",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				return fakerest.StatusJSON(rsp, http.StatusAccepted, fakeSettingsJSON())
			}),
		fakerest.ReplyToGet("/api/atlas/v1.0/groups/64a6c94435fc02228a1d6aa1",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				if eps.projectDeleted {
					return fakerest.StatusEmpty(rsp, http.StatusNotFound)
				}
				return fakerest.OKJSON(rsp, fakeProjectJSON(path.Base(req.URL.Path), 1))
			}),
		fakerest.ReplyToDelete("/api/atlas/v1.0/groups/64a6c94435fc02228a1d6aa1",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				eps.projectDeleted = true
				return fakerest.RemovedJSONReply(req, rsp)
			}),
	})
	return &eps
}

type existingDeploymentServer struct {
	*fakerest.Server
	providerName      string
	deploymentDeleted bool
}

func newExistingDeploymentServer(providerName string) *existingDeploymentServer {
	eds := existingDeploymentServer{deploymentDeleted: false, providerName: providerName}
	handleGetDeployment := func(req *http.Request, rsp *http.Response) (*http.Response, error) {
		if eds.providerName == "" {
			eds.providerName = "AWS"
		}
		if eds.deploymentDeleted {
			return fakerest.StatusEmpty(rsp, http.StatusNotFound)
		}
		return fakerest.OKJSON(rsp, fakeDeploymentJSON(eds.providerName))
	}
	eds.Server = fakerest.NewServer(fakerest.Script{
		fakerest.ReplyToGet("/api/atlas/v1.5/groups/64a6c94435fc02228a1d6aa1/clusters/test-deployment-aws/globalWrites/",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				if eds.deploymentDeleted {
					return fakerest.StatusEmpty(rsp, http.StatusNotFound)
				}
				return fakerest.OKJSON(rsp, fakeGlobalWritesJSON())
			}),
		fakerest.ReplyToGet(
			"/api/atlas/v1.5/groups/64a6c94435fc02228a1d6aa1/clusters/test-deployment-aws",
			handleGetDeployment,
		),
		fakerest.ReplyToGet(
			"/api/atlas/v1.5/groups/test-project/clusters/test-deployment-aws",
			handleGetDeployment,
		),
		// If you setup the providername as AWS the operator expects TENANT, and when you do use TENAT it expects AWS
		// so we need to handle a patch call
		fakerest.ReplyToPatch("/api/atlas/v1.5/groups/64a6c94435fc02228a1d6aa1/clusters/test-deployment-aws",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				eds.providerName = "TENANT"
				return fakerest.OKJSON(rsp, fakeDeploymentJSON(eds.providerName))
			}),
		fakerest.ReplyToDelete("/api/atlas/v1.5/groups/64a6c94435fc02228a1d6aa1/clusters/test-deployment-aws",
			func(req *http.Request, rsp *http.Response) (*http.Response, error) {
				eds.deploymentDeleted = true
				return fakerest.RemovedJSONReply(req, rsp)
			}),
	})
	return &eds
}

var baseProject string = `{"clusterCount":%d,"created":"2023-07-06T14:01:40Z","id":"64a6c94435fc02228a1d6aa1",` +
	`"links":[],"name":"%s","orgId":"649ae72e2236fa671134962d"}`

func fakeProjectJSON(name string, clusters int) string {
	return fmt.Sprintf(baseProject, clusters, name)
}

func fakeUserSecurityJSON() string {
	return `{"customerX509":{},"ldap":{"authenticationEnabled":false,"authorizationEnabled":false,"userToDNMapping":null}}`
}

func fakeAccessListJSON() string {
	return `{"links":[],"results":[{"cidrBlock":"0.0.0.0/0","groupId":"64a6c94435fc02228a1d6aa1","links":[]}],"totalCount":1}`
}

func fakeAtlasEmptyListJSON() string {
	return `{"links":[],"results":[],"totalCount":0}`
}

func fakeEncryptionStatusJSON() string {
	return `{"awsKms":{"customerMasterKeyID":null,"enabled":false,"region":null,"valid":false},` +
		`"azureKeyVault":{"azureEnvironment":null,"clientID":null,"enabled":false,"keyIdentifier":null,` +
		`"keyVaultName":null,"resourceGroupName":null,"subscriptionID":null,"tenantID":null,"valid":false},` +
		`"googleCloudKms":{"enabled":false,"keyVersionResourceID":null,"valid":false}}`
}

func fakeAuditLogStatusJSON() string {
	return `{"auditAuthorizationSuccess":false,"configurationType":"NONE","enabled":false}`
}

func fakeSettingsJSON() string {
	return `{"isCollectDatabaseSpecificsStatisticsEnabled":true,"isDataExplorerEnabled":true,` +
		`"isExtendedStorageSizesEnabled":false,"isPerformanceAdvisorEnabled":true,` +
		`"isRealtimePerformancePanelEnabled":true,"isSchemaAdvisorEnabled":true}`
}

var baseCluster = `{
		"autoScaling": {"compute": {"enabled": false,"scaleDownEnabled": false},"diskGBEnabled": false},
		"backupEnabled": false,
		"biConnector": {"enabled": false,"readPreference": "PRIMARY"},
		"clusterType": "REPLICASET",
		"connectionStrings": {},
		"createDate": "2019-08-24T14:15:22Z",
		"diskSizeGB": 10,
		"encryptionAtRestProvider": "NONE",
		"groupId": "64a6c94435fc02228a1d6aa1",
		"id": "700000000000000000000000",
		"labels": [],
		"links": [],
		"mongoDBMajorVersion": "4.2",
		"mongoDBVersion": "string",
		"mongoURI": "string",
		"mongoURIUpdated": "2019-08-24T14:15:22Z",
		"mongoURIWithOptions": "string",
		"name": "%s",
		"numShards": 1,
		"paused": false,
		"pitEnabled": false,
		"providerBackupEnabled": false,
		"providerSettings": {
			"instanceSizeName": "M2",
            "providerName": "%s",
            "regionName": "US_EAST_1"
		},
		"replicationFactor": 3,
		"replicationSpec": {},
		"replicationSpecs": [
			{
			  "numShards": 1,
			  "regionConfigs": [
				{
				  "analyticsSpecs": {"instanceSize": "M2","nodeCount": 0},
				  "electableSpecs": {"instanceSize": "M2","nodeCount": 3},
				  "readOnlySpecs": {"instanceSize": "M2","nodeCount": 0},
				  "priority": 7,
				  "providerName": "%s",
				  "regionName": "US_EAST_1"
				}
			  ]
			}
		  ],
		"rootCertType": "ISRGROOTX1",
		"srvAddress": "string",
		"stateName": "IDLE",
		"tags": [],
		"terminationProtectionEnabled": false,
		"versionReleaseSystem": "LTS"
	  }`

func fakeDeploymentJSON(providerName string) string {
	return fmt.Sprintf(baseCluster, "test-deployment-aws", providerName, providerName)
}

func fakeGlobalWritesJSON() string {
	return `{"customZoneMapping": {},"managedNamespaces": []}`
}
