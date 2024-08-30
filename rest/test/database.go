package test

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type DatabaseMock struct {
	mock.Mock
}

func (m *DatabaseMock) Query(c context.Context, s string, a ...any) (pgx.Rows, error) {
	args := m.Called(c, s, a)
	return args.Get(0).(pgx.Rows), args.Error(1)
}

func (m *DatabaseMock) Exec(c context.Context, s string, a ...any) (pgconn.CommandTag, error) {
	args := m.Called(c, s, a)
	return args.Get(0).(pgconn.CommandTag), args.Error(1)
}

func (m *DatabaseMock) QueryRow(c context.Context, s string, a ...any) pgx.Row {
	args := m.Called(c, s, a)
	return args.Get(0).(pgx.Row)
}

func (m *DatabaseMock) Ping(c context.Context) error {
	args := m.Called(c)
	return args.Error(0)
}
