package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

func terminateProcess(pid int) error {
	// Signal the process group (-pid), not just the process, so that the process
	// and all its children are signaled. Else, child procs can keep running and
	// keep the stdout/stderr fd open and cause cmd.Wait to hang.
	return syscall.Kill(-pid, syscall.SIGTERM)
}

func setProcessGroupID(cmd *exec.Cmd) {
	// Set process group ID so the cmd and all its children become a new
	// process group. This allows Stop to SIGTERM the cmd's process group
	// without killing this process (i.e. this code here).
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

type Status struct {
	Cmd      string
	PID      int
	Complete bool     // false if stopped or signaled
	Exit     int      // exit code of process
	Error    error    // Go error
	StartTs  int64    // Unix ts (nanoseconds), zero if Cmd not started
	StopTs   int64    // Unix ts (nanoseconds), zero if Cmd not started or running
	Runtime  float64  // seconds, zero if Cmd not started
	Stdout   []string // buffered STDOUT; see Cmd.Status for more info
	Stderr   []string // buffered STDERR; see Cmd.Status for more info
}

type Cmd struct {
	Name  string
	Args  []string
	Env   []string
	Dir   string
	Stdin io.Reader
	*sync.Mutex
	started    bool          // cmd.Start called, no error
	stopped    bool          // Stop called
	done       bool          // run() done
	final      bool          // status finalized in Status
	startTime  time.Time     // if started true
	stdout     *OutputBuffer // low-level stdout buffering and streaming
	stderr     *OutputBuffer // low-level stderr buffering and streaming
	status     Status
	statusChan chan Status   // nil until Start() called
	doneChan   chan struct{} // closed when done running
	buffered   bool          // buffer STDOUT and STDERR to Status.Stdout and Std
}

func NewCmd(name string, args ...string) *Cmd {
	return &Cmd{
		Name:     name,
		Args:     args,
		buffered: true,
		Mutex:    &sync.Mutex{},
		status: Status{
			Cmd:      name,
			PID:      0,
			Complete: false,
			Exit:     -1,
			Error:    nil,
			Runtime:  0,
		},
		doneChan: make(chan struct{}),
	}
}

func (c *Cmd) Start() <-chan Status {
	c.Lock()
	defer c.Unlock()

	if c.statusChan != nil {
		return c.statusChan
	}

	c.statusChan = make(chan Status, 1)
	go c.run()
	return c.statusChan
}

func (c *Cmd) Stop() error {
	c.Lock()
	defer c.Unlock()

	// Nothing to stop if Start hasn't been called, or the proc hasn't started,
	// or it's already done.
	if c.statusChan == nil || !c.started || c.done {
		return nil
	}

	// Flag that command was stopped, it didn't complete. This results in
	// status.Complete = false
	c.stopped = true

	// Signal the process group (-pid), not just the process, so that the process
	// and all its children are signaled. Else, child procs can keep running and
	// keep the stdout/stderr fd open and cause cmd.Wait to hang.
	return terminateProcess(c.status.PID)
}

func (c *Cmd) Status() Status {
	c.Lock()
	defer c.Unlock()

	// Return default status if cmd hasn't been started
	if c.statusChan == nil || !c.started {
		return c.status
	}

	if c.done {
		// No longer running
		if !c.final {
			if c.buffered {
				c.status.Stdout = c.stdout.Lines()
				c.status.Stderr = c.stderr.Lines()
				c.stdout = nil // release buffers
				c.stderr = nil
			}
			c.final = true
		}
	} else {
		// Still running
		c.status.Runtime = time.Now().Sub(c.startTime).Seconds()
		if c.buffered {
			c.status.Stdout = c.stdout.Lines()
			c.status.Stderr = c.stderr.Lines()
		}
	}

	return c.status
}

func (c *Cmd) Done() <-chan struct{} {
	return c.doneChan
}

func (c *Cmd) run() {
	defer func() {
		c.statusChan <- c.Status() // unblocks Start if caller is waiting
		close(c.doneChan)
	}()

	// //////////////////////////////////////////////////////////////////////
	// Setup command
	// //////////////////////////////////////////////////////////////////////
	cmd := exec.Command(c.Name, c.Args...)

	// Platform-specific SysProcAttr management
	setProcessGroupID(cmd)

	// Write stdout and stderr to buffers that are safe to read while writing
	// and don't cause a race condition.
	c.stdout = NewOutputBuffer()
	c.stderr = NewOutputBuffer()
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr

	// Set the runtime environment for the command as per os/exec.Cmd.  If Env
	// is nil, use the current process' environment.
	cmd.Env = c.Env
	cmd.Dir = c.Dir
	cmd.Stdin = c.Stdin

	// //////////////////////////////////////////////////////////////////////
	// Start command
	// //////////////////////////////////////////////////////////////////////
	now := time.Now()
	if err := cmd.Start(); err != nil {
		c.Lock()
		c.status.Error = err
		c.status.StartTs = now.UnixNano()
		c.status.StopTs = time.Now().UnixNano()
		c.done = true
		c.Unlock()
		return
	}

	// Set initial status
	c.Lock()
	c.startTime = now              // command is running
	c.status.PID = cmd.Process.Pid // command is running
	c.status.StartTs = now.UnixNano()
	c.started = true
	c.Unlock()

	// //////////////////////////////////////////////////////////////////////
	// Wait for command to finish or be killed
	// //////////////////////////////////////////////////////////////////////
	err := cmd.Wait()
	now = time.Now()

	// Get exit code of the command. According to the manual, Wait() returns:
	// "If the command fails to run or doesn't complete successfully, the error
	// is of type *ExitError. Other error types may be returned for I/O problems."
	exitCode := 0
	signaled := false
	if err != nil && fmt.Sprintf("%T", err) == "*exec.ExitError" {
		// This is the normal case which is not really an error. It's string
		// representation is only "*exec.ExitError". It only means the cmd
		// did not exit zero and caller should see ExitError.Stderr, which
		// we already have. So first we'll have this as the real/underlying
		// type, then discard err so status.Error doesn't contain a useless
		// "*exec.ExitError". With the real type we can get the non-zero
		// exit code and determine if the process was signaled, which yields
		// a more specific error message, so we set err again in that case.
		exiterr := err.(*exec.ExitError)
		err = nil
		if waitStatus, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			exitCode = waitStatus.ExitStatus() // -1 if signaled
			if waitStatus.Signaled() {
				signaled = true
				err = errors.New(exiterr.Error()) // "signal: terminated"
			}
		}
	}

	// Set final status
	c.Lock()
	if !c.stopped && !signaled {
		c.status.Complete = true
	}
	c.status.Runtime = now.Sub(c.startTime).Seconds()
	c.status.StopTs = now.UnixNano()
	c.status.Exit = exitCode
	c.status.Error = err
	c.done = true
	c.Unlock()
}

type OutputBuffer struct {
	buf   *bytes.Buffer
	lines []string
	*sync.Mutex
}

func NewOutputBuffer() *OutputBuffer {
	out := &OutputBuffer{
		buf:   &bytes.Buffer{},
		lines: []string{},
		Mutex: &sync.Mutex{},
	}
	return out
}

func (rw *OutputBuffer) Write(p []byte) (n int, err error) {
	rw.Lock()
	n, err = rw.buf.Write(p) // and bytes.Buffer implements io.Writer
	rw.Unlock()
	return // implicit
}

func (rw *OutputBuffer) Lines() []string {
	rw.Lock()
	// Scanners are io.Readers which effectively destroy the buffer by reading
	// to EOF. So once we scan the buf to lines, the buf is empty again.
	s := bufio.NewScanner(rw.buf)
	for s.Scan() {
		rw.lines = append(rw.lines, s.Text())
	}
	rw.Unlock()
	return rw.lines
}

const (
	// DEFAULT_LINE_BUFFER_SIZE is the default size of the OutputStream line buffer.
	// The default value is usually sufficient, but if ErrLineBufferOverflow errors
	// occur, try increasing the size by calling OutputBuffer.SetLineBufferSize.
	DEFAULT_LINE_BUFFER_SIZE = 16384

	// DEFAULT_STREAM_CHAN_SIZE is the default string channel size for a Cmd when
	// Options.Streaming is true. The string channel size can have a minor
	// performance impact if too small by causing OutputStream.Write to block
	// excessively.
	DEFAULT_STREAM_CHAN_SIZE = 1000
)
