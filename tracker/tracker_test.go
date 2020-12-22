package tracker_test

type stubOS struct {
	applicationName  string
	realTime         bool
	getAppNameCalled int
	nowCalled        int
	shouldLog        int
	logChan          chan string
}