package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nacuellar25111992/bricks-backend-poc/internal/signals"
	"github.com/nacuellar25111992/bricks-backend-poc/internal/version"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/nacuellar25111992/bricks-backend-poc/internal/api"
)

func main() {

	// flags definition.
	fs := pflag.NewFlagSet("default", pflag.ContinueOnError)
	fs.String("host", "", "Host to bind service to")
	fs.Int("port", 9898, "HTTP port to bind service to")
	fs.String("log-level", "debug", "log level debug, info, warn, error, flat or panic")
	fs.StringSlice("backend-url", []string{}, "backend service URL")
	fs.Duration("http-client-timeout", 2*time.Minute, "client timeout duration")
	fs.Duration("http-server-timeout", 30*time.Second, "server read and write timeout duration")
	fs.Duration("http-server-shutdown-timeout", 5*time.Second, "server graceful shutdown timeout duration")
	fs.String("config-path", ".", "config dir path")
	fs.String("config", "config.yaml", "config file name")
	fs.Bool("h2c", false, "allow upgrading to H2C")
	fs.Bool("unhealthy", false, "when set, healthy state is never reached")
	fs.Bool("unready", false, "when set, ready state is never reached")

	versionFlag := fs.BoolP("version", "v", false, "get version number")

	// parse flags.
	err := fs.Parse(os.Args[1:])
	switch {
	case err == pflag.ErrHelp:
		os.Exit(0)
	case err != nil:
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
		fs.PrintDefaults()
		os.Exit(2)
	case *versionFlag:
		fmt.Println(version.VERSION)
		os.Exit(0)
	}

	// bind flags and environment variables.
	viper.BindPFlags(fs)
	viper.RegisterAlias("backendUrl", "backend-url")
	hostname, _ := os.Hostname()
	viper.SetDefault("jwt-secret", "elarbolgigante")
	viper.Set("hostname", hostname)
	viper.Set("version", version.VERSION)
	viper.Set("revision", version.REVISION)
	viper.SetEnvPrefix("BRICKS_BACKEND_POC")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// load config from file.
	_, err = os.Stat(filepath.Join(viper.GetString("config-path"), viper.GetString("config")))
	if err == nil {

		configName := strings.Split(viper.GetString("config"), ".")[0]
		configPath := viper.GetString("config-path")

		viper.SetConfigType("yaml")
		viper.SetConfigName(configName)
		viper.AddConfigPath(configPath)

		err = viper.ReadInConfig()
		if err != nil {

			fmt.Printf("error reading config file. %v", err)
			os.Exit(3)
		}
	}

	// configure logging.
	logger, err := initZap(viper.GetString("log-level"))
	if err != nil {

		fmt.Printf("error reading config file. %v", err)
		os.Exit(3)
	}
	defer logger.Sync()
	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()

	// validate port.
	if _, err := strconv.Atoi(viper.GetString("port")); err != nil {
		port, _ := fs.GetInt("port")
		viper.Set("port", strconv.Itoa(port))
	}

	// load http server config.
	var srvCfg api.Config
	err = viper.Unmarshal(&srvCfg)
	if err != nil {
		logger.Panic("config unmarshal failed", zap.Error(err))
	}

	// log version and port.
	logger.Info("starting bricks-backend-poc",
		zap.String("version", viper.GetString("version")),
		zap.String("revision", viper.GetString("revision")),
		zap.String("port", srvCfg.Port),
	)

	// start HTTP server.
	stopChannel := signals.SetupSignalHandler()

	srv := api.NewServer(&srvCfg, logger)
	srv.ListenAndServe(stopChannel)
}

func initZap(logLevel string) (*zap.Logger, error) {

	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)

	switch logLevel {
	case "debug":
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case "panic":
		level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	}

	zapEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	zapConfig := zap.Config{
		Level:       level,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zapEncoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return zapConfig.Build()
}
