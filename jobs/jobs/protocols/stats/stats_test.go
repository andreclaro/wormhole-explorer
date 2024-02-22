package stats_test

import (
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/mock"
	"github.com/wormhole-foundation/wormhole-explorer/jobs/jobs/protocols/internal/commons/mocks"
	"github.com/wormhole-foundation/wormhole-explorer/jobs/jobs/protocols/stats"
	"go.uber.org/zap"
	"io"
	"net/http"
	"testing"
)

func Test_ProtocolsStatsJob_Succeed(t *testing.T) {
	var mockErr error
	statsFetcher := &mockStatsFetch{}
	statsFetcher.On("Get", mock.Anything).Return(stats.Stats{}, mockErr)
	statsFetcher.On("ProtocolName", mock.Anything).Return("protocol_test")
	mockWriterDB := &mocks.MockWriterApi{}
	mockWriterDB.On("WritePoint", mock.Anything, mock.Anything).Return(mockErr)

	job := stats.NewProtocolsStatsJob(mockWriterDB, zap.NewNop(), "v1", statsFetcher)
	resultErr := job.Run(context.Background())
	assert.Nil(t, resultErr)
}

func Test_ProtocolsStatsJob_FailFetching(t *testing.T) {
	var mockErr error
	statsFetcher := &mockStatsFetch{}
	statsFetcher.On("Get", mock.Anything).Return(stats.Stats{}, errors.New("mocked_error_fetch"))
	statsFetcher.On("ProtocolName", mock.Anything).Return("protocol_test")
	mockWriterDB := &mocks.MockWriterApi{}
	mockWriterDB.On("WritePoint", mock.Anything, mock.Anything).Return(mockErr)

	job := stats.NewProtocolsStatsJob(mockWriterDB, zap.NewNop(), "v1", statsFetcher)
	resultErr := job.Run(context.Background())
	assert.NotNil(t, resultErr)
	assert.Equal(t, "mocked_error_fetch", resultErr.Error())
}

func Test_ProtocolsStatsJob_FailedUpdatingDB(t *testing.T) {
	var mockErr error
	statsFetcher := &mockStatsFetch{}
	statsFetcher.On("Get", mock.Anything).Return(stats.Stats{}, mockErr)
	statsFetcher.On("ProtocolName", mock.Anything).Return("protocol_test")
	mockWriterDB := &mocks.MockWriterApi{}
	mockWriterDB.On("WritePoint", mock.Anything, mock.Anything).Return(errors.New("mocked_error_update_db"))

	job := stats.NewProtocolsStatsJob(mockWriterDB, zap.NewNop(), "v1", statsFetcher)
	resultErr := job.Run(context.Background())
	assert.NotNil(t, resultErr)
	assert.Equal(t, "mocked_error_update_db", resultErr.Error())
}

func Test_HttpRestClientStats_FailRequestCreation(t *testing.T) {

	a := stats.NewHttpRestClientStats("protocol_test", "localhost", zap.NewNop(),
		mockHttpClient(func(req *http.Request) (*http.Response, error) {
			return nil, nil
		}))
	_, err := a.Get(nil) // passing ctx nil to force request creation error
	assert.NotNil(t, err)
}

func Test_HttpRestClientStats_FailedRequestExecution(t *testing.T) {

	a := stats.NewHttpRestClientStats("protocol_test", "localhost", zap.NewNop(),
		mockHttpClient(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("mocked_http_client_do")
		}))
	_, err := a.Get(context.Background())
	assert.NotNil(t, err)
	assert.Equal(t, "mocked_http_client_do", err.Error())
}

func Test_HttpRestClientStats_Status500(t *testing.T) {

	a := stats.NewHttpRestClientStats("protocol_test", "localhost", zap.NewNop(),
		mockHttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString("response_body_test")),
			}, nil
		}))
	_, err := a.Get(context.Background())
	assert.NotNil(t, err)
	assert.Equal(t, "failed retrieving client stats from url:localhost - status_code:500 - response_body:response_body_test", err.Error())
}

func Test_HttpRestClientStats_Status200_FailedReadBody(t *testing.T) {

	a := stats.NewHttpRestClientStats("protocol_test", "localhost", zap.NewNop(),
		mockHttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       &mockFailReadCloser{},
			}, nil
		}))
	_, err := a.Get(context.Background())
	assert.NotNil(t, err)
	assert.Equal(t, "failed reading response body from client stats. url:localhost - status_code:200: mocked_fail_read", err.Error())
}

func Test_HttpRestClientStats_Status200_FailedParsing(t *testing.T) {

	a := stats.NewHttpRestClientStats("protocol_test", "localhost", zap.NewNop(),
		mockHttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("this should be a json")),
			}, nil
		}))
	_, err := a.Get(context.Background())
	assert.NotNil(t, err)
	assert.Equal(t, "failed unmarshalling response body from client stats. url:localhost - status_code:200 - response_body:this should be a json: invalid character 'h' in literal true (expecting 'r')", err.Error())
}

func Test_HttpRestClientStats_Status200_Succeed(t *testing.T) {

	a := stats.NewHttpRestClientStats("protocol_test", "localhost", zap.NewNop(),
		mockHttpClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("{\"total_value_locked\":\"123\",\"total_messages\":\"456\"}")),
			}, nil
		}))
	resp, err := a.Get(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, float64(123), resp.TotalValueLocked)
	assert.Equal(t, uint64(456), resp.TotalMessages)
}

type mockStatsFetch struct {
	mock.Mock
}

func (m *mockStatsFetch) Get(ctx context.Context) (stats.Stats, error) {
	args := m.Called(ctx)
	return args.Get(0).(stats.Stats), args.Error(1)
}

func (m *mockStatsFetch) ProtocolName() string {
	args := m.Called()
	return args.String(0)
}

type mockHttpClient func(req *http.Request) (*http.Response, error)

func (m mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return m(req)
}

type mockFailReadCloser struct {
}

func (m *mockFailReadCloser) Read(p []byte) (n int, err error) {
	return 0, errors.New("mocked_fail_read")
}

func (m *mockFailReadCloser) Close() error {
	return nil
}
