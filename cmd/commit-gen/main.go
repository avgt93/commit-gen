package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/yourusername/commit-gen/internal/cache"
	"github.com/yourusername/commit-gen/internal/config"
	"github.com/yourusername/commit-gen/internal/generator"
	"github.com/yourusername/commit-gen/internal/git"
	"github.com/yourusername/commit-gen/internal/hook"
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
		// Default to generate command if no args
		if len(args) == 0 {
			generateCmd.Run(cmd, args)
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
		cacheDir := filepath.Join(os.Getenv("HOME"), ".cache", "commit-gen")
		sessionCache := cache.GetCache(24*time.Hour, cacheDir)

		gen := generator.NewGenerator(cfg, sessionCache)

		message, err := gen.Generate()
		if err != nil {
			color.Red("Error: %v", err)

			// Prompt user to start OpenCode if needed
			if err.Error() == fmt.Sprintf("opencode server is not running at %s:%d\n\nPlease start it with: opencode serve", cfg.OpenCode.Host, cfg.OpenCode.Port) {
				fmt.Println("\nTo fix this, run:")
				color.Cyan("  opencode serve")
				fmt.Println("\nIn another terminal, then try again.")
				return err
			}

			return err
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		isHook, _ := cmd.Flags().GetBool("hook")

		if dryRun || isHook {
			// Just output the message
			fmt.Println(message)
		} else {
			// Write to commit message file
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
		if err := hook.Install(); err != nil {
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

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/commit-gen/config.yaml)")

	// Add commands
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(previewCmd)
	rootCmd.AddCommand(versionCmd)

	// Cache subcommands
	cacheCmd.AddCommand(cacheStatusCmd)
	cacheCmd.AddCommand(cacheClearCmd)
	rootCmd.AddCommand(cacheCmd)

	// Generate command flags
	generateCmd.Flags().StringP("style", "s", "conventional", "Commit message style (conventional, imperative, detailed)")
	generateCmd.Flags().Bool("dry-run", false, "Show message without writing to git")
	generateCmd.Flags().Bool("hook", false, "Internal flag for git hook usage")
}

func initConfig() {
	config.Initialize(cfgFile)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
