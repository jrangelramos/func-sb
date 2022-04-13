// +build windows

package e2e

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ActiveState/termtest"
)

type TestShellInteractiveCmdRunner struct {
	TestShell *TestShellCmdRunner
	T         *testing.T

	// Sleep interval between each subcommand
	commandSleepInterval time.Duration

	// Sleep interval after last command completion.
	// Required time to give to process to complete before EOF
	completionSleepInterval time.Duration

	// Timeout before kill the cmd in case of some failure
	completionTimeout time.Duration
}

func NewTestShellInteractiveCmdRunner(t *testing.T) *TestShellInteractiveCmdRunner {
	testShell := NewKnFuncShellCli(t)
	return &TestShellInteractiveCmdRunner{
		TestShell:               testShell,
		T:                       t,
		commandSleepInterval:    time.Second * 1,
		completionSleepInterval: time.Second * 2,
		completionTimeout:       time.Second * 15,
	}
}

// Prepare creates a go function used to start kn func (binary) that requires user interaction such as `func config command`
func (f *TestShellInteractiveCmdRunner) PrepareRun(funcCommand ...string) func(args ...string) TestShellCmdResult {

	return func(userInput ...string) TestShellCmdResult {

		// Prepare Command args
		finalArgs := f.TestShell.BinaryArgs
		if finalArgs == nil {
			finalArgs = funcCommand
		} else if funcCommand != nil {
			finalArgs = append(finalArgs, funcCommand...)
		}
		if f.TestShell.ShouldDumpCmdLine {
			f.T.Log(f.TestShell.Binary, strings.Join(finalArgs, " "))
		}

		opts := termtest.Options{
			DefaultTimeout: 0,
			Environment:    append(os.Environ(), f.TestShell.Env...),
			CmdName:        f.TestShell.Binary,
			Args:           finalArgs,
		}
		if f.TestShell.SourceDir != "" {
			opts.WorkDirectory = f.TestShell.SourceDir
		}
		cp, err := termtest.NewTest(f.T, opts)
		if err != nil {
			f.T.Fatal(err)
		}

		defer cp.Close()

		for _, subcmd := range userInput {
			time.Sleep(f.commandSleepInterval)
			cp.Send(subcmd)
		}
		cp.ExpectExitCode(0, f.commandSleepInterval)
		err = cp.Close()
		if err != nil {
			f.T.Logf("error on PTY close: %v\n", err)
		}


		// Collect results
		result := TestShellCmdResult{
			Stdout: cp.Snapshot(),
			Error:  err,
		}
		if err == nil && f.TestShell.ShouldDumpOnSuccess {
			f.T.Log(result.Stdout)
		}
		if err != nil {
			f.T.Log(result.Stdout)
			f.T.Log(err.Error())
			f.T.Log(result.Stderr)
		}
		return result
	}
}
