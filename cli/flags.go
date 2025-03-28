package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Flager interface {
	ApplyFlags(ccmd *cobra.Command)
}

type FlagSet struct {
	Local      func(pfs *pflag.FlagSet)
	Persistent func(pfs *pflag.FlagSet)
}

func (fs *FlagSet) ApplyFlags(ccmd *cobra.Command) {
	if fs.Local != nil {
		fs.Local(ccmd.Flags())
	}
	if fs.Persistent != nil {
		fs.Persistent(ccmd.PersistentFlags())
	}
}
