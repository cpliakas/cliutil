# cliutil

[![Build Status](https://travis-ci.org/cpliakas/cliutil.svg?branch=master)](https://travis-ci.org/cpliakas/cliutil)

Helper functions that simplify writing CLI tools in Golang using the
[Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper)
libraries.

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

### Key/Value Parser

Parses a strings like `key1=value1 key2="some other value"` into a
`map[string]string`.

```go

func parseValues() {
	s := `key1=value1 key2="some other value"`
	m := cliutil.ParseKeyValue(s)
	fmt.Println(m["key1"])  // value1
	fmt.Println(m["key2"])  // some other value
}

```

### EventListener

Listens for shutdown events, useful for long-running processes.

```go
func main() {

	// Start the event listener. A message is sent to the shutdown channel when
	// a SIGINT or SIGTERM signal is received.
	sig := make(chan os.Signal)
	shutdown := cliutil.EventListener(sig)

	// Do something long-running in a goroutine.
	go doStuff()

	// Wait for the shutdown signal.
	<-shutdown
	log.Println("shutdown signal received, exiting")
}

func doStuff() {
	// do stuff here
}
```
