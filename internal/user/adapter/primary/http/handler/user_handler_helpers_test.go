package handler_test

import (
	"testing"

	handler "github.com/lechitz/AionApi/internal/user/adapter/primary/http/handler"
	"github.com/stretchr/testify/require"
)

func TestCheckRequiredFields(t *testing.T) {
	require.NoError(t, handler.CheckRequiredFields(map[string]string{
		"password": "abc",
		"new":      "def",
	}))

	err := handler.CheckRequiredFields(map[string]string{
		"password": "",
		"new":      "def",
	})
	require.Error(t, err)
}
