# RCmd

(If you've read the documentation of [github.com/metrumresearchgroup/command](https://github.com/metrumresearchgroup/command), this will all be familiar to you.)

This package wraps the functionality of MRG's `*command.Cmd` with some environmental frameworking to make calling R easier.

## Goal

This package wraps the functionality of MRG's `*command.Cmd` with some environmental frameworking to make calling R easier.

The project's goal is ease of use when configuring, starting/stopping, and directly capturing output of a `R` and `Rscript` calls.

## Use Cases

Not knowing where to start with automating R with the sometimes daunting `*exec.Cmd`, you can use this library to simplify the process. The entry point is `New` (it always requires a context). You can also `Kill()` a process without knowing/caring how it gets done.

## Usage

Simple interactive use:

```go
// skipping out on errors, just noise in an example.
cmd, _ := rcmd.New(context.Background(), "", "--quiet", "-e", "2+2")
co, _ := cmd.CombinedOutput()
```

Run a shiny app (see shiny/demo):
```go
    app, _ := shiny.ConfigureApp(context.Background(), "testdata", "app.R", 0)

	command.InteractiveIO().Apply(app.Cmd)

	_ = app.Start()

	// kill after a wait, or use the .Wait() function if your shiny app 
	// has an exit route.
	if err := app.Kill(); err != nil {
		panic(err)
	}
}
```

Programmatic input, standard output/err:

```go
reader, writer, _ := os.Pipe()
c := rcmd.New(context.Background(), "", "--interactive", "--quiet")
command.WireIO(reader, os.Stdout, os.Stderr).Apply(c)
_ = c.Start()
_, _ = fmt.Fprintln(writer, "2+2")
_, _ = fmt.Fprintln(writer, `q("no")`)
_ = c.Wait()
```

## Availability of base functionality

Everything in `*exec.Cmd` is available. See the [official docs](https://pkg.go.dev/os/exec#Cmd) for expanded help:

```go
type Cmd struct {
	Path string
	Args []string
	Env []string
	Dir string
	Stdin io.Reader
	Stdout io.Writer
	Stderr io.Writer
	ExtraFiles []*os.File
	SysProcAttr *syscall.SysProcAttr
	Process *os.Process
	ProcessState *os.ProcessState
}

func (c *Cmd) CombinedOutput() ([]byte, error)
func (c *Cmd) Output() ([]byte, error)
func (c *Cmd) Run() error
func (c *Cmd) Start() error
func (c *Cmd) StderrPipe() (io.ReadCloser, error)
func (c *Cmd) StdinPipe() (io.WriteCloser, error)
func (c *Cmd) StdoutPipe() (io.ReadCloser, error)
func (c *Cmd) String() string
func (c *Cmd) Wait() error
```

## Additional functionality

We added additional "Kill" functionality to the library for your convenience. As always, you can also cancel the context you're handing off to the Cmd if you want a shortcut.

```go
// Kill ends a process. Its operation depends on whether you created the Cmd
// with a context or not.
Kill() error

// KillTimer waits for the duration stated and then sends back the results
// of calling Kill via the errCh channel.
KillTimer(d time.Duration, errCh chan<- error)

// KillAfter waits until the time stated and then sends back the results
// of calling Kill via the errCh channel.
KillAfter(t time.Time, errCh chan<- error)
```

## Testing

This package only depends upon our own [wrapt](https://github.com/metrumresearchgroup/wrapt/) testing library. Running `make test` is sufficient to verify its contents.

We include .golangci.yml configuration and a .drone.yaml for quality purposes.
