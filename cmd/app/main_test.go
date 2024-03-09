package main_test

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tedmo/go-rest-template/internal/testdb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tedmo/go-rest-template/internal/app"
	"github.com/tedmo/go-rest-template/internal/http"
	"github.com/tedmo/go-rest-template/internal/postgres"
)

func TestApp(t *testing.T) {

	testServer := NewTestServer(t)

	testClient := testServer.Client()
	baseURL := testServer.URL

	t.Cleanup(func() {
		testServer.Close()
	})

	t.Run("create user", func(t *testing.T) {
		reqBody := `{"name": "test"}`
		req, err := nethttp.NewRequest(nethttp.MethodPost, baseURL+"/users", strings.NewReader(reqBody))
		require.NoError(t, err)

		resp, err := testClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, nethttp.StatusCreated, resp.StatusCode)

		var respBody http.Response[app.User]
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)

		user := respBody.Data
		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, "test", user.Name)
	})

	t.Run("find user by id", func(t *testing.T) {
		req, err := nethttp.NewRequest(nethttp.MethodGet, baseURL+"/users/1", nethttp.NoBody)
		require.NoError(t, err)

		resp, err := testClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, nethttp.StatusOK, resp.StatusCode)

		var respBody http.Response[app.User]
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)

		user := respBody.Data
		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, "test", user.Name)
	})

	t.Run("find users", func(t *testing.T) {
		req, err := nethttp.NewRequest(nethttp.MethodGet, baseURL+"/users", nethttp.NoBody)
		require.NoError(t, err)

		resp, err := testClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, nethttp.StatusOK, resp.StatusCode)

		var respBody http.Response[[]app.User]
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)

		users := respBody.Data
		require.Len(t, users, 1)
		user := users[0]
		assert.Equal(t, int64(1), user.ID)
		assert.Equal(t, "test", user.Name)
	})

}

func NewTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	ctx := context.Background()

	testDB, err := testdb.New(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = testDB.Close(ctx)
	})

	migrationsPath, err := app.RootPath("migrations")
	require.NoError(t, err)

	err = testDB.Migrate(migrationsPath)
	require.NoError(t, err)

	db, err := testDB.Open(ctx)
	require.NoError(t, err)
	t.Cleanup(db.Close)

	server := &http.Server{UserService: postgres.NewUserRepo(db)}

	return httptest.NewServer(server.Routes())
}
