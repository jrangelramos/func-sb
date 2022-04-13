package e2e

/*
func TestBash(t *testing.T) {
	opts := termtest.Options{
		CmdName: "/bin/bash",
	}
	cp, err := termtest.NewTest(t, opts)
	require.NoError(t, err, "create console process")
	defer cp.Close()

	cp.SendLine("echo hello world")
	//cp.Expect("hello world")
	cp.SendLine("exit")
	//cp.ExpectExitCode(0)
	cp.SendCtrlC()
	println(cp.Snapshot())
}

*/