package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	openDB(os.Getenv("LOCAL_MYSQL_PSWD"))
	os.Exit(m.Run())
}
