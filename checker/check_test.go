package checker

import (
	"log"
	"testing"
)

type TestCheckMock struct {
	configured bool
	data map[string]interface{}
}
func NewTestCheckMock() Check {
	return &TestCheckMock{}
}
func (m *TestCheckMock) Configure(options map[string]interface{}) error {
	m.configured = true
	m.data = options
	return nil
}
func (m *TestCheckMock) Run() error {
	return nil
}

func TestCheckRegistry_CreateCheck(t *testing.T) {
	registry := NewCheckRegistry()
	registry.Add("mock", NewTestCheckMock)

	check, err := registry.CreateAndConfigure("mock", map[string]interface{}{"foo": "bar"})
	if err != nil {
		log.Fatalf("Configuration failed %s", err)
	}
	if !check.(*TestCheckMock).configured {
		log.Fatal("Get not configured.")
	}
	if foo, ok := check.(*TestCheckMock).data["foo"]; !ok || foo != "bar" {
		log.Fatalf("Configured key is not valid, expected 'bar' got '%s'", foo)
	}
}

// TODO: Add tests
// * failed Configure()
// * missing provider
