package cmd

import (
	"errors"
	"fmt"
	"github.com/frida/frida-go/frida"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strings"
)

var scriptContent string

var rootCmd = &cobra.Command{
	Use:   "fduplicator",
	Short: "Intercept FD writes",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		dev := frida.USBDevice()
		if dev == nil {
			return errors.New("no USB device detected")
		}
		defer dev.Clean()

		spawn, err := cmd.Flags().GetBool("spawn")
		if err != nil {
			return err
		}

		app, err := cmd.Flags().GetString("app")
		if err != nil {
			return err
		}

		pd, err := cmd.Flags().GetInt("pid")
		if err != nil {
			return err
		}

		if app == "" && pd == -1 {
			return errors.New("you need to specify either --app or --pid")
		}

		pid := 0

		var session *frida.Session
		if app != "" {
			if spawn {
				spawnedPid, err := dev.Spawn(app, nil)
				if err != nil {
					return err
				}
				pid = spawnedPid
				s, err := dev.Attach(spawnedPid, nil)
				if err != nil {
					return err
				}
				session = s
			} else {
				s, err := dev.Attach(app, nil)
				if err != nil {
					return err
				}
				session = s
			}
		} else {
			s, err := dev.Attach(pd, nil)
			if err != nil {
				return err
			}
			session = s
		}

		script, err := session.CreateScript(scriptContent)
		if err != nil {
			return err
		}

		script.On("message", func(message string) {
			if !strings.Contains(message, "can't decode byte") {
				msg, err := frida.ScriptMessageToMessage(message)
				if err != nil {
					return
				}
				mappedPayload := msg.Payload.(map[string]any)
				fmt.Printf("%s", mappedPayload["data"])
			}
		})

		if err := script.Load(); err != nil {
			return err
		}

		if spawn {
			dev.Resume(pid)
		}

		script.ExportsCall("start")

		go func() {
			for {
				_ = script.ExportsCall("read")
			}
		}()

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		for sig := range c {
			_ = sig
			fmt.Println("[*] Exiting...\nUnloading script")
			if err := script.Unload(); err != nil {
				return err
			}
			fmt.Println("Script unloaded")
			return nil
		}
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.Flags().StringP("app", "a", "", "Application name to attach to")
	rootCmd.Flags().IntP("pid", "p", -1, "PID of process to attach to")
	rootCmd.Flags().BoolP("spawn", "s", false, "Spawn the app/file")
}

func Execute(script string) error {
	scriptContent = script
	return rootCmd.Execute()
}
