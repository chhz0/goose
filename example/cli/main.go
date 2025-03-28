package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/chhz0/goose/cli"
)

func main() {
	rootOpts := &RootOption{}

	exec, err := cli.NewCommand(
		&cli.Command{
			Use:   "goosecli",
			Short: "create a cli application.",
			Long:  "goosecli is a cli toolkit based on cobra package.",
			Inits: func() []func() {
				var init1 = func() {
					fmt.Println("goosecli init1....")
				}
				var init2 = func() {
					fmt.Println("goosecli init2....")
				}

				return []func(){init1, init2}
			},
			PreRun: func(ctx context.Context, args []string) error {
				fmt.Println("goosecli pre run")
				return nil
			},
			Run: func(ctx context.Context, args []string) error {
				fmt.Println("goosecli run ...")
				return nil
			},
			FlagSet: rootOpts.flags(),
			Commands: []cli.Commander{
				newPrintCmd(),
				newEchoCmd(),
			},
		},
	)
	if err != nil {
		panic(err)
	}

	if err := exec.Execute(context.Background()); err != nil {
		panic(err)
	}
}

func newPrintCmd() cli.Commander {
	printOpts := &PrintOption{}
	return &cli.Command{
		Use:   "print",
		Short: "Print anything to the screen.",
		Long: `print is for printing anything back to the screen.
For many years people have printed back to the screen.`,
		Run: func(ctx context.Context, args []string) error {
			fmt.Printf("print: %s\n", args)
			return nil
		},
		FlagSet: printOpts.flags(),
	}
}

func newEchoCmd() cli.Commander {
	echoOpts := &EchoOption{}
	return &cli.Command{
		Use:   "echo",
		Short: "Echo anything to the screen.",
		Long: `echo is for echoing anything back.
Echo works a lot like print, except it has a child command.`,
		Run: func(ctx context.Context, args []string) error {
			fmt.Println("Echo: " + strings.Join(args, " "))
			return nil
		},
		Commands: []cli.Commander{
			newTimeCmd(),
		},
		FlagSet: echoOpts.flags(),
	}
}

func newTimeCmd() cli.Commander {
	timeOpts := &TimesOption{}
	return &cli.Command{
		Use:   "times [# times] [string to echo]",
		Short: "Echo anything to the screen more times.",
		Long: `echo things multiple times back to the user by providing
a count and a string.`,
		Run: func(ctx context.Context, args []string) error {
			for i := 0; i < 5; i++ {
				fmt.Println("Echo times: " + strings.Join(args, " "))
			}
			return nil
		},
		FlagSet: timeOpts.flags(),
	}
}
