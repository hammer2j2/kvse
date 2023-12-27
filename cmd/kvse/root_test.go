package kvse_test

import (
	"os"
	"testing"

	"github.com/hammer2j2/kvse/cmd/kvse"
)

var (
	err error
)

func TestRootCmdExecute(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	// when calling the root command with no args
	// it should not error
	err = kvse.RootCmd.Execute()
	if err != nil {
		t.Errorf("RootCmd.Execute() returned error: %v", err)
	}
}
