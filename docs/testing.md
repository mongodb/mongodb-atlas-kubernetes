# Testing Guidelines

## High level advice

_Prefer reliable, fast and focused tests over end to end (e2e) testing_.

### Start with a spec

Before starting the implementation, **focus on more detailed behavior the operator code must have**. Having a clear written down specification of the change helps. Use this spec to inform which tests need to be satisfied.

### Exploratory testing as needed

For changes that involve dependencies, such as Atlas, a Cloud provider API or Kubernetes, it might be the case that just checking their API documentation might not be enough.

It is possible that a test is needed to learn how the interaction between the operator and the dependency should play out in the context of the given specification.

But **do not assume this test will become the final e2e test** for the change.

We will refer to this testing phase as _exploratory testing_.

### Prefer TDD

**When the detailed behavior on part of the operator is clear, you can start writing the unit or integration test that covers it**. Unit tests are preferable whenever possible, as they will allow faster and more reliable development feedback loops. 

**Keep your tests focused on a single check**. Note a single _exploratory test_ scenario will usually imply several steps or interactions to be tested individually.

**Try to write the test before**, following the specification and/or any details learnt from the exploratory testing. Test Driven Development (TDD) will force you to think on the spec and the interface of the change or feature design before its implementation. It will probably inform you to use more arguments and simpler functions so you can test the specified behavior more easily in isolation.

**Avoid testing implementation details**, focus on specified outcomes and behaviors. [More on this later](#prefer-external-tests)

Testing against a real Kubernetes, Atlas or a Cloud Provider is not possible on a unit test. But we have support for mocking those dependencies sufficiently so that we can test most of the operator behavior without resorting to an e2e test or integration test. Those support mechanisms are explained in the [Test isolation support section](#test-isolation-support).

### Evaluate additional testing requirements

Once unit tests work as expected, you should be pretty confident your code does what the test spec says. Still you might want to add an integration test if the feature covers a complex workflow or set of steps.

It might even be the case now that the initial exploratory testing code must be converted to an e2e test. Note this should not be done lightly, as e2e tests add costs to the development process and CI pipeline feedback. If unit tests and/or integration tests cover the change or feature properly and the e2e feels redundant, it is preferably to not include it. On the other hand, if the e2e might catch issues that could have been missed otherwise, or need to detect dependency behavior changes, you should include the e2e.

Note the e2e is the "last line of defense", but unit and integration tests and should be the primary gatekeepers against bugs, issues or regressions.

### <a name="prefer-external-tests"></a>Prefer external tests

When writing tests for a feature or change, it is best if those can exercise the code as an external consumer, without any knowledge of or access to our code internal implementation. This makes for tests on behavior rather than implementation details. In Go this usually means prefer `package xxx_test` rather than `package xxx` directly.

Still, it might be difficult to do this on some small changes or when adding test code after the fact. When testing from the inside, try to avoid testing implementation details as much as possible. For example, a new change requires a new behavior which is managed by a new internal function. You could test from the public function that uses the new function behavior, but that might be too expensive if it would require a lot of new mocking just to get to the new behavior. In such a case, it might be simpler and more practical to just test the new internal function directly against its spec.

## <a name="test-isolation-support"></a>Test isolation support

To be able to test our code behavior without involving 3rd party dependencies, the _first trick is to use simple functions or types wherever possible_. Some small changes or features might be writable as a self contained function, or function tree or a type with a method set we can test in isolation. If that is the case, we should just go for that, even if we also need to add or extend later integration or e2e tests to verify this new feature plays well with the rest of the system. TDD usually helps to identify ways to do this.

On the other hand, many changes will involve a tighter integration with some dependency that makes us test it at the edges of the dependency interface. We have some support for that as well.

### Kubernetes

For Kubernetes, the [operator-runtime library has a fake package](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client/fake) that allows you to create a client.Client using [`NewClientBuilder`](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client/fake#NewClientBuilder) and fake many Kubernetes interactions. Just note some limitations:

- You must pre-set the relevant Kubernetes status in the client before each test.
- This fake client _cannot simulate complex Kubernetes status scenarios_. But remember your unit tests should be focused on a single check, not validating a full workflow or step sequence.
- This fake client does not support error injection. You might want to resort to TDD and use simple functions for that instead.

### Atlas

For Atlas, the operator currently uses the [Atlas Go SDK](https://github.com/mongodb/atlas-sdk-go). This [client main struct is composed of a set of interfaces](https://pkg.go.dev/go.mongodb.org/atlas@v0.31.0/mongodbatlas#Client), one per Service behind Atlas. A simple way to mock such a client for unit tests, that will probably only call a couple of API endpoints from Atlas at a time, is to replace those service implementations by the mock:

- The mock is an struct that implements the Service interface by calling methods set in that struct for each of the interface methods defined. [See sample below](#sample-snippets).
- Under `test/atlas` you should find all mocked services already being unit tested using this approach.
- And from them, you can also [find examples of how they are being used in unit tests](#sample-usage-mock).
- If a service is mocked but the method you need is not configurable in the mock struct, it means it has not been used before in a test.
- You can add service method implementations following the same structure as the ones already in place. [Add a new function to be set in the struct and called in the method, and just remember to remove the context first parameter](#sample-projects-mock), ignored in mocks.
- If a service is not mocked yet, create a new entry for it at `test/atlas` following the same conventions as already present service mocks there.
- You can use the [impl tool](https://github.com/josharian/impl) to generate the skeleton implementation. Alternatively, you can try using your IDE for filling the implementation with “unimplemented” method calls so that the service interface is satisfied for the compiler. After that, just fill in the methods needed for the test and ignore the rest.
- Once you have everything in place, you just need to pass implementations to the mock that work as expected for the unit test case at hand.
- Note that, unlike the Kubernetes fake above, these Atlas mocks allow for easy error injection as needed.

Alternatively, you can also mock Atlas at the HTTP Client [http.RoundTripper](https://pkg.go.dev/net/http#RoundTripper) implementation. This is achieved by passing a [custom transport](https://github.com/mongodb/mongodb-atlas-kubernetes/blob/main/pkg/util/httputil/transportclient.go) as a [ClientOpt](https://github.com/mongodb/mongodb-atlas-kubernetes/blob/main/pkg/util/httputil/decoratedclient.go#L5) at the [atlas client creation function](https://github.com/mongodb/mongodb-atlas-kubernetes/blob/main/pkg/controller/atlas/client.go#L18). This is usually not recommended, as the test setup is much more complex in this case compared to mocking the client at its service surface. It requires [creating a round tripper type and implementation per test](#sample-http-mock).

### <a name="sample-snippets"></a>Sample snippets

<a name="sample-projects-mock"></a>Sample projects service mock struct and a sample method implementation:

```go
type ProjectsServiceMock struct {
	GetAllProjectsFn func(*mongodbatlas.ListOptions) (*mongodbatlas.Projects, *mongodbatlas.Response, error)
	DeleteFn         func(string) (*mongodbatlas.Response, error)
	...
}

...

func (ps *ProjectsServiceMock) GetAllProjects(_ context.Context, listOptions *mongodbatlas.ListOptions) (*mongodbatlas.Projects, *mongodbatlas.Response, error) {
	if ps.GetAllProjectsFn == nil {
		panic("GetAllProjects was not set for test")
	}
	return ps.GetAllProjectsFn(listOptions)
}
```

<a name="sample-usage-mock"></a>Sample usage of the mock:

```go
client := &mongodbatlas.Client{
Projects: &atlastest.ProjectsServiceMock{
		GetAllProjectsFn: func(listOptions *mongodbatlas.ListOptions) (*mongodbatlas.Projects, *mongodbatlas.Response, error) {
			return &mongodbatlas.Projects{
				Results:    projectTriplets,
				TotalCount: numberOfProjects,
		}, &mongodbatlas.Response{}, nil
	},
}
```

<a name="sample-http-mock"></a>Sample mocking with HTTP client:

```go
func testAtlasClient(t *testing.T, connection atlas.Connection, rt http.RoundTripper) mongodbatlas.Client {
	t.Helper()
client, err := atlas.Client(fakeDomain, connection, nil, httputil.CustomTransport(rt))
...
}

type deploymentDeletionRoundTripper struct {
	called bool
}

func (rt *deploymentDeletionRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	expectedPath := fmt.Sprintf("/%s/api/atlas/v1.5/groups/%s/clusters/%s", fakeDomain, fakeProjectID, fakeDeployment)
	if req.Method == http.MethodDelete && req.URL.Path == expectedPath {
		rt.called = true
		return reply(req, http.StatusNoContent, ""), nil
	}
	panic(fmt.Sprintf("not implemented for %s path=%q", req.Method, req.URL.Path))
}
```

### Other

We plan on adding more mocks or conventions for Cloud Provider APIs or other dependencies as needed. Feel free to propose such new mocks on demand, particularly if they fit well with the other mocks, and are easy to set up and use.

## Test execution

There are 2 types of code in the repository:

- **Production code** that may be imported by other projects. This includes:
	- All **non** `*_test.go` files at `cmd/manager` and `pkg/`.
- Test code used to verify the correctness of the production code.
	- All `*_test.go` files anywhere and all files under `test/`.

Test code can be further decomposed into:

- **Unit tests**, which can be run by simply doing `go test ./...` and should always succeed without any preparations.
	- Still you should normally use `make unit-test` so that default flags such as race detection and coverage are also included.
	- Includes all `*_test.go` files, tests on `test/int/` and `test/e2e/` folders will be skipped by default.
- **Non-unit tests requiring a setup**, such as **integration** and **e2e** tests, live under the `test/` folder and need to be invoked with a explicit *environment variable*, such as `AKO_INT_TEST=1` for integration tests or `AKO_E2E_TEST=1` for e2e tests. In short: 
	- Run **integration tests** with `make int-test label=...`, using a label to limit the tests to be run.
	- Run **e2e tests** with `make e2e label=...`.
	- Note you will need to load extra environment variables, including credentials, to be able to run most of these tests.
	- Includes files under `test/int/` and `test/e2e/` folders.
- **Helper test code** is code used by unit and non-unit tests code which is not part of the production code. For example, mocks and helpers used to make tests easier to write, more succinct and reliable.
	- Such test code might include its own unit tests that will run as any other unit test via `go test ./...` or `make unit-test`.
	- Note this code requires no build tags, as it is only imported by `*_test.go` code it will not become part of the imported production code.
	- Includes basically any **non** `*_test.go` under `test/` excluding `test/int/` and `test/e2e/` folders, such as `test/atlas/` mocks.
	- For historical reasons, many helpers are today under the `test/e2e` folder. The plan it to move them to `test/helper` or some other folder under `test/`, so that their unit tests, if any,  will always be run by default.
