package app

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	kubeConfigFlags = genericclioptions.NewConfigFlags(false)
	rootCmd         = &cobra.Command{
		Use:   "kubectl-lua [command]",
		Short: "Query and operate resources with Lua",
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(replCmd)
}

func Execute() {
	flags := pflag.NewFlagSet("kubectl-lua", pflag.ExitOnError)
	pflag.CommandLine = flags

	kubeConfigFlags.AddFlags(flags)
	flags.AddFlagSet(rootCmd.PersistentFlags())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
