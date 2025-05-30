package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"gopkg.in/yaml.v2"
)

type ConfigFile struct {
	AppName   string `yaml:"appname"`
	SecretKey string `yaml:"secretkey"`
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage UploadThing configuration",
	Long:  `Configure your UploadThing credentials and settings.`,
}

var setSecretCmd = &cobra.Command{
	Use:   "set-secret <secret-key>",
	Short: "Set your UploadThing secret key",
	Long:  `Set your UploadThing secret API key for authentication.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		secretKey := args[0]
		err := setSecretKey(secretKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error setting secret key: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Secret key updated successfully!")
	},
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current UploadThing configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := showConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error showing config: %v\n", err)
			os.Exit(1)
		}
	},
}

var setConfigPathCmd = &cobra.Command{
	Use:   "set-config-path <path>",
	Short: "Set configuration file path",
	Long:  `Set a custom path for the configuration file.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configPath := args[0]
		fmt.Printf("Custom config path set to: %s\n", configPath)
		fmt.Println("Note: This feature is not yet implemented.")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setSecretCmd)
	configCmd.AddCommand(showConfigCmd)
	configCmd.AddCommand(setConfigPathCmd)
}

func setSecretKey(secretKey string) error {
	configDir, configFile, err := getConfigPaths()
	if err != nil {
		return err
	}

	if err := ensureConfigDir(configDir); err != nil {
		return err
	}

	var cfg ConfigFile
	if data, err := os.ReadFile(configFile); err == nil {
		yaml.Unmarshal(data, &cfg)
	}

	cfg.SecretKey = secretKey

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("unable to marshal config: %w", err)
	}

	err = os.WriteFile(configFile, data, 0600)
	if err != nil {
		return fmt.Errorf("unable to write config file: %w", err)
	}

	return nil
}

func showConfig() error {
	_, configFile, err := getConfigPaths()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No configuration file found.")
			fmt.Println("Use 'ut config set-secret <key>' to set up your UploadThing secret key.")
			return nil
		}
		return fmt.Errorf("unable to read config file: %w", err)
	}

	var cfg ConfigFile
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("unable to parse config file: %w", err)
	}

	fmt.Println("Current Configuration:")
	fmt.Printf("  Config file: %s\n", configFile)
	if cfg.AppName != "" {
		fmt.Printf("  App Name: %s\n", cfg.AppName)
	}
	if cfg.SecretKey != "" {
		maskedKey := maskSecretKey(cfg.SecretKey)
		fmt.Printf("  Secret Key: %s\n", maskedKey)
	} else {
		fmt.Println("  Secret Key: (not set)")
	}

	return nil
}

func getConfigPaths() (string, string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", "", fmt.Errorf("unable to find home directory: %w", err)
	}

	configDir := filepath.Join(home, ".ut-cli")
	configFile := filepath.Join(configDir, "config.yml")

	return configDir, configFile, nil
}

func ensureConfigDir(configDir string) error {
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0700)
		if err != nil {
			return fmt.Errorf("unable to create config directory: %w", err)
		}
	}
	return nil
}

func maskSecretKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

func readSecret() (string, error) {
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	fmt.Println()
	return string(bytePassword), nil
}
