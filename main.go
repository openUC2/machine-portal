package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/carlmjohnson/versioninfo"
	"github.com/labstack/echo/v4"
	"github.com/urfave/cli/v3"

	"github.com/openUC2/machine-portal/internal/app/server"
	"github.com/openUC2/machine-portal/internal/app/server/conf"
)

func main() {
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

const (
	defaultPort            = 3001
	defaultShutdownTimeout = 5 * time.Second
)

var cmd = &cli.Command{
	Name:    "machine-portal",
	Version: toolVersion,
	Usage:   "Provides a landing page",
	Action:  serverMain,
	Flags: []cli.Flag{
		// HTTP server
		&cli.IntFlag{
			Name:    "http-port",
			Value:   defaultPort,
			Usage:   "port for HTTP server",
			Sources: cli.EnvVars("HTTP_PORT"),
		},
		&cli.StringFlag{
			Name:    "http-base-path",
			Value:   "/",
			Usage:   "base path for HTTP routes",
			Sources: cli.EnvVars("HTTP_BASEPATH"),
		},
		&cli.IntFlag{
			Name:    "http-gzip-level",
			Value:   1,
			Usage:   "port for HTTP server",
			Sources: cli.EnvVars("HTTP_GZIPLEVEL"),
		},
		&cli.DurationFlag{
			Name:    "http-shutdown-timeout",
			Value:   defaultShutdownTimeout,
			Usage:   "timeout for graceful shutdown before hard shutdown",
			Sources: cli.EnvVars("SHUTDOWNTIMEOUT"),
		},
	},
}

func serverMain(ctx context.Context, cmd *cli.Command) error {
	e := echo.New()

	// Get config
	config, err := conf.GetConfig()
	if err != nil {
		return err
	}
	config.HTTP.Port = cmd.Int("http-port")
	config.HTTP.BasePath = cmd.String("http-base-path")
	config.HTTP.GzipLevel = cmd.Int("http-gzip-level")

	// Prepare server
	s, err := server.New(config, e.Logger)
	if err != nil {
		return err
	}
	if err = s.Register(e); err != nil {
		return err
	}

	// Run server
	ctxRun, cancelRun := signal.NotifyContext(
		ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT,
	)
	go func() {
		if err = s.Run(e); err != nil {
			e.Logger.Error(err)
		}
		cancelRun()
	}()
	<-ctxRun.Done()
	cancelRun()

	// Shut down server
	shutdownTimeout := cmd.Duration("http-shutdown-timeout")
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelShutdown()
	e.Logger.Infof("attempting to shut down gracefully within %.1f sec", shutdownTimeout.Seconds())
	if err := s.Shutdown(ctxShutdown, e); err != nil {
		e.Logger.Warn("forcibly closing http server due to failure of graceful shutdown")
		if closeErr := s.Close(e); closeErr != nil {
			return closeErr
		}
	}
	e.Logger.Info("finished shutdown")
	return nil
}

// Versioning

const (
	// fallbackVersion is the version reported which the Forklift tool reports itself as if its actual
	// version is unknown.
	fallbackVersion = "dev"
)

var (
	toolVersion = determineVersion(buildSummary, fallbackVersion)
	// buildSummary should be overridden by ldflags, such as with GoReleaser's "Summary".
	buildSummary = ""
)

// determineVersion returns either a semver, a pseudoversion, or a Git hash based on information
// available from Go's `debug.ReadBuildInfo()`.
func determineVersion(override, fallback string) string {
	if override != "" {
		return override
	}

	const dirtySuffix = "-dirty"
	// Determine any version tags, if available
	if info, ok := debug.ReadBuildInfo(); ok &&
		info.Main.Version != "" && info.Main.Version != "(devel)" {
		v := info.Main.Version
		if versioninfo.DirtyBuild {
			v += dirtySuffix
		}
		return v
	}
	if v := versioninfo.Version; v != "unknown" && v != "(devel)" {
		if versioninfo.DirtyBuild {
			v += dirtySuffix
		}
		return v
	}

	// Fall back to whatever is available
	if r := versioninfo.Revision; r != "unknown" && r != "" {
		if versioninfo.DirtyBuild {
			r += dirtySuffix
		}
		return r
	}
	return fallback
}
