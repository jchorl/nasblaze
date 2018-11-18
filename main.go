package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/jchorl/watchdog"
	"github.com/ncw/rclone/fs"

	"github.com/jchorl/nasblaze/drives"
	"github.com/jchorl/nasblaze/rclone"
)

const (
	driveSize    = "931.3G"
	mountpoint   = "/sync"
	bucketName   = "backup-c3bac1e7-d888-4940-8778-f03adcddfe45"
	configPath   = "/home/j/nasblaze/rclone.conf"
	otherLogFile = "/home/j/logs/nasblaze.OTHER"
)

var (
	hardDelete = flag.Bool("hard-delete", false, "whether to hard delete files on b2")
	dryRun     = flag.Bool("dry-run", true, "whether to run in dry run mode")
	filters    = arrayFlags{}
)

func main() {
	flag.Var(&filters, "exclude", "Exclude flags passed directly to rclone")
	flag.Parse()

	// set rclone log level
	fs.Config.LogLevel = fs.LogLevelInfo
	logFile, err := initGolangLog()
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	glog.Info("Starting")

	// mount the drive, if necessary
	err = drives.MountDriveBySize(driveSize, mountpoint)
	if err != nil {
		glog.Fatalf("Error mounting drive: %s", err)
	}

	err = rclone.Sync(configPath, bucketName, mountpoint, filters, *dryRun, *hardDelete)
	if err != nil {
		glog.Fatalf("Error syncing: %s", err)
	}

	err = drives.UnmountDriveByMountpoint(mountpoint)
	if err != nil {
		// don't fail, there is a bug where the disk unmounts
		// but reports that it did not
		glog.Infof("Error unmounting drive: %s", err)
	}

	wdClient := watchdog.Client{"https://watchdog.joshchorlton.com"}
	wdClient.Ping("nasblaze", watchdog.Watch_WEEKLY)
	glog.Info("Complete")
	glog.Flush()
}

// rclone uses standard golang log package
// this will ensure that those logs get written to file
func initGolangLog() (*os.File, error) {
	f, err := os.OpenFile(otherLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		glog.Errorf("Error opening other log file: %v", err)
		return nil, err
	}

	log.SetOutput(f)
	log.Println("Other logging initted")
	return f, nil
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
