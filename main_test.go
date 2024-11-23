package main

import (
	"os"
	"testing"

	_ "github.com/anton2920/gofa/time"
	"github.com/anton2920/gofa/trace"
)

func TestMain(m *testing.M) {
	trace.BeginProfile()
	ret := m.Run()
	trace.EndAndPrintProfile()

	os.Exit(ret)
}
