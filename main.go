package main

import (
	"flag"
	"strings"

	"github.com/golang/glog"
	watchdog "github.com/jchorl/watchdog/client"
	watchdog_types "github.com/jchorl/watchdog/types"

	"github.com/jchorl/nasblaze/drives"
	"github.com/jchorl/nasblaze/rclone"
)

const (
	driveSize  = "931.3G"
	mountpoint = "/sync"
	bucketName = "backup-c3bac1e7-d888-4940-8778-f03adcddfe45"
	configPath = "/home/j/nasblaze/rclone.conf"
)

var (
	hardDelete = flag.Bool("hard-delete", false, "whether to hard delete files on b2")
	dryRun     = flag.Bool("dry-run", true, "whether to run in dry run mode")
	filters    = arrayFlags{}
)

func main() {
	flag.Var(&filters, "exclude", "Exclude flags passed directly to rclone")
	flag.Parse()

	glog.Info("Starting")

	// mount the drive, if necessary
	err := drives.MountDriveBySize(driveSize, mountpoint)
	if err != nil {
		glog.Fatalf("Error mounting drive: %s", err)
	}

	err = rclone.Sync(configPath, bucketName, mountpoint, filters, *dryRun, *hardDelete)
	if err != nil {
		glog.Fatalf("Error syncing: %s", err)
	}

	wdClient := watchdog.Client{"https://watchdog.joshchorlton.com"}
	wdClient.Ping("nasblaze", watchdog_types.Watch_WEEKLY)
	glog.Info("Complete")
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
