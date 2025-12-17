//go:build !windows

package runner

import (
	"os"
	"os/exec"
	"syscall"
)

func prepareCommand(cmd *exec.Cmd) {
	// Set process group so we can kill all child processes.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func interruptProcess(proc *os.Process) error {
	// First try graceful shutdown with SIGINT to process group.
	if err := syscall.Kill(-proc.Pid, syscall.SIGINT); err != nil {
		// Process might already be dead or we can't signal the group.
		if err != syscall.ESRCH {
			// Try SIGINT to single process.
			return proc.Signal(os.Interrupt)
		}
	}
	return nil
}

func killProcess(proc *os.Process) error {
	// Force kill process group.
	if err := syscall.Kill(-proc.Pid, syscall.SIGKILL); err != nil {
		if err != syscall.ESRCH {
			// Try killing the process itself.
			return proc.Kill()
		}
	}
	// Also try killing the process itself just in case
	_ = proc.Kill()
	return nil
}
