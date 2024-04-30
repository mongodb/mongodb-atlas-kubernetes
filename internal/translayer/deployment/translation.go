package deployment

import "go.mongodb.org/atlas-sdk/v20231115008/admin"

type Conn struct {
	Name             string
	ConnURL          string
	SrvConnURL       string
	PrivateURL       string
	SrvPrivateURL    string
	Serverless       bool
	PrivateEndpoints []PrivateEndpoint
}

type PrivateEndpoint struct {
	URL      string
	SrvURL   string
	ShardURL string
}

func clustersToConns(clusters []admin.AdvancedClusterDescription) []Conn {
	conns := []Conn{}
	for _, c := range clusters {
		conns = append(conns, Conn{
			Name:             c.GetName(),
			ConnURL:          c.ConnectionStrings.GetStandard(),
			SrvConnURL:       c.ConnectionStrings.GetStandardSrv(),
			PrivateURL:       c.ConnectionStrings.GetPrivate(),
			SrvPrivateURL:    c.ConnectionStrings.GetPrivateSrv(),
			Serverless:       false,
			PrivateEndpoints: fillClusterPrivateEndpoints(c.ConnectionStrings.GetPrivateEndpoint()),
		})
	}
	return conns
}

func fillClusterPrivateEndpoints(cpeList []admin.ClusterDescriptionConnectionStringsPrivateEndpoint) []PrivateEndpoint {
	pes := []PrivateEndpoint{}
	if len(cpeList) == 0 {
		return pes
	}
	for _, cpe := range cpeList {
		pes = append(pes, PrivateEndpoint{
			URL:      cpe.GetConnectionString(),
			SrvURL:   cpe.GetSrvConnectionString(),
			ShardURL: cpe.GetSrvShardOptimizedConnectionString(),
		})
	}
	return pes
}

func serverlessToConns(serverless []admin.ServerlessInstanceDescription) []Conn {
	conns := []Conn{}
	for _, s := range serverless {
		conns = append(conns, Conn{
			Name:             s.GetName(),
			ConnURL:          "",
			SrvConnURL:       s.ConnectionStrings.GetStandardSrv(),
			Serverless:       true,
			PrivateEndpoints: fillServerlessPrivateEndpoints(s.ConnectionStrings.GetPrivateEndpoint()),
		})
	}
	return conns
}

func fillServerlessPrivateEndpoints(cpeList []admin.ServerlessConnectionStringsPrivateEndpointList) []PrivateEndpoint {
	pes := []PrivateEndpoint{}
	if len(cpeList) == 0 {
		return pes
	}
	for _, cpe := range cpeList {
		pes = append(pes, PrivateEndpoint{
			SrvURL: cpe.GetSrvConnectionString(),
		})
	}
	return pes
}

func connSet(conns ...[]Conn) []Conn {
	return set(func(conn Conn) string { return conn.Name }, conns...)
}

func set[T any](nameFn func(T) string, lists ...[]T) []T {
	hash := map[string]struct{}{}
	result := []T{}
	for _, list := range lists {
		for _, item := range list {
			name := nameFn(item)
			if _, found := hash[name]; !found {
				hash[name] = struct{}{}
				result = append(result, item)
			}
		}
	}
	return result
}
