package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/golang/glog"
	watchdog "github.com/jchorl/watchdog/client"
	watchdogTypes "github.com/jchorl/watchdog/types"
)

type child struct {
	Name       string  `json:"name"`
	Mountpoint *string `json:"mountpoint"`
	Size       string  `json:"size"`
}

type blockDevice struct {
	Children []child `json:"children"`
}

type lsblkOutput struct {
	BlockDevices []blockDevice `json:"blockDevices"`
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var excludeFlags arrayFlags

func main() {
	flag.Var(&excludeFlags, "exclude", "Exclude flags passed directly to rclone")
	flag.Parse()

	err := mountDrive()
	if err != nil {
		glog.Fatalf("Error mounting drive: %s", err)
	}

	bucketName := "backup-c3bac1e7-d888-4940-8778-f03adcddfe45"

	cmdName := "rclone"
	cmdArgs := []string{"--config=./rclone.conf", "sync", "/sync", fmt.Sprintf("backblaze:%s", bucketName), "--exclude", "exclude/"}
	var cmdOut []byte
	cmdStr := fmt.Sprintf("%s %s", cmdName, strings.Join(cmdArgs, " "))
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).CombinedOutput(); err != nil {
		glog.Errorln(string(cmdOut))
		glog.Fatalf("There was an error running command \"%s\": %s", cmdStr, err)
	}

	glog.Infof("Successfully ran \"%s\"", cmdStr)
	glog.Infof(string(cmdOut))

	wdClient := watchdog.Client{"https://watchdog.joshchorlton.com"}
	wdClient.Ping("nasblaze", watchdogTypes.Watch_DAILY)
}

func setupSyncDir() {
	err := os.MkdirAll("/sync/exclude/", 0700)
	if err != nil {
		glog.Fatalf("Error create /sync/exclude dir: %s", err)
	}
	err = os.MkdirAll("/sync/include/", 0700)
	if err != nil {
		glog.Fatalf("Error create /sync/include dir: %s", err)
	}

	os.OpenFile("/sync/include/file", os.O_RDONLY|os.O_CREATE, 0666)
	os.OpenFile("/sync/exclude/excludefile", os.O_RDONLY|os.O_CREATE, 0666)
}

func mountDrive() error {
	var (
		cmdOut []byte
		err    error
	)
	cmdName := "lsblk"
	cmdArgs := []string{"--json"}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		return err
	}

	parsed := lsblkOutput{}
	if err := json.Unmarshal(cmdOut, &parsed); err != nil {
		return err
	}

	sizeToFind := "931.3G"
	for _, blockDevice := range parsed.BlockDevices {
		for _, child := range blockDevice.Children {
			if child.Size == sizeToFind {
				glog.Infof("Found drive at %s", child.Name)
				if child.Mountpoint != nil && *child.Mountpoint == "/sync" {
					glog.Infof("Child is already mounted at /sync")
					return nil
				} else if child.Mountpoint != nil {
					return fmt.Errorf("Drive is already mounted somewhere else: %s", child.Mountpoint)
				} else {
					// mount that drive
					glog.Infof("Mounting child at /sync")
					cmd := "mount"
					args := []string{fmt.Sprintf("/dev/%s", child.Name), "/sync"}
					if err := exec.Command(cmd, args...).Run(); err != nil {
						return err
					}
					return nil
				}
			}
		}
	}

	return fmt.Errorf("Failed to find a drive with size: %s", sizeToFind)
}
