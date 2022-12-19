package utils

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	password := "zufar2000"
	hashedPassword, err := HashPassword(password)
	log.Print(hashedPassword)
	require.NoError(t, err)
	require.NotEmpty(t, password)

	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)
}