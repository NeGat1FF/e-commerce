// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/NeGat1FF/e-commerce/product-service/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// ProductRepository is an autogenerated mock type for the ProductRepository type
type ProductRepository struct {
	mock.Mock
}

// AddStock provides a mock function with given fields: ctx, id, quantity
func (_m *ProductRepository) AddStock(ctx context.Context, id int64, quantity int64) error {
	ret := _m.Called(ctx, id, quantity)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) error); ok {
		r0 = rf(ctx, id, quantity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateProduct provides a mock function with given fields: ctx, product
func (_m *ProductRepository) CreateProduct(ctx context.Context, product models.Product) error {
	ret := _m.Called(ctx, product)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Product) error); ok {
		r0 = rf(ctx, product)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteProduct provides a mock function with given fields: ctx, id
func (_m *ProductRepository) DeleteProduct(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetProductByID provides a mock function with given fields: ctx, id
func (_m *ProductRepository) GetProductByID(ctx context.Context, id int64) (models.UserProduct, error) {
	ret := _m.Called(ctx, id)

	var r0 models.UserProduct
	if rf, ok := ret.Get(0).(func(context.Context, int64) models.UserProduct); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(models.UserProduct)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProductsByCategory provides a mock function with given fields: ctx, category, page, limit
func (_m *ProductRepository) GetProductsByCategory(ctx context.Context, category string, page int, limit int) ([]models.UserProduct, error) {
	ret := _m.Called(ctx, category, page, limit)

	var r0 []models.UserProduct
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []models.UserProduct); ok {
		r0 = rf(ctx, category, page, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.UserProduct)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) error); ok {
		r1 = rf(ctx, category, page, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStock provides a mock function with given fields: ctx, id
func (_m *ProductRepository) GetStock(ctx context.Context, id int64) (int64, error) {
	ret := _m.Called(ctx, id)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, int64) int64); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReduceStock provides a mock function with given fields: ctx, id, quantity
func (_m *ProductRepository) ReduceStock(ctx context.Context, id int64, quantity int64) error {
	ret := _m.Called(ctx, id, quantity)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) error); ok {
		r0 = rf(ctx, id, quantity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateProduct provides a mock function with given fields: ctx, id, updateFields
func (_m *ProductRepository) UpdateProduct(ctx context.Context, id int64, updateFields map[string]interface{}) error {
	ret := _m.Called(ctx, id, updateFields)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, map[string]interface{}) error); ok {
		r0 = rf(ctx, id, updateFields)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewProductRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewProductRepository creates a new instance of ProductRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProductRepository(t mockConstructorTestingTNewProductRepository) *ProductRepository {
	mock := &ProductRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
