package sandbox

import (
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

type Permission int

const (
	ReadMeta  Permission = 1 << 0
	ReadData  Permission = 1 << 1
	WriteMeta Permission = 1 << 2
	WriteData Permission = 1 << 3
)

func (p Permission) MarshalYAML() (any, error) {
	var perms []string
	if p&ReadMeta != 0 {
		perms = append(perms, "read-meta")
	}
	if p&ReadData != 0 {
		perms = append(perms, "read-data")
	}
	if p&WriteMeta != 0 {
		perms = append(perms, "write-meta")
	}
	if p&WriteData != 0 {
		perms = append(perms, "write-data")
	}
	return perms, nil
}

func (p *Permission) UnmarshalYAML(value *yaml.Node) error {
	var perms []string
	if err := value.Decode(&perms); err != nil {
		return err
	}

	var result Permission
	for _, perm := range perms {
		switch perm {
		case "read-meta":
			result |= ReadMeta
		case "read-data":
			result |= ReadData
		case "write-meta":
			result |= WriteMeta
		case "write-data":
			result |= WriteData
		default:
			return fmt.Errorf("unknown permission: %s", perm)
		}
	}
	*p = result
	return nil
}

type FolderPermission struct {
	Path       string     `yaml:"path"`
	Permission Permission `yaml:"permissions"`
}

type Config struct {
	Command        string     `yaml:"command"`
	IncludeParents Permission `yaml:"include_parents"`

	FoldersInParents []FolderPermission `yaml:"folders_in_parents"`

	ExportEnvironment []string `yaml:"export_env"`

	// Working config
	cwd         string
	home        string
	configDir   string
	configFile  string
	commandName string
}

func FromCommand(command string) (*Config, error) {
	commandName := filepath.Base(command)
	home, err := getHomeDir()
	if err != nil {
		slog.Error("Unable to resolve home directory")
		return nil, fmt.Errorf("resolving home directory: %w", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("resolving working directory: %w", err)
	}
	cfg := Config{
		Command:     command,
		home:        home,
		cwd:         filepath.Clean(cwd),
		configDir:   filepath.Join(home, ".config", "sb"),
		configFile:  filepath.Join(home, ".config", "sb", commandName+".yaml"),
		commandName: commandName,
	}
	slog.Debug("loading config", "path", cfg.configFile)
	if err := fromConfigFile(&cfg); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("loading config %s: %w", cfg.configFile, err)
		}
		slog.Debug("no config file found, using defaults", "path", cfg.configFile)
	} else {
		slog.Debug("config loaded", "command", cfg.Command, "path", cfg.configFile)
	}
	return &cfg, nil
}

func fromConfigFile(cfg *Config) error {
	fd, err := os.Open(cfg.configFile)
	if err != nil {
		return fmt.Errorf("opening config file %s: %w", cfg.configFile, err)
	}
	// Read-only file: a close error cannot lose data.
	defer func() { _ = fd.Close() }()
	decoder := yaml.NewDecoder(fd)
	if err := decoder.Decode(&cfg); err != nil {
		return fmt.Errorf("decoding %s: %w", cfg.configFile, err)
	}
	return nil
}

func getHomeDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		home, ok := os.LookupEnv("HOME")
		if !ok {
			return "", fmt.Errorf("cannot determine home directory ($HOME not set): %w", err)
		}
		slog.Warn("user.Current() failed, falling back to $HOME environment variable", "error", err)
		return home, nil
	}
	return u.HomeDir, nil
}
