package main

type SSHIdentity struct {
	sshKeyMarshalled string
}

func InitializeGameBasedOnIdentity(id SSHIdentity) SudokuGameWrapperState {
	// stub: return a new wrapper
	return NewWrapper()

	// todo: search key in db
	// if match, return existing wrapper
}
