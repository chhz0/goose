package cli

import (
	"context"

	"github.com/spf13/cobra"
)

type Exec interface {
	Execute(ctx context.Context) error
}

type executor struct {
	cobra *cobra.Command
}

// Execute implements Exec.
func (e *executor) Execute(ctx context.Context) error {
	return e.cobra.ExecuteContext(ctx)
}

// Commander
type Commander interface {
	Usage() string
	ShortDesc() string
	LongDesc() string

	InitFunc()
	PreFunc(ctx context.Context, args []string) error
	RunFunc(ctx context.Context, args []string) error

	Flags() Flager
	Commanders() []Commander
	Cobra() *cobra.Command
}

type Command struct {
	Use   string
	Short string
	Long  string

	Inits  func() []func()
	PreRun func(ctx context.Context, args []string) error
	Run    func(ctx context.Context, args []string) error

	Args     *cobra.PositionalArgs
	FlagSet  Flager
	Commands []Commander

	cobra *cobra.Command
}

func NewCommand(rcmd Commander) (Exec, error) {
	rbuilder := &commandBuilder{
		commander: rcmd,
	}

	var addCmdBuilder func(cb *commandBuilder, cmder Commander)
	addCmdBuilder = func(cb *commandBuilder, cmder Commander) {
		cb2 := &commandBuilder{
			commander: cmder,
		}
		cb.subCmdBuilders = append(cb.subCmdBuilders, cb2)
		for _, c := range cmder.Commanders() {
			addCmdBuilder(cb2, c)
		}
	}

	for _, cmder := range rcmd.Commanders() {
		addCmdBuilder(rbuilder, cmder)
	}

	if err := rbuilder.build(); err != nil {
		return nil, err
	}

	return &executor{rbuilder.cobra}, nil
}

func (c *Command) Usage() string {
	return c.Use
}

func (c *Command) ShortDesc() string {
	return c.Short
}

func (c *Command) LongDesc() string {
	return c.Long
}

func (c *Command) InitFunc() {
	if c.Inits != nil {
		cobra.OnInitialize(c.Inits()...)
	}
}

func (c *Command) PreFunc(ctx context.Context, args []string) error {
	if c.PreRun != nil {
		return c.PreRun(ctx, args)
	}
	return nil
}

func (c *Command) RunFunc(ctx context.Context, args []string) error {
	if c.Run != nil {
		return c.Run(ctx, args)
	}
	return nil
}

func (c *Command) Flags() Flager {
	if c.FlagSet != nil {
		return c.FlagSet
	}
	return &FlagSet{}
}

func (c *Command) Commanders() []Commander {
	return c.Commands
}

func (c *Command) Cobra() *cobra.Command {
	return c.cobra
}

type commandBuilder struct {
	cobra     *cobra.Command
	commander Commander

	subCmdBuilders []*commandBuilder
}

func (cb *commandBuilder) build() error {
	cb.cobra = &cobra.Command{
		Use:   cb.commander.Usage(),
		Short: cb.commander.ShortDesc(),
		Long:  cb.commander.LongDesc(),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return cb.commander.PreFunc(cmd.Context(), args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cb.commander.RunFunc(cmd.Context(), args)
		},
		SilenceErrors:              true,
		SilenceUsage:               true,
		SuggestionsMinimumDistance: 2,
	}

	cb.commander.InitFunc()
	cb.commander.Flags().ApplyFlags(cb.cobra)

	for _, sub := range cb.subCmdBuilders {
		if err := sub.build(); err != nil {
			return err
		}
		cb.cobra.AddCommand(sub.cobra)
	}

	return nil
}
