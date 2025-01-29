package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/frida/frida-go/frida"
	"github.com/spf13/cobra"
	"os"
)

const (
	chSize = 512
)

type output struct {
	fd   int
	data []byte
}

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

		ch := make(chan *output, chSize)

		go func() {
			for out := range ch {
				fmt.Printf("[fd=%d] %s", out.fd, out.data)
			}
		}()

		dev.On("output", func(pid, fd int, data []byte) {
			ch <- &output{
				fd:   fd,
				data: data,
			}
		})

		dev.On("process_crashed", func(crash *frida.Crash) {
			fmt.Printf("process crashed: %s\n", crash)
			os.Exit(1)
		})

		dev.On("lost", func() {
			fmt.Println("device lost")
			os.Exit(1)
		})

		opts := frida.NewSpawnOptions()
		opts.SetStdio(frida.StdioPipe)

		pid, err := dev.Spawn(app, opts)
		if err != nil {
			return err
		}

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
