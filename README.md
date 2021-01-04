# cliutil

[![Build Status](https://travis-ci.org/cpliakas/cliutil.svg?branch=main)](https://travis-ci.org/cpliakas/cliutil)
[![Go Reference](https://pkg.go.dev/badge/github.com/cpliakas/cliutil.svg)](https://pkg.go.dev/github.com/cpliakas/cliutil)

Helper functions that simplify writing CLI tools in Golang using the [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) libraries.

## Installation

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
go get github.com/cpliakas/cliutil
```

Next, include cliutil in your application:

```go
import "github.com/cpliakas/cliutil"
```

## Usage

### Flagger

Convenience functions that make it easier to add options to commands when using [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper).

```go

var myCfg *viper.Viper

func init() {

	// Assumes rootCmd and myCmd are defined. We are adding flags to myCmd.
	// See https://github.com/spf13/cobra#create-additional-commands
	rootCmd.AddCommand(myCmd)
    
	// Configure the AutomaticEnv capability to read configuration from
	// environment variables prefixed with "MYAPP_".
	// See https://github.com/spf13/viper#working-with-environment-variables
	myCfg = cliutil.InitConfig("MYAPP")

	// Add flags to myCmd. Use the myCfg.Get* methods to get the options passed
	// via command line. See https://github.com/spf13/viper for usage docs.
	flags := cliutil.NewFlagger(myCmd, myCfg)
	flags.String("log-level", "l", "info", "the minimum log level")
	flags.Int("max-num", "n", 100, "the maximum number of something")
}
```

Or ...

```go

var myCfg *viper.Viper

func init() {
	var flags *cliutil.Flagger
	myCfg, flagger = cliutil.AddCommand(rootCmd, myCmd, "MYAPP")

	flags.String("log-level", "l", "info", "the minimum log level")
	flags.Int("max-num", "n", 100, "the maximum number of something")
}

```

### Option Struct Tags

Set and get options via the `cliutil` struct tag, as shown with the `PrintInput` struct below:

```go
package cmd

import (
	"fmt"

	"github.com/cpliakas/cliutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type PrintInput struct {
	Text string `cliutil:"option=text short=t default='default value' usage='text printed to stdout'"`
}

var printCfg *viper.Viper

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print text to STDOUT",
	Run: func(cmd *cobra.Command, args []string) {
		input := &PrintInput{}
		cliutil.GetOptions(input, printCfg)
		fmt.Println(input.Text)
	},
}

func init() {
	var flags *cliutil.Flagger
	printCfg, flags = cliutil.AddCommand(rootCmd, printCmd, "MYAPP")
	flags.SetOptions(&PrintInput{})
}
```

Assuming `rootCmd` exists and defines the `myapp` command:

```
$> ./myapp print --text hello
hello
```

### Key/Value Parser

Parses strings like `key1=value1 key2="some other value"` into a `map[string]string`.

```go

func parseValues() {
	s := `key1=value1 key2="some other value"`
	m := cliutil.ParseKeyValue(s)
	fmt.Println(m["key1"])  // prints "value1"
	fmt.Println(m["key2"])  // prints "some other value"
}

```

### Event Listener

Listens for shutdown events, useful for long-running processes.

```go
func main() {

	// Start the event listener. A message is sent to the shutdown channel when
	// a SIGINT or SIGTERM signal is received.
	listener := cliutil.NewEventListener().Run()

	// Do something long-running in a goroutine.
	go doStuff()

	// Wait for the shutdown signal.
	listener.Wait()
	log.Println("shutdown signal received, exiting")
}

func doStuff() {
	// do stuff here
}
```

### Leveled Logger with Context

A simple, leveled logger with log tags derived from context. The defaults are inspired by the [best practices](https://dev.splunk.com/enterprise/docs/developapps/logging/loggingbestpractices/) suggested by Splunk.

```go
func main() {
	ctx, logger, _ := cliutil.NewLoggerWithContext(context.Background(), cliutil.LogDebug)
	logger.Debug(ctx, "transaction id created")
	// 2020/04/29 14:24:50.516125 DEBUG message="transaction id created" transid=bqkoscmg10l5tdt068i0

	err := doStuff()
	logger.FatalIfError(ctx, "error doing stuff", err)
	// no-op, will only log the message if err != nil.

	ctx = cliutil.ContextWithLogTag(ctx, "stuff", "done doing it")
	logger.Notice(ctx, "shutdown")
	// 2020/04/29 14:24:50.516140 NOTICE message="shutdown" transid=bqkoscmg10l5tdt068i0 stuff="done doing it"
}

func doStuff() error {

	// do stuff here, returns any errors

	return nil
}
```