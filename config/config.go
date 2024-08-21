package config

import (
	"log/slog"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type bigipConfig struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	BasicAuth bool   `yaml:"basic_auth"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
}

type exporterConfig struct {
	BindAddress string `yaml:"bind_address"`
	BindPort    int    `yaml:"bind_port"`
	Partitions  string `yaml:"partitions"`
	Config      string `yaml:"config"`
	Namespace   string `yaml:"namespace"`
	LogLevel    string `yaml:"log_level"`
}

// Config is a container for settings modifiable by the user.
type Config struct {
	Bigip    bigipConfig    `yaml:"bigip"`
	Exporter exporterConfig `yaml:"exporter"`
}

func init() {
	registerFlags()
	bindFlags()
	bindEnvs()
	flag.Parse()

	if viper.GetString("exporter.config") != "" {
		readConfigFile(viper.GetString("exporter.config"))
	}

	logLevel := viper.GetString("exporter.log_level")
	if levelCode, validLevel := parseLevel(logLevel); validLevel {
		slog.SetLogLoggerLevel(levelCode)
	} else {
		slog.Warn("Invalid log level: using info")
	}
}

func parseLevel(level string) (slog.Level, bool) {
	switch level {
	case "debug":
		return slog.LevelDebug, true
	case "info":
		return slog.LevelInfo, true
	case "warn":
		return slog.LevelWarn, true
	case "error":
		return slog.LevelError, true
	default:
		return slog.LevelInfo, false
	}
}

func registerFlags() {
	flag.Bool("bigip.basic_auth", false, "Use HTTP Basic authentication")
	flag.String("bigip.host", "localhost", "The host on which f5 resides")
	flag.Int("bigip.port", 443, "The port which f5 listens to")
	flag.String("bigip.username", "user", "Username")
	flag.String("bigip.password", "pass", "Password")
	flag.String("exporter.bind_address", "localhost", "Exporter bind address")
	flag.Int("exporter.bind_port", 9142, "Exporter bind port")
	flag.String("exporter.partitions", "", "A comma separated list of partitions which to export. (default: all)")
	flag.String("exporter.config", "", "bigip_exporter configuration file name.")
	flag.String("exporter.namespace", "bigip", "bigip_exporter namespace.")
	flag.String("exporter.log_level", "info", "Available options are trace, debug, info, warning, error and critical")
}

func bindFlags() {
	flag.VisitAll(func(f *flag.Flag) {
		err := viper.BindPFlag(f.Name, f)
		if err != nil {
			slog.Warn("Failed to bind flag", "error", err)
		}
	})
}

func bindEnvs() {
	viper.SetEnvPrefix("be")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	flag.VisitAll(func(f *flag.Flag) {
		err := viper.BindEnv(f.Name)
		if err != nil {
			varName := "BE_" + strings.ToUpper(strings.ReplaceAll(f.Name, ".", "_"))
			slog.Warn("Failed to bind environment variable", "var", varName, "error", err)
		}
	})
}

func readConfigFile(fileName string) {
	if file, err := os.Open(fileName); err != nil {
		slog.Warn("Failed to open configuration file", "error", err)
	} else {
		viper.SetConfigType("yaml")
		if err = viper.ReadConfig(file); err != nil {
			slog.Warn("Failed to read configuration file", "error", err)
		}
	}
}

// GetConfig returns an instance of Config containing the resulting parameters
// to the program.
func GetConfig() *Config {
	return &Config{
		Bigip: bigipConfig{
			Username:  viper.GetString("bigip.username"),
			Password:  viper.GetString("bigip.password"),
			BasicAuth: viper.GetBool("bigip.basic_auth"),
			Host:      viper.GetString("bigip.host"),
			Port:      viper.GetInt("bigip.port"),
		},
		Exporter: exporterConfig{
			BindAddress: viper.GetString("exporter.bind_address"),
			BindPort:    viper.GetInt("exporter.bind_port"),
			Partitions:  viper.GetString("exporter.partitions"),
			Config:      viper.GetString("exporter.config"),
			Namespace:   viper.GetString("exporter.namespace"),
			LogLevel:    viper.GetString("exporter.log_level"),
		},
	}
}
