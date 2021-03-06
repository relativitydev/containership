// Code generated by mockery v2.8.0. DO NOT EDIT.

package processor

import mock "github.com/stretchr/testify/mock"

// MockRegistryClient is an autogenerated mock type for the RegistryClient type
type MockRegistryClient struct {
	mock.Mock
}

// copy provides a mock function with given fields: sourceFQN, destFQN, creds
func (_m *MockRegistryClient) copy(sourceFQN string, destFQN string, creds RegistryCredentials) error {
	ret := _m.Called(sourceFQN, destFQN, creds)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, RegistryCredentials) error); ok {
		r0 = rf(sourceFQN, destFQN, creds)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// delete provides a mock function with given fields: destFQN, creds
func (_m *MockRegistryClient) delete(destFQN string, creds RegistryCredentials) error {
	ret := _m.Called(destFQN, creds)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, RegistryCredentials) error); ok {
		r0 = rf(destFQN, creds)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// listTags provides a mock function with given fields: repository, creds
func (_m *MockRegistryClient) listTags(repository string, creds RegistryCredentials) ([]string, error) {
	ret := _m.Called(repository, creds)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string, RegistryCredentials) []string); ok {
		r0 = rf(repository, creds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, RegistryCredentials) error); ok {
		r1 = rf(repository, creds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
