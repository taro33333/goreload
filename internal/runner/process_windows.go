//go:build windows

package runner

import (
	"os"
	"os/exec"
)

func prepareCommand(cmd *exec.Cmd) {
	// Windows doesn't support Setpgid in the same way.
	// For now, we rely on default behavior.
}

func interruptProcess(proc *os.Process) error {
	// Windows doesn't support signals well. os.Interrupt is the best we can do.
	return proc.Signal(os.Interrupt)
}

func killProcess(proc *os.Process) error {
	// Force kill.
	return proc.Kill()
}
