package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/avgt93/commit-gen/internal/cache"
	"github.com/avgt93/commit-gen/internal/config"
	"github.com/avgt93/commit-gen/internal/generator"
	"github.com/avgt93/commit-gen/internal/git"
	"github.com/avgt93/commit-gen/internal/hook"
	"github.com/avgt93/commit-gen/internal/opencode"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   color.CyanString("commit-gen"),
	Short: color.GreenString("Generate commit messages using OpenCode AI"),
	Long: color.YellowString(`commit-gen is a CLI tool that generates descriptive commit messages
based on your staged git changes using OpenCode's AI capabilities.

Simply run 'git commit -m ""' and it will fill in the message for you!`),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			err := cmd.Help()
			if err != nil {
				fmt.Println(err)
			}
		}
	},
}

var generateCmd = &cobra.Command{
	Use:   color.CyanString("generate"),
	Short: color.GreenString("Generate a commit message from staged changes"),
	Long: color.YellowString(`Generate a commit message from your currently staged git changes.
The message will be generated using OpenCode's AI based on the diff.`),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()

		if modeFlag, _ := cmd.Flags().GetString("mode"); modeFlag != "" {
			cfg.OpenCode.Mode = modeFlag
		}

		ignoreCheck, _ := cmd.Flags().GetBool("ignore-server-check")
		if err := checkBackendAvailability(cfg, ignoreCheck); err != nil {
			return err
		}

		cacheDir := filepath.Join(os.Getenv("HOME"), ".cache", "commit-gen")
		sessionCache := cache.GetCache(24*time.Hour, cacheDir)

		gen := generator.NewGenerator(cfg, sessionCache)

		message, err := gen.Generate()
		if err != nil {
			color.Red("Error: %v", err)
			return err
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		isHook, _ := cmd.Flags().GetBool("hook")

		if dryRun || isHook {
			fmt.Println(message)
		} else {
			if err := git.WriteCommitMessage(message); err != nil {
				return fmt.Errorf("failed to write commit message: %w", err)
			}
			color.Green("✓ Commit message generated:")
			fmt.Printf("  %s\n", message)
		}

		return nil
	},
}

var installCmd = &cobra.Command{
	Use:   color.CyanString("install"),
	Short: color.GreenString("Install git hook for automatic commit message generation"),
	Long: color.YellowString(`Installs a prepare-commit-msg git hook in the current repository.
This allows automatic commit message generation when running 'git commit -m ""'.`),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		if err := hook.Install(cfg.Git.Editor); err != nil {
			color.Red("Error: %v", err)
			return err
		}
		color.Green("✓ Git hook installed successfully")
		fmt.Println("Now you can use: git commit -m \"\"")
		return nil
	},
}

var uninstallCmd = &cobra.Command{
	Use:   color.CyanString("uninstall"),
	Short: color.GreenString("Remove the git hook"),
	Long:  color.YellowString(`Removes the prepare-commit-msg git hook from the current repository.`),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := hook.Uninstall(); err != nil {
			color.Red("Error: %v", err)
			return err
		}
		color.Green("✓ Git hook removed successfully")
		return nil
	},
}

var configCmd = &cobra.Command{
	Use:   color.CyanString("config"),
	Short: color.GreenString("Manage configuration"),
	Long:  color.YellowString(`View and modify commit-gen configuration.`),

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()

		color.Cyan("OpenCode Configuration:")
		fmt.Printf("  Mode: %s\n", cfg.OpenCode.Mode)
		fmt.Printf("  Host: %s (server mode only)\n", cfg.OpenCode.Host)
		fmt.Printf("  Port: %d (server mode only)\n", cfg.OpenCode.Port)
		fmt.Printf("  Timeout: %ds\n", cfg.OpenCode.Timeout)

		color.Cyan("\nGeneration Configuration:")
		fmt.Printf("  Style: %s\n", cfg.Generation.Style)
		fmt.Printf("  Provider: %s\n", cfg.Generation.Model.Provider)
		fmt.Printf("  Model: %s\n", cfg.Generation.Model.ModelID)

		color.Cyan("\nCache Configuration:")
		fmt.Printf("  Enabled: %v (server mode only)\n", cfg.Cache.Enabled)
		fmt.Printf("  TTL: %s\n", cfg.Cache.TTL)

		color.Cyan("\nGit Configuration:")
		fmt.Printf("  Editor: %s\n", cfg.Git.Editor)
		fmt.Printf("  Staged Only: %v\n", cfg.Git.StagedOnly)
		fmt.Printf("  Max Diff Size: %d bytes (%dKB)\n", cfg.Git.MaxDiffSize, cfg.Git.MaxDiffSize/1024)

		return nil
	},
}

var previewCmd = &cobra.Command{
	Use:   color.CyanString("preview"),
	Short: color.GreenString("Preview changes and generated commit message"),
	Long:  color.YellowString(`Shows your staged changes and what commit message would be generated.`),
	RunE: func(cmd *cobra.Command, args []string) error {
		diff, err := git.GetStagedDiff()
		if err != nil {
			color.Red("Error: %v", err)
			return err
		}

		if diff == "" {
			color.Yellow("No staged changes found")
			return nil
		}

		color.Cyan("=== Staged Changes ===")
		fmt.Println(diff)

		color.Cyan("\n=== Generated Commit Message ===")

		cfg := config.Get()

		if modeFlag, _ := cmd.Flags().GetString("mode"); modeFlag != "" {
			cfg.OpenCode.Mode = modeFlag
		}

		ignoreCheck, _ := cmd.Flags().GetBool("ignore-server-check")
		if err := checkBackendAvailability(cfg, ignoreCheck); err != nil {
			return err
		}

		cacheDir := filepath.Join(os.Getenv("HOME"), ".cache", "commit-gen")
		sessionCache := cache.GetCache(24*time.Hour, cacheDir)

		gen := generator.NewGenerator(cfg, sessionCache)
		message, err := gen.Generate()
		if err != nil {
			color.Red("Error generating message: %v", err)
			return err
		}

		color.Green(message)
		return nil
	},
}

var cacheCmd = &cobra.Command{
	Use:   color.CyanString("cache"),
	Short: color.GreenString("Manage session cache"),
	Long:  color.YellowString(`View and manage the OpenCode session cache.`),
}

var cacheStatusCmd = &cobra.Command{
	Use:   color.CyanString("status"),
	Short: color.GreenString("Show cache status"),
	RunE: func(cmd *cobra.Command, args []string) error {
		cacheDir := filepath.Join(os.Getenv("HOME"), ".cache", "commit-gen")
		sessionCache := cache.GetCache(24*time.Hour, cacheDir)

		total, valid, err := sessionCache.Status()
		if err != nil {
			color.Red("Error: %v", err)
			return err
		}

		color.Cyan("Cache Status:")
		fmt.Printf("  Total entries: %d\n", total)
		fmt.Printf("  Valid entries: %d\n", valid)
		fmt.Printf("  Location: %s\n", cacheDir)

		return nil
	},
}

var cacheClearCmd = &cobra.Command{
	Use:   color.CyanString("clear"),
	Short: color.GreenString("Clear the session cache"),
	RunE: func(cmd *cobra.Command, args []string) error {
		cacheDir := filepath.Join(os.Getenv("HOME"), ".cache", "commit-gen")
		sessionCache := cache.GetCache(24*time.Hour, cacheDir)

		if err := sessionCache.Clear(); err != nil {
			color.Red("Error: %v", err)
			return err
		}

		color.Green("✓ Cache cleared successfully")
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   color.CyanString("version"),
	Short: color.GreenString("Show version information"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("commit-gen version %s\n", version)
	},
}

var healthCmd = &cobra.Command{
	Use:   color.CyanString("health"),
	Short: color.GreenString("Check if the OpenCode backend is available"),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()

		color.Cyan("Commit-gen:")
		fmt.Printf("  Version: %s\n", version)

		if _, err := os.Stat(filepath.Join(os.Getenv("HOME"), ".config", "commit-gen")); err == nil {
			color.Cyan("Configuration file:")
			fmt.Printf("  Location: %s\n", filepath.Join(os.Getenv("HOME"), ".config", "commit-gen"))
			fmt.Printf("  Exists: true\n")
		} else {
			color.Cyan("Configuration file:")
			fmt.Printf("  Location: %s\n", filepath.Join(os.Getenv("HOME"), ".config", "commit-gen"))
			fmt.Printf("  Exists: false\n")
		}

		color.Cyan("Configuration:")
		fmt.Printf("  Mode: %s\n", cfg.OpenCode.Mode)
		fmt.Printf("  Host: %s\n", cfg.OpenCode.Host)
		fmt.Printf("  Port: %d\n", cfg.OpenCode.Port)
		fmt.Printf("  Timeout: %ds\n", cfg.OpenCode.Timeout)
		fmt.Printf("  Cache: %v\n", cfg.Cache.Enabled)
		fmt.Printf("  Max Diff Size: %d bytes\n", cfg.Git.MaxDiffSize)

		color.Cyan("OpenCode Backend Check:")

		if cfg.OpenCode.Mode == "server" {
			client := opencode.NewClient(cfg.OpenCode.Host, cfg.OpenCode.Port, cfg.OpenCode.Timeout)

			healthy, err := client.CheckHealth()
			if err != nil {
				color.Red("✗ OpenCode server is not running")
				return err
			}

			if healthy {
				color.Green("✓ OpenCode server is running")
			} else {
				color.Red("✗ OpenCode server is not running")
			}
		} else {
			runner := opencode.NewRunner(cfg.OpenCode.Timeout)
			available, err := runner.CheckAvailable()
			if err != nil || !available {
				color.Red("✗ opencode binary not found in PATH")
				return err
			}
			color.Green("✓ opencode binary is available (run mode)")
		}

		return nil
	},
}

var initConfigCmd = &cobra.Command{
	Use:   color.CyanString("init"),
	Short: color.GreenString("Initialize the configuration file"),
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(filepath.Join(os.Getenv("HOME"), ".config", "commit-gen")); err == nil {
			color.Red("Error: configuration file already exists")
			return
		}
		if err := config.Initialize(""); err != nil {
			color.Red("Error: %v", err)
			return
		}
		color.Green("✓ Configuration file initialized successfully")
		fmt.Println("Now you can use: git commit -m \"\"")
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
	rootCmd.AddCommand(initConfigCmd)
	rootCmd.AddCommand(healthCmd)

	cacheCmd.AddCommand(cacheStatusCmd)
	cacheCmd.AddCommand(cacheClearCmd)
	rootCmd.AddCommand(cacheCmd)

	generateCmd.Flags().StringP("style", "s", "conventional", "Commit message style (conventional, imperative, detailed)")
	generateCmd.Flags().StringP("mode", "m", "", "Operation mode: 'run' (default) or 'server'")
	generateCmd.Flags().Bool("dry-run", false, "Show message without writing to git")
	generateCmd.Flags().Bool("hook", false, "Internal flag for git hook usage")
	generateCmd.Flags().Bool("ignore-server-check", false, "Skip checking if OpenCode backend is available")

	previewCmd.Flags().StringP("mode", "m", "", "Operation mode: 'run' (default) or 'server'")
	previewCmd.Flags().Bool("ignore-server-check", false, "Skip checking if OpenCode backend is available")
}

/**
 * checkBackendAvailability checks if the appropriate backend is available based on mode.
 *
 * @param cfg - The application configuration
 * @param ignoreCheck - If true, skip the check
 * @returns An error if the backend is not available
 */
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

/**
 * checkOpenCodeRunner verifies that the opencode binary is available in PATH.
 *
 * @returns An error if the binary is not found
 */
func checkOpenCodeRunner() error {
	runner := opencode.NewRunner(10)
	available, err := runner.CheckAvailable()
	if err != nil || !available {
		return fmt.Errorf("opencode binary not found in PATH. Please install opencode first")
	}
	return nil
}

/**
 * checkOpenCodeHealth checks if the OpenCode server is running and starts it if needed.
 *
 * @param cfg - The application configuration
 * @returns An error if the server is not running and cannot be started
 */
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

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

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

/**
 * initConfig initializes the configuration at startup.
 */
func initConfig() {
	_ = config.Initialize(cfgFile)
}

/**
 * main is the entry point for the CLI application.
 */
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
