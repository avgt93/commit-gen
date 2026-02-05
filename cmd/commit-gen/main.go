// Package main is the CLI entry point for commit-gen.
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
	Use:   "commit-gen",
	Short: "Generate commit messages using OpenCode AI",
	Long: `commit-gen is a CLI tool that generates descriptive commit messages
based on your staged git changes using OpenCode's AI capabilities.

Simply run 'git commit -m ""' and it will fill in the message for you!`,
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
	Use:   "generate",
	Short: "Generate a commit message from staged changes",
	Long: `Generate a commit message from your currently staged git changes.
The message will be generated using OpenCode's AI based on the diff.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		ignoreCheck, _ := cmd.Flags().GetBool("ignore-server-check")
		if err := checkOpenCodeHealth(cfg, ignoreCheck); err != nil {
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
	Use:   "install",
	Short: "Install git hook for automatic commit message generation",
	Long: `Installs a prepare-commit-msg git hook in the current repository.
This allows automatic commit message generation when running 'git commit -m ""'.`,
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
	Use:   "uninstall",
	Short: "Remove the git hook",
	Long:  `Removes the prepare-commit-msg git hook from the current repository.`,
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
	Use:   "config",
	Short: "Manage configuration",
	Long:  `View and modify commit-gen configuration.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()

		color.Cyan("OpenCode Configuration:")
		fmt.Printf("  Host: %s\n", cfg.OpenCode.Host)
		fmt.Printf("  Port: %d\n", cfg.OpenCode.Port)
		fmt.Printf("  Timeout: %ds\n", cfg.OpenCode.Timeout)

		color.Cyan("\nGeneration Configuration:")
		fmt.Printf("  Style: %s\n", cfg.Generation.Style)
		fmt.Printf("  Provider: %s\n", cfg.Generation.Model.Provider)
		fmt.Printf("  Model: %s\n", cfg.Generation.Model.ModelID)

		color.Cyan("\nCache Configuration:")
		fmt.Printf("  Enabled: %v\n", cfg.Cache.Enabled)
		fmt.Printf("  TTL: %s\n", cfg.Cache.TTL)

		color.Cyan("\nGit Configuration:")
		fmt.Printf("  Editor: %s\n", cfg.Git.Editor)
		fmt.Printf("  Staged Only: %v\n", cfg.Git.StagedOnly)

		return nil
	},
}

var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview changes and generated commit message",
	Long:  `Shows your staged changes and what commit message would be generated.`,
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
		ignoreCheck, _ := cmd.Flags().GetBool("ignore-server-check")
		if err := checkOpenCodeHealth(cfg, ignoreCheck); err != nil {
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
	Use:   "cache",
	Short: "Manage session cache",
	Long:  `View and manage the OpenCode session cache.`,
}

var cacheStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show cache status",
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
	Use:   "clear",
	Short: "Clear the session cache",
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
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("commit-gen version %s\n", version)
	},
}

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check if the OpenCode server is running",
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

		fmt.Printf("  Host: %s\n", cfg.OpenCode.Host)
		fmt.Printf("  Port: %d\n", cfg.OpenCode.Port)
		fmt.Printf("  Timeout: %ds\n", cfg.OpenCode.Timeout)
		fmt.Printf("  Cache: %v\n", cfg.Cache.Enabled)

		color.Cyan("OpenCode Health Check:")
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

		return nil
	},
}

var initConfigCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the configuration file",
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
	generateCmd.Flags().Bool("dry-run", false, "Show message without writing to git")
	generateCmd.Flags().Bool("hook", false, "Internal flag for git hook usage")
	generateCmd.Flags().Bool("ignore-server-check", false, "Skip checking if OpenCode server is running")

	previewCmd.Flags().Bool("ignore-server-check", false, "Skip checking if OpenCode server is running")
}

func checkOpenCodeHealth(cfg *config.Config, ignoreCheck bool) error {
	if ignoreCheck {
		return nil
	}

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

func initConfig() {
	_ = config.Initialize(cfgFile)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
