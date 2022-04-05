package app

import (
	"github.com/spf13/cobra"

	luavm "github.com/ciiiii/kubectl-lua/vm"
)

var replCmd = &cobra.Command{
	Use:   "repl",
	Short: "Enter a kube injected Lua interpreter",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := kubeConfigFlags.ToRESTConfig()
		if err != nil {
			return err
		}
		vm := luavm.NewLuaVM(config)
		return vm.REPL()
	},
}
