package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDbUrl(t *testing.T) {
	require := require.New(t)

	for url, exp := range map[string][2]string{
		"api.db": [2]string{"sqlite", "api.db"},
		"postgres://bob:secret@1.2.3.4:5432/mydb?sslmode=disable": [2]string{"postgres", "bob:secret@1.2.3.4:5432/mydb?sslmode=disable"},
		"mysql://bob:secret@tcp(1.2.3.4:3006)/mydb":               [2]string{"mysql", "bob:secret@tcp(1.2.3.4:3006)/mydb"},
	} {
		driver, conn := parseDbUrl(url)
		require.Equal(exp[0], driver)
		require.Equal(exp[1], conn)
	}
}
