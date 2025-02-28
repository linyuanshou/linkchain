package config

import (
	"path/filepath"

	"github.com/lianxiangcloud/linkchain/libs/log"
)

var (
	// DefaultWalletDir wallet home dir
	DefaultWalletDir   = "wallet"
	defaultDataDir     = "data"
	defaultKeyStoreDir = "keystore"
	defaultLogDir      = "logs"
	defaultLogFileName = "wallet.log"
)

// BaseConfig define
type BaseConfig struct {
	Password       string `mapstructure:"password"`
	PasswordFile   string `mapstructure:"password_file"`
	KeystoreFile   string `mapstructure:"keystore_file"`
	KdfRounds      int    `mapstructure:"kdf_rounds"`
	Detach         bool   `mapstructure:"detach"`
	MaxConcurrency int    `mapstructure:"max_concurrency"`
	Pidfile        string `mapstructure:"pidfile"`
	LogLevel       string `mapstructure:"log_level"`
	// The root directory for all data.
	// This should be set in viper so it can unmarshal into this struct
	RootDir string `mapstructure:"home"`
	// Database backend: leveldb | memdb
	DBBackend string `mapstructure:"db_backend"`
	// Database directory
	DBPath string `mapstructure:"db_path"`
	// KeyStore directory
	KeyStorePath string `mapstructure:"keystore_dir"`
	// LogPath directory
	LogPath string `mapstructure:"log_dir"`
	LogFile string `mapstructure:"log_file"`
}

// DefaultBaseConfig return default config
func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		Password:       "",
		KdfRounds:      1,
		Detach:         false,
		MaxConcurrency: 1,
		LogLevel:       "debug",
		DBBackend:      "leveldb",
		DBPath:         defaultDataDir,
		LogPath:        defaultLogDir,
		KeyStorePath:   defaultKeyStoreDir,
	}
}

// LogDir returns the full path to the log directory
func (cfg BaseConfig) LogDir() string {
	return rootify(cfg.LogPath, cfg.RootDir)
}

// DBDir returns the full path to the database directory
func (cfg BaseConfig) DBDir() string {
	return rootify(cfg.DBPath, cfg.RootDir)
}

// KeyStoreDir returns the full path to the keystore directory
func (cfg BaseConfig) KeyStoreDir() string {
	return rootify(cfg.KeyStorePath, cfg.RootDir)
}

// helper function to make config creation independent of root dir
func rootify(path, root string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}

// DaemonConfig daemon config
type DaemonConfig struct {
	PeerRPC string `mapstructure:"peer_rpc"`
	Login   string `mapstructure:"login"`
	Trusted bool   `mapstructure:"trusted"`
	Testnet bool   `mapstructure:"testnet"`
}

// RPCConfig rpc config
type RPCConfig struct {
	IpcEndpoint  string   `mapstructure:"ipc_endpoint"`
	HTTPEndpoint string   `mapstructure:"http_endpoint"`
	HTTPModules  []string `mapstructure:"http_modules"`
	HTTPCores    []string `mapstructure:"http_cores"`
	VHosts       []string `mapstructure:"vhosts"`
	WSEndpoint   string   `mapstructure:"ws_endpoint"`
	WSModules    []string `mapstructure:"ws_modules"`
	WSOrigins    []string `mapsturcture:"ws_origins"`
	WSExposeAll  bool     `mapstructure:"ws_expose_all"`
}

// DefaultDaemonConfig returns default daemon config
func DefaultDaemonConfig() *DaemonConfig {
	return &DaemonConfig{
		PeerRPC: "http://127.0.0.1:11000",
		Login:   "",
		Trusted: true,
		Testnet: true,
	}
}

// DefaultRPCConfig returns default rpc config
func DefaultRPCConfig() *RPCConfig {
	return &RPCConfig{
		IpcEndpoint:  "wallet.ipc",
		HTTPEndpoint: "127.0.0.1:18082",
		HTTPModules:  []string{"personal", "ltk"},
		HTTPCores:    []string{"*"},
		VHosts:       []string{"*"},
		WSEndpoint:   "127.0.0.1:18083",
		WSModules:    []string{"personal", "ltk"},
		WSExposeAll:  true,
		WSOrigins:    []string{"*"},
	}
}

// DefaultRotateConfig returns default roate config
func DefaultRotateConfig() *log.RotateConfig {
	return &log.RotateConfig{
		Filename:   defaultLogFileName,
		Daily:      true,
		MaxDays:    7,
		Rotate:     true,
		RotatePerm: "0444",
		Perm:       "0664",
	}
}

func (cfg *Config) IPCFile() string {
	return rootify(cfg.RPC.IpcEndpoint, cfg.RootDir)
}

// Config config struct
type Config struct {
	// Top level options use an anonymous struct
	BaseConfig `mapstructure:",squash"`
	Daemon     *DaemonConfig     `mapstructure:"daemon"`
	RPC        *RPCConfig        `mapstructure:"rpc"`
	Log        *log.RotateConfig `mapstructure:"log"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		BaseConfig: DefaultBaseConfig(),
		Daemon:     DefaultDaemonConfig(),
		RPC:        DefaultRPCConfig(),
		Log:        DefaultRotateConfig(),
	}
}

// DefaultLogLevel returns a default log level of "info"
func DefaultLogLevel() string {
	return "debug"
}
