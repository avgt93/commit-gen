package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/avgt93/commit-gen/internal/config"
	"github.com/avgt93/commit-gen/internal/opencode"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "commit-gen",
	Short: "Generate commit messages using OpenCode AI",
	Long: `commit-gen is a CLI tool that generates descriptive commit messages
based on your staged git changes using OpenCode's AI capabilities.

Simply run 'git commit -m ""' and it will fill in the message for you!`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/commit-gen/config.yaml)")

	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(previewCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(healthCmd)

	cacheCmd.AddCommand(cacheStatusCmd)
	cacheCmd.AddCommand(cacheClearCmd)
	rootCmd.AddCommand(cacheCmd)

	generateCmd.Flags().StringP("style", "s", "conventional", "Commit message style (conventional, imperative, detailed)")
	generateCmd.Flags().StringP("mode", "m", "", "Operation mode: 'run' (default) or 'server'")
	generateCmd.Flags().BoolP("no-confirm", "n", false, "Skip confirmation prompt and use generated message directly")
	generateCmd.Flags().Bool("dry-run", false, "Show message without writing to git")
	generateCmd.Flags().Bool("hook", false, "Internal flag for git hook usage")
	generateCmd.Flags().Bool("ignore-server-check", false, "Skip checking if OpenCode backend is available")

	previewCmd.Flags().StringP("mode", "m", "", "Operation mode: 'run' (default) or 'server'")
	previewCmd.Flags().Bool("ignore-server-check", false, "Skip checking if OpenCode backend is available")
}

func initConfig() {
	_ = config.Initialize(cfgFile)
}

func checkBackendAvailability(cfg *config.Config, ignoreCheck bool) error {
	if ignoreCheck {
		return nil
	}

	mode := cfg.OpenCode.Mode
	if mode == "" {
		mode = "run"
	}

	if mode == "server" {
		return checkOpenCodeHealth(cfg)
	}
	return checkOpenCodeRunner()
}

func checkOpenCodeRunner() error {
	runner := opencode.NewRunner(10)
	available, err := runner.CheckAvailable()
	if err != nil || !available {
		return fmt.Errorf("opencode binary not found in PATH. Please install opencode first")
	}
	return nil
}

func checkOpenCodeHealth(cfg *config.Config) error {
	client := opencode.NewClient(
		cfg.OpenCode.Host,
		cfg.OpenCode.Port,
		cfg.OpenCode.Timeout,
	)

	healthy, err := client.CheckHealth()
	if err == nil && healthy {
		return nil
	}

	cmd := exec.Command(
		"opencode",
		"serve",
		"--port", strconv.Itoa(cfg.OpenCode.Port),
	)

	setSysProcAttr(cmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		ErrServerNotRunning := errors.New("opencode server is not running")
		return fmt.Errorf(
			"%w at %s:%d: %v",
			ErrServerNotRunning,
			cfg.OpenCode.Host,
			cfg.OpenCode.Port,
			err,
		)
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Printf("opencode server exited: %v\n", err)
		}
	}()

	time.Sleep(2 * time.Second)

	healthy, err = client.CheckHealth()
	if err != nil || !healthy {
		return fmt.Errorf("opencode server failed to become healthy")
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
