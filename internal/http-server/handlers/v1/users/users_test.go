package users_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlexZahvatkin/segments-users-service/internal/http-server/handlers/v1/users"
	"github.com/AlexZahvatkin/segments-users-service/internal/http-server/handlers/v1/users/mocks"
	slogdiscard "github.com/AlexZahvatkin/segments-users-service/internal/lib/logger/handlers"
	"github.com/AlexZahvatkin/segments-users-service/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddUserHandler(t *testing.T) {
	cases := []struct {
		name string
		requestBody string
		statusCode int
	} {
		{
			name: "Valid name",
			requestBody: `{"name": "example"}`,
			statusCode: http.StatusCreated,
		},
		{
			name: "No name in request body",
			requestBody: `{"wrong": "example"}`,
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Short name in request body",
			requestBody: `{"name": "s"}`,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) { 
			t.Parallel()

			userAdderMock := mocks.NewUserAdder(t)
			userAdderMock.On("AddUser", mock.Anything, mock.Anything).Return(models.User{}, nil).Maybe();

			handler := users.AddUserHandler(slogdiscard.NewDiscardLogger(), userAdderMock)
			req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(tc.requestBody)))
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			require.Equal(t, tc.statusCode, rr.Code)
		})
	}
}