package main

import (
	"net/url"
	"testing"
)

func TestOpenDB(t *testing.T) {
	for _, path := range []string{
		"api.db",
		"postgres://bob:secret@1.2.3.4:5432/mydb?sslmode=disable",
		"mysql://bob:secret@tcp(1.2.3.4)/mydb",
	} {
		dbUrl, err := url.Parse(path)
		if err != nil {
			t.Fatal(err)
		}

		if dbUrl.IsAbs() && len(dbUrl.Scheme) > 0 {
			pref := len(dbUrl.Scheme) + len("://")
			conn := path[pref:]
			t.Logf(`sql.OpenDB("%s", "%s")`, dbUrl.Scheme, conn)
		} else {
			t.Logf(`sql.OpenDB("%s", "%s")`, "sqlite", path)
		}
	}
}
