// Slurp is an intermediary to the stored build/blob, used specifically
// to speed up publishing microbox builds.
//
// Usage
//
// To start slurp as a server, simply run:
//
//  slurp
//
// For more specific usage information, refer to the help doc (slurp -h):
//  Usage:
//    slurp [flags]
//
//  Flags:
//    -a, --api-address="https://127.0.0.1:1566": Listen uri for the API (scheme defaults to https)
//    -t, --api-token="secret": Token for API Access
//    -b, --build-dir="/var/db/slurp/build/": Build staging directory
//    -c, --config-file="": Configuration file to load
//    -i, --insecure[=true]: Disable tls certificate verification when connecting to storage
//    -l, --log-level="info": Log level to output [fatal|error|info|debug|trace]
//    -s, --ssh-addr="127.0.0.1:1567": Address ssh server will listen on (ip:port combo)
//    -k, --ssh-host="/var/db/slurp/slurp_rsa": SSH host (private) key file
//    -S, --store-addr="hoarders://127.0.0.1:7410": Storage host address
//    -T, --store-token="": Storage auth token
//    -v, --version[=false]: Print version info and exit
//
package main

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/spf13/cobra"

	"github.com/mu-box/slurp/api"
	"github.com/mu-box/slurp/backend"
	"github.com/mu-box/slurp/config"
	"github.com/mu-box/slurp/ssh"
)

var (
	// slurp provides the slurp cli/server functionality
	slurp = &cobra.Command{
		Use:               "slurp",
		Short:             "slurp - build intermediary",
		Long:              ``,
		SilenceErrors:     true,
		SilenceUsage:      true,
		PersistentPreRunE: readConfig,
		PreRunE:           preFlight,
		RunE:              startSlurp,
	}

	// to be populated by linker
	version string
	commit  string
)

// add cli options to slurp
func init() {
	config.AddFlags(slurp)
}

func readConfig(ccmd *cobra.Command, args []string) error {
	if err := config.LoadConfigFile(); err != nil {
		config.Log.Fatal("Failed to read config - %v", err)
		return fmt.Errorf("")
	}

	return nil
}

func preFlight(ccmd *cobra.Command, args []string) error {
	if config.Version {
		fmt.Printf("slurp %s (%s)\n", version, commit)
		return fmt.Errorf("")
	}

	return nil
}

// start slurp
func startSlurp(ccmd *cobra.Command, args []string) error {
	config.Log = lumber.NewConsoleLogger(lumber.LvlInt(config.LogLevel))

	// initialize backend
	err := backend.Initialize()
	if err != nil {
		config.Log.Fatal("Backend init failed - %v", err)
		return fmt.Errorf("")
	}

	// start ssh server
	err = ssh.Start()
	if err != nil {
		config.Log.Fatal("SSH server start failed - %v", err)
		return fmt.Errorf("")
	}

	// start api
	err = api.StartApi()
	if err != nil {
		config.Log.Fatal("Api start failed - %v", err)
		return fmt.Errorf("")
	}

	return nil
}

func main() {
	slurp.Execute()
}
