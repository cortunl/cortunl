package runner

import (
	"github.com/cortunl/cortunl/constants"
	"github.com/dropbox/godropbox/errors"
	"os"
	"os/exec"
)

type Runner struct {
	exited  bool
	command *exec.Cmd
	OnError func(error)
}

func (r *Runner) Run(cmd *exec.Cmd, onExit func()) (err error) {
	if r.command != nil {
		panic("runner: Already started")
	}

	r.command = cmd
	r.command.Stdout = os.Stdout
	r.command.Stderr = os.Stdout

	err = r.command.Start()
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "runner: Failed to start process"),
		}
		r.exited = true
		return
	}

	go func() {
		err = r.command.Wait()
		if err != nil && !r.exited {
			err = &constants.ExecError{
				errors.Wrap(err, "runner: Process exec error"),
			}

			if r.OnError != nil {
				r.OnError(err)
			}
		}
		r.exited = true

		onExit()
	}()

	return
}

func (r *Runner) Stop() (err error) {
	if r.command == nil || r.exited {
		return
	}

	r.exited = true
	err = r.command.Process.Kill()
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "runner: Failed to stop process"),
		}
		return
	}

	return
}

func (r *Runner) Wait() (err error) {
	if r.command == nil {
		return
	}

	err = r.command.Wait()
	if err != nil {
		err = &constants.ExecError{
			errors.Wrap(err, "runner: Process exec error"),
		}
		return
	}

	return
}
