package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"
)

var settings = Settings{}

type DaemonConf struct {
	PID             string        // Where to place the pidfile (not yet supported)
	KeyPollInterval time.Duration // Omit to use fs watcher (not yet supported)
}

type Settings struct {
	Logs string // Where to store logs

	DaemonMode bool       // Run in daemon mode
	Daemon     DaemonConf // Daemon configuration

	MinUptime time.Duration // Duration after which to activate
	MaxUptime time.Duration // Duration after which to force a shutdown
	Locks     string        // Locks root (either a single file or directory)
	LockAge   time.Duration // Max age for a lock
}

func configure() {
	daemonConfig := flag.String("daemon", "", "JSON configuration file to load. Overrides flags")
	minUptime := flag.Duration("minuptime", 0, "will not exit 0 before uptime >= <minuptime> (depreciated)")
	locks := flag.String("lock", "/tmp/lockfiles/", "where to look for lockfiles (depreciated)")
	lockDur := flag.Duration("duration", time.Minute*10, "duration for which lock files are considered valid (depreciated)")
	logDir := flag.String("logs", "/var/log/dozy", "Where to place log files (depreciated)")
	_ = flag.Duration("sleep", 0, "duration to sleep at the end of script before exit 0 (not used anymore)")

	flag.Parse()

	if len(*daemonConfig) > 0 {
		configFile, err := os.Open(*daemonConfig)
		if err != nil {
			panic(err)
		}
		parser := json.NewDecoder(configFile)
		if err = parser.Decode(&settings); err != nil {
			panic(err)
		}
		if len(settings.Daemon.PID) == 0 {
			settings.DaemonMode = true
		}
	} else {
		settings.MinUptime = *minUptime
		settings.Locks = *locks
		settings.LockAge = *lockDur
		settings.Logs = *logDir
	}
}
