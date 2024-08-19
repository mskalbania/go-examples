package logging

import (
	"context"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"
)

func RunSlogExample() {
	textHandlerWithOpts()
	jsonHandlerUnstructured()
	jsonHandlerStructured()
	jsonHandlerLoggingStruct()
	dynamicLoggingBasedOnFlag()
}

func textHandlerWithOpts() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,            //will add file and line
		Level:     slog.LevelError, //this is how to control log level
	}))
	logger.Info("info")
	logger.Error("error")
	logger.Warn("warn")
	logger.Debug("trace")
}

func jsonHandlerUnstructured() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("request", "method", "GET", "path", "/", "status", 200)
}

func jsonHandlerStructured() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	version, _ := debug.ReadBuildInfo()
	child := logger.With(slog.Group("build",
		"go_version", version.GoVersion,
		"main", version.Main.Path,
	))
	//now all child entries will include build info
	child.Info("test")

	logger.Info(
		"request",
		//without enforcing types
		slog.String("method", "GET"),
		slog.String("path", "/"),
		slog.Int("status", 200),
	)
	logger.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"request",
		//now types are enforced to be Attr
		slog.String("method", "GET"),
		slog.String("path", "/"),
		slog.Int("status", 200),
		//and grouping
		slog.Group("headers",
			slog.String("content-type", "application/json"),
			slog.String("user-agent", "curl/7.64.1"),
		),
	)
}

func jsonHandlerLoggingStruct() {
	u := user{
		username: "admin",
		password: "secret",
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("user data", "user", u)
}

func dynamicLoggingBasedOnFlag() {
	pflag.StringP("log", "l", "INFO", "log level")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	//parse pflag and set correct logging level
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: stringToSlogLevel(viper.GetString("log")),
	}))
	l.Info("info")
	l.Debug("debug")

}

type user struct {
	username string
	password string
}

func (u user) LogValue() slog.Value { //now logging will hide sensitive data
	return slog.GroupValue(
		slog.String("username", u.username),
	)
}

func stringToSlogLevel(levelStr string) slog.Leveler {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
