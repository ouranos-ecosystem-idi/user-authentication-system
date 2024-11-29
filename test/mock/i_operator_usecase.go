// Code generated by mockery v2.42.3. DO NOT EDIT.

package mocks

import (
	traceability "authenticator-backend/domain/model/traceability"

	mock "github.com/stretchr/testify/mock"
)

// IOperatorUsecase is an autogenerated mock type for the IOperatorUsecase type
type IOperatorUsecase struct {
	mock.Mock
}

// GetOperator provides a mock function with given fields: getOperatorInput
func (_m *IOperatorUsecase) GetOperator(getOperatorInput traceability.GetOperatorInput) (traceability.OperatorModel, error) {
	ret := _m.Called(getOperatorInput)

	if len(ret) == 0 {
		panic("no return value specified for GetOperator")
	}

	var r0 traceability.OperatorModel
	var r1 error
	if rf, ok := ret.Get(0).(func(traceability.GetOperatorInput) (traceability.OperatorModel, error)); ok {
		return rf(getOperatorInput)
	}
	if rf, ok := ret.Get(0).(func(traceability.GetOperatorInput) traceability.OperatorModel); ok {
		r0 = rf(getOperatorInput)
	} else {
		r0 = ret.Get(0).(traceability.OperatorModel)
	}

	if rf, ok := ret.Get(1).(func(traceability.GetOperatorInput) error); ok {
		r1 = rf(getOperatorInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOperators provides a mock function with given fields: getOperatorsInput
func (_m *IOperatorUsecase) GetOperators(getOperatorsInput traceability.GetOperatorsInput) ([]traceability.OperatorModel, error) {
	ret := _m.Called(getOperatorsInput)

	if len(ret) == 0 {
		panic("no return value specified for GetOperators")
	}

	var r0 []traceability.OperatorModel
	var r1 error
	if rf, ok := ret.Get(0).(func(traceability.GetOperatorsInput) ([]traceability.OperatorModel, error)); ok {
		return rf(getOperatorsInput)
	}
	if rf, ok := ret.Get(0).(func(traceability.GetOperatorsInput) []traceability.OperatorModel); ok {
		r0 = rf(getOperatorsInput)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]traceability.OperatorModel)
		}
	}

	if rf, ok := ret.Get(1).(func(traceability.GetOperatorsInput) error); ok {
		r1 = rf(getOperatorsInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutOperator provides a mock function with given fields: operatorModel
func (_m *IOperatorUsecase) PutOperator(operatorModel traceability.OperatorModel) (traceability.OperatorModel, error) {
	ret := _m.Called(operatorModel)

	if len(ret) == 0 {
		panic("no return value specified for PutOperator")
	}

	var r0 traceability.OperatorModel
	var r1 error
	if rf, ok := ret.Get(0).(func(traceability.OperatorModel) (traceability.OperatorModel, error)); ok {
		return rf(operatorModel)
	}
	if rf, ok := ret.Get(0).(func(traceability.OperatorModel) traceability.OperatorModel); ok {
		r0 = rf(operatorModel)
	} else {
		r0 = ret.Get(0).(traceability.OperatorModel)
	}

	if rf, ok := ret.Get(1).(func(traceability.OperatorModel) error); ok {
		r1 = rf(operatorModel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIOperatorUsecase creates a new instance of IOperatorUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIOperatorUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IOperatorUsecase {
	mock := &IOperatorUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
