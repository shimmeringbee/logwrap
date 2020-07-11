package logwrap

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockImpl struct {
	mock.Mock
}

func (l *MockImpl) Impl(ctx context.Context, msg Message) {
	l.Called(ctx, msg)
}
