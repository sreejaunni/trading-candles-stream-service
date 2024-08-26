package datastore

import (
	"binance-candlestick-service/internal/business/ohlc"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockDatabaseExecutor struct {
	mock.Mock
}

func (m *MockDatabaseExecutor) NamedExec(query string, arg interface{}) (sql.Result, error) {
	args := m.Called(query, arg)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockDatabaseExecutor) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	return nil, nil
}

type mockResult struct {
	LastInsertIdVal int64
	RowsAffectedVal int64
}

func (r *mockResult) LastInsertId() (int64, error) {
	return r.LastInsertIdVal, nil
}

func (r *mockResult) RowsAffected() (int64, error) {
	return r.RowsAffectedVal, nil
}

func TestSaveOHLC_Success(t *testing.T) {
	mockDB := new(MockDatabaseExecutor)
	repo := NewOHLCRepo(mockDB)

	ohlcData := ohlc.OHLC{
		Symbol:         "BTCUSDT",
		Open:           30000,
		High:           31000,
		Low:            29000,
		Close:          30500,
		Volume:         1500,
		StartTimeStamp: time.Now().UTC(),
		EndTimeStamp:   time.Now().UTC(),
	}
	result := &mockResult{
		LastInsertIdVal: 1,
		RowsAffectedVal: 1,
	}
	mockDB.On("NamedExec", mock.Anything, mock.Anything).Return(result, nil)

	err := repo.SaveOHLC(ohlcData)

	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestSaveOHLC_Failure(t *testing.T) {
	mockDB := new(MockDatabaseExecutor)
	repo := NewOHLCRepo(mockDB)

	ohlcData := ohlc.OHLC{
		Symbol:         "BTCUSDT",
		Open:           30000,
		High:           31000,
		Low:            29000,
		Close:          30500,
		Volume:         1500,
		StartTimeStamp: time.Now().UTC(),
		EndTimeStamp:   time.Now().UTC(),
	}

	result := &mockResult{}

	mockDB.On("NamedExec", mock.Anything, mock.Anything).Return(result, errors.New("database error"))

	err := repo.SaveOHLC(ohlcData)

	assert.Error(t, err)
	assert.Equal(t, "error saving OHLC data: database error", err.Error())
	mockDB.AssertExpectations(t)
}
