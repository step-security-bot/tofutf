package agent

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/go-logr/logr"
	"github.com/hashicorp/go-multierror"
	"github.com/leg100/otf"
)

// Environment is an implementation of an execution environment
var _ otf.Environment = (*Environment)(nil)

// Execution is made available to the Run Job so that it can interact with OTF
// services and write to the local filesystem, use the logger, etc.

// Environment provides an execution environment for a running a run job,
// providing a working directory, capturing logs etc.
type Environment struct {
	otf.JobService

	RunService                  otf.RunService
	ConfigurationVersionService otf.ConfigurationVersionService
	StateVersionService         otf.StateVersionService

	logr.Logger

	// Current working directory
	path string

	// Whether cancelation has been triggered
	canceled bool

	// Cancel context func for currently running func
	cancel context.CancelFunc

	// Current process
	proc *os.Process

	// CLI process output is written to this
	out io.WriteCloser
}

func NewEnvironment(
	logger logr.Logger,
	rs otf.RunService,
	cvs otf.ConfigurationVersionService,
	svs otf.StateVersionService,
	js otf.JobService,
	job otf.Job) (*Environment, error) {

	path, err := os.MkdirTemp("", "otf-plan")
	if err != nil {
		return nil, err
	}

	out := &otf.JobWriter{
		ID:         job.GetID(),
		JobService: js,
		Logger:     logger,
	}

	return &Environment{
		Logger:                      logger,
		RunService:                  rs,
		ConfigurationVersionService: cvs,
		StateVersionService:         svs,
		out:                         out,
		path:                        path,
	}, nil
}

// Execute executes a run (or anything with a Do(env)) and regardless of whether
// it fails, it'll close the environment logs.
func (e *Environment) Execute(run *otf.Run, job otf.Job) (err error) {
	if err := job.Do(run, e); err != nil {
		err = multierror.Append(fmt.Errorf("executing run: %w", err))
	}

	// Mark the logs as fully uploaded
	if err := e.out.Close(); err != nil {
		err = multierror.Append(fmt.Errorf("closing logs: %w", err))
	}

	return err
}

func (e *Environment) GetConfigurationVersionService() otf.ConfigurationVersionService {
	return e.ConfigurationVersionService
}

func (e *Environment) GetStateVersionService() otf.StateVersionService {
	return e.StateVersionService
}

func (e *Environment) GetRunService() otf.RunService {
	return e.RunService
}

func (e *Environment) GetPath() string {
	return e.path
}

// Cancel terminates execution. Force controls whether termination is graceful
// or not. Performed on a best-effort basis - the func or process may have
// finished before they are cancelled, in which case only the next func or
// process will be stopped from executing.
func (e *Environment) Cancel(force bool) {
	e.canceled = true

	e.cancelCLI(force)
	e.cancelFunc(force)
}

// RunCLI executes a CLI process in the executor.
func (e *Environment) RunCLI(name string, args ...string) error {
	if e.canceled {
		return fmt.Errorf("execution canceled")
	}

	cmd := exec.Command(name, args...)
	cmd.Dir = e.path
	cmd.Stdout = e.out

	stderr := new(bytes.Buffer)
	errWriter := io.MultiWriter(e.out, stderr)
	cmd.Stderr = errWriter

	e.proc = cmd.Process

	if err := cmd.Run(); err != nil {
		e.Error(err, "running CLI step", "stderr", stderr.String())
		return err
	}

	return nil
}

// RunFunc invokes a func in the executor.
func (e *Environment) RunFunc(fn otf.EnvironmentFunc) error {
	if e.canceled {
		return fmt.Errorf("execution canceled")
	}

	// Create and store cancel func so func's context can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	e.cancel = cancel

	return fn(ctx, e)
}

func (e *Environment) cancelCLI(force bool) {
	if e.proc == nil {
		return
	}

	if force {
		e.proc.Signal(os.Kill)
	} else {
		e.proc.Signal(os.Interrupt)
	}
}

func (e *Environment) cancelFunc(force bool) {
	// Don't cancel a func's context unless forced
	if !force {
		return
	}
	if e.cancel == nil {
		return
	}
	e.cancel()
}