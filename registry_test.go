package uax

import "testing"

func TestRegistryCoverage(t *testing.T) {
	r := NewRegistry()
	rule := Rule{ID: "test", Brand: "TestBrand"}

	r.RegisterEngine(rule)
	r.RegisterOS(rule)
	r.RegisterBot(rule)
	r.RegisterApp(rule)
	r.RegisterAutomation(rule)
	r.RegisterHTTPClient(rule)

	// Call register with empty list to cover that branch
	r.RegisterEngine()

	r.RegisterBrowser(Rule{ID: "low", Priority: 10}, Rule{ID: "high", Priority: 20})

	r.Reset()
}
