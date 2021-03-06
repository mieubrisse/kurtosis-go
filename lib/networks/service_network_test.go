/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package networks

import (
	"github.com/kurtosis-tech/kurtosis-go/lib/services"
	"os"
	"testing"
	"time"
)

const (
	testServiceName = "test-service"
	testNetworkName = "test-network"
	testConfiguration = "test-configuration"
)

type TestService struct {}

// ======================== Test Initializer Core ========================
type TestInitializerCore struct {}
func (t TestInitializerCore) GetUsedPorts() map[int]bool {
	return make(map[int]bool)
}

func (t TestInitializerCore) GetServiceFromIp(ipAddr string) services.Service {
	return TestService{}
}


func (t TestInitializerCore) GetFilesToMount() map[string]bool {
	return make(map[string]bool)
}

func (t TestInitializerCore) InitializeMountedFiles(filepathsToMount map[string]*os.File, dependencies []services.Service) error {
	return nil
}

func (t TestInitializerCore) GetStartCommand(mountedFileFilepaths map[string]string, ipPlaceholder string, dependencies []services.Service) ([]string, error) {
	return make([]string, 0), nil
}

func (t TestInitializerCore) GetTestVolumeMountpoint() string {
	return "/foo/bar"
}

func getTestInitializerCore() services.ServiceInitializerCore {
	return TestInitializerCore{}
}


// ======================== Test Availability Checker Core ========================
type TestAvailabilityCheckerCore struct {}
func (t TestAvailabilityCheckerCore) IsServiceUp(toCheck services.Service, dependencies []services.Service) bool {
	return true
}
func (t TestAvailabilityCheckerCore) GetTimeout() time.Duration {
	return 30 * time.Second
}
func getTestCheckerCore() services.ServiceAvailabilityCheckerCore {
	return TestAvailabilityCheckerCore{}
}

// ======================== Tests ========================
func TestDisallowingNonexistentConfigs(t *testing.T) {
	builder := NewServiceNetworkBuilder(nil, "/foo/bar")
	network := builder.Build()
	_, err := network.AddService(testConfiguration, testServiceName, make(map[ServiceID]bool))
	if err == nil {
		t.Fatal("Expected error when declaring a service with a configuration that doesn't exist")
	}
}

func TestDisallowingNonexistentDependencies(t *testing.T) {
	var configId ConfigurationID = testConfiguration
	builder := NewServiceNetworkBuilder(nil, "/foo/bar")
	err := builder.AddConfiguration(configId, "test", getTestInitializerCore(), getTestCheckerCore())
	if err != nil {
		t.Fatal("Adding a configuration shouldn't fail")
	}
	network := builder.Build()

	dependencies := map[ServiceID]bool{
		testServiceName: true,
	}

	_, err = network.AddService(configId, testServiceName, dependencies)
	if err == nil {
		t.Fatal("Expected error when declaring a dependency on a service ID that doesn't exist")
	}
}
