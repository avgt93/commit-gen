package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a commit message from staged changes",
	Long: `Generate a commit message from your currently staged git changes.
The message will be generated using OpenCode's AI based on the diff.`,
	RunE: runGenerate,
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install git hook for automatic commit message generation",
	Long: `Installs a prepare-commit-msg git hook in the current repository.
This allows automatic commit message generation when running 'git commit -m ""'.`,
	RunE: runInstall,
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove the git hook",
	Long:  `Removes the prepare-commit-msg git hook from the current repository.`,
	RunE:  runUninstall,
}

var reinstallCmd = &cobra.Command{
	Use:   "reinstall",
	Short: "Reinstall the git hook",
	Long:  `Reinstalls the prepare-commit-msg git hook in the current repository.`,
	RunE:  runReinstall,
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `View and modify commit-gen configuration.`,
	RunE:  runConfig,
}

var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview changes and generated commit message",
	Long:  `Shows your staged changes and what commit message would be generated.`,
	RunE:  runPreview,
}

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage session cache",
	Long:  `View and manage the OpenCode session cache.`,
}

var cacheStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show cache status",
	RunE:  runCacheStatus,
}

var cacheClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the session cache",
	RunE:  runCacheClear,
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
	Short: "Check if the OpenCode backend is available",
	RunE:  runHealth,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the configuration file",
	Long: `Creates a configuration file at ~/.config/commit-gen/config.yaml
with default settings. Run this command once to set up commit-gen.`,
	Run: runInit,
}

// runGenerate generates a commit message from staged changes.
func runGenerate(cmd *cobra.Command, args []string) error {
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
	noConfirm, _ := cmd.Flags().GetBool("no-confirm")

	if isHook {
		fmt.Println(message)
		return nil
	}

	if dryRun {
		fmt.Println(message)
		return nil
	}

	shouldConfirm := cfg.Generation.Confirm && !noConfirm

	if shouldConfirm {
		message, err = confirmMessage(message, cfg)
		if err != nil {
			return err
		}
		if message == "" {
			color.Yellow("Commit cancelled")
			return nil
		}
	}

	if err := git.WriteCommitMessage(message); err != nil {
		return fmt.Errorf("failed to write commit message: %w", err)
	}
	color.Green("✓ Commit message generated:")
	fmt.Printf("  %s\n", message)

	return nil
}

// confirmMessage prompts the user to confirm, edit, or cancel the message.
// Returns the final message or empty string if cancelled.
func confirmMessage(message string, cfg *config.Config) (string, error) {
	color.Cyan("Generated commit message:")
	fmt.Printf("  %s\n\n", message)

	for {
		color.White("[y] Accept  [e] Edit  [r] Regenerate  [c] Cancel")
		fmt.Print("Choice: ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}

		choice := strings.ToLower(strings.TrimSpace(input))

		switch choice {
		case "y", "yes", "":
			return message, nil

		case "e", "edit":
			edited, err := editMessage(message, cfg)
			if err != nil {
				color.Red("Error editing message: %v", err)
				continue
			}
			return edited, nil

		case "r", "regenerate":
			return "", fmt.Errorf("regenerate requested")

		case "c", "cancel", "n", "no":
			return "", nil

		default:
			color.Yellow("Invalid choice. Please enter y, e, r, or c.")
		}
	}
}

// editMessage opens the user's editor to edit the commit message.
func editMessage(message string, cfg *config.Config) (string, error) {
	tmpFile, err := os.CreateTemp("", "commit-msg-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	if _, err := tmpFile.WriteString(message); err != nil {
		_ = tmpFile.Close()
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}
	_ = tmpFile.Close()

	editor := cfg.Git.Editor
	if editor == "" || editor == "cat" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor failed: %w", err)
	}

	edited, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("failed to read edited message: %w", err)
	}

	return strings.TrimSpace(string(edited)), nil
}

// runInstall installs the git hook.
func runInstall(cmd *cobra.Command, args []string) error {
	if err := hook.Install(); err != nil {
		color.Red("Error: %v", err)
		return err
	}
	color.Green("✓ Git hook installed successfully")
	fmt.Println("Now you can use: git commit")
	fmt.Println("The generated message will open in your editor for confirmation.")
	return nil
}

// runUninstall removes the git hook.
func runUninstall(cmd *cobra.Command, args []string) error {
	if err := hook.Uninstall(); err != nil {
		color.Red("Error: %v", err)
		return err
	}
	color.Green("✓ Git hook removed successfully")
	return nil
}

// runReinstall reinstalls the git hook.
func runReinstall(cmd *cobra.Command, args []string) error {
	if err := hook.Uninstall(); err != nil {
		color.Red("Error: %v", err)
		return err
	}
	return runInstall(cmd, args)
}

// runConfig displays the current configuration.
func runConfig(cmd *cobra.Command, args []string) error {
	cfg := config.Get()

	color.Cyan("OpenCode Configuration:")
	fmt.Printf("  Mode: %s\n", cfg.OpenCode.Mode)
	fmt.Printf("  Host: %s (server mode only)\n", cfg.OpenCode.Host)
	fmt.Printf("  Port: %d (server mode only)\n", cfg.OpenCode.Port)
	fmt.Printf("  Timeout: %ds\n", cfg.OpenCode.Timeout)

	color.Cyan("\nGeneration Configuration:")
	fmt.Printf("  Style: %s\n", cfg.Generation.Style)
	fmt.Printf("  Confirm: %v\n", cfg.Generation.Confirm)
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
}

// runPreview shows staged changes and the generated commit message.
func runPreview(cmd *cobra.Command, args []string) error {
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
}

// runCacheStatus displays cache statistics.
func runCacheStatus(cmd *cobra.Command, args []string) error {
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
}

// runCacheClear clears all cached sessions.
func runCacheClear(cmd *cobra.Command, args []string) error {
	cacheDir := filepath.Join(os.Getenv("HOME"), ".cache", "commit-gen")
	sessionCache := cache.GetCache(24*time.Hour, cacheDir)

	if err := sessionCache.Clear(); err != nil {
		color.Red("Error: %v", err)
		return err
	}

	color.Green("✓ Cache cleared successfully")
	return nil
}

// runHealth checks if the OpenCode backend is available.
func runHealth(cmd *cobra.Command, args []string) error {
	cfg := config.Get()

	color.Cyan("Commit-gen:")
	fmt.Printf("  Version: %s\n", version)

	configPath, _ := config.GetConfigPath()
	color.Cyan("Configuration file:")
	fmt.Printf("  Location: %s\n", configPath)
	fmt.Printf("  Exists: %v\n", config.ConfigExists())

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
}

// runInit initializes the configuration file.
func runInit(cmd *cobra.Command, args []string) {
	if config.ConfigExists() {
		configPath, _ := config.GetConfigPath()
		color.Yellow("Configuration file already exists at: %s", configPath)
		fmt.Println("Use 'commit-gen config' to view current settings.")
		return
	}

	if err := config.Initialize(""); err != nil {
		color.Red("Error initializing config: %v", err)
		return
	}

	if err := config.CreateConfig(); err != nil {
		color.Red("Error creating config file: %v", err)
		return
	}

	configPath, _ := config.GetConfigPath()
	color.Green("✓ Configuration file created successfully")
	fmt.Printf("  Location: %s\n", configPath)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Edit the config file to customize settings (optional)")
	fmt.Println("  2. Run 'commit-gen install' in your git repository")
	fmt.Println("  3. Use 'git commit' to generate commit messages")
}
