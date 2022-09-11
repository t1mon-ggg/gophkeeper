package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env"
	"github.com/creasty/defaults"

	"github.com/t1mon-ggg/gophkeeper/pkg/logging"
	"github.com/t1mon-ggg/gophkeeper/pkg/logging/zerolog"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

//Config - struct for handling configuration
type Config struct {
	WebBind  string         `env:"WEB_ADDRESS" json:"web_address" default:":8443"`
	GRPCBind string         `env:"GRPC_ADDRESS" json:"grpc_address" default:":3200"`
	DSN      string         `env:"DSN_ADDRESS" json:"postgre_dsn" default:"postgresql://gophkeeper:gophkeeper@127.0.0.1:5432/gophkeeper"`
	LogLevel uint8          `json:"LogLevel,omitempty" default:"-"`
	file     string         `json:"-" default:"config.json"`
	logger   logging.Logger `default:"-"`
}

//NewConfig - config initialization
func New() *Config {
	c := new(Config)
	c.logger = zerolog.New()
	c.SetDefaults()
	c.SetByFile()
	c.SetByFlags()
	c.SetByEnv()
	c.log().Infof("resulted configuration:\"\n\t\t\t\tWebBind:\t%v\n\t\t\t\tGRPCBind:\t%v\n\t\t\t\tDSN:\t\t%v\n\t\t\t\tLogging Level:\t%v\n\t\t\t", nil, c.WebBind, c.GRPCBind, c.DSN, c.LogLevel)
	return c
}

// log - returns default application Logger
func (c *Config) log() logging.Logger {
	return c.logger.WithPrefix("configuration")
}

//SetDefaults - set configuration values to default state
func (c *Config) SetDefaults() *Config {
	err := defaults.Set(c)
	if err != nil {
		c.log().Fatal(err)
	}
	c.LogLevel = uint8(logging.DebugLevel)
	return c
}

//SetByEnv - set configuration values from evironment
//  WEB_ADDRESS - evironment for web api instance binding. Default is ":8443"
//  GRPC_ADDRESS - evironment for grpc instance binding. Default is ":3200"
//  DSN_ADDRESS - evironment for postgres connection string. Default is "postgresql://gophkeeper:gophkeeper@127.0.0.1:5432/gophkeeper"
func (c *Config) SetByEnv() *Config {
	cc := Config{}
	err := env.Parse(&cc)
	if err != nil {
		c.log().Warnf("read environment failed with error: %v", err)
		return c
	}
	if cc.WebBind != "" {
		c.WebBind = cc.WebBind
	}
	if cc.DSN != "" {
		c.DSN = cc.DSN
	}
	if cc.GRPCBind != "" {
		c.GRPCBind = cc.GRPCBind
	}
	return c
}

//Configuring flags
var (
	webBindFlag  = flag.String("http", "", "Set up web address binding.\nExample: -http=\":8443\"")
	grpcBindFlag = flag.String("grpc", "", "Set up grpc address binding.\nExample: -grpc=\":3200\"")
	logLevelFlag = flag.String("LogLevel", "debug", "Set up logging level.\nExample: -LogLevel=info\nAvailible level are trace, debug, info, warn, error, fatal, panic.")
	dsn          = flag.String("dsn", "", "Set up postgre dsn connection string\nExample: -dsn=postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]")
	config       = flag.String("config", "", "Set up configuration file path.\n Example: -config=config.json")
)

//SetByFlags - set configuration values from cli flags
//	-debug
//	-http
//	-grpc
//	-LogLevel
//  -version
//  -dsn
func (c *Config) SetByFlags() *Config {
	flag.Parse()
	if flag.Parsed() {
		if isFlagPassed("version") {
			fmt.Printf("%s\n%s\n%s", buildVersion, buildDate, buildCommit)
			os.Exit(0)
		}
		if webBindFlag != nil && *webBindFlag != "" {
			c.WebBind = *webBindFlag
		}
		if grpcBindFlag != nil && *grpcBindFlag != "" {
			c.GRPCBind = *grpcBindFlag
		}
		if dsn != nil && *dsn != "" {
			c.DSN = *dsn
		}
		if config != nil && *config != "" {
			c.file = *config
		}
		if logLevelFlag != nil {
			if isFlagPassed("LogLevel") {
				switch *logLevelFlag {
				case "trace":
					c.LogLevel = uint8(logging.TraceLevel)
				case "debug":
					c.LogLevel = uint8(logging.DebugLevel)
				case "info":
					c.LogLevel = uint8(logging.InfoLevel)
				case "warn":
					c.LogLevel = uint8(logging.WarnLevel)
				case "error":
					c.LogLevel = uint8(logging.ErrorLevel)
				case "fatal":
					c.LogLevel = uint8(logging.FatalLevel)
				case "panic":
					c.LogLevel = uint8(logging.PanicLevel)
				default:
					c.LogLevel = uint8(logging.DebugLevel)
				}
			}
		}
	}
	return c
}

//SetByFile - set configuration values from configuration file
func (c *Config) SetByFile() *Config {
	cfg, err := os.ReadFile(c.file)
	if err != nil {
		c.log().Warn(err, "configuration file can not be read")
		return c
	}
	cc := new(Config)
	err = json.Unmarshal(cfg, cc)
	if err != nil {
		c.log().Warn(err, "configuration can not be parsed")
		return c
	}
	cc.file = c.file
	return cc
}

//SetByFile - set configuration values from configuration file
func (c *Config) Level() logging.Level {
	return logging.Level(c.LogLevel)
}

//isFlagPassed - checking the using of the flag
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
