package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/frida/frida-go/frida"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "fnoios",
	Short: "iOS read output",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("missing app name")
		}

		app := args[0]

		dev := frida.USBDevice()
		if dev == nil {
			return errors.New("no USB device detected")
		}
		defer dev.Clean()

		dev.On("output", func(pid, fd int, data []byte) {
			fmt.Printf("[fd=%d] %s", fd, string(data))
		})

		opts := frida.NewSpawnOptions()
		opts.SetStdio(frida.StdioPipe)

		pid, err := dev.Spawn(app, opts)
		if err != nil {
			return err
		}

		session, err := dev.Attach(pid, nil)
		if err != nil {
			return err
		}

		session.On("detached", func(reason frida.SessionDetachReason, crash *frida.Crash) {
			fmt.Printf("detached: %s\n", reason)
		})

		if err := dev.Resume(pid); err != nil {
			return err
		}

		r := bufio.NewReader(os.Stdin)
		r.ReadLine()
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
}

func Execute() error {
	return rootCmd.Execute()
}
