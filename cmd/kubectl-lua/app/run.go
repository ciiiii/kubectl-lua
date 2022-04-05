package app

import (
	"github.com/spf13/cobra"

	luavm "github.com/ciiiii/kubectl-lua/vm"
)

var runCmd = &cobra.Command{
	Use:   "run <lua_file>",
	Short: "Run a lua file with kube injected Lua vm",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := kubeConfigFlags.ToRESTConfig()
		if err != nil {
			return err
		}
		vm := luavm.NewLuaVM(config)
		return vm.Load(args[0])
	},
}
