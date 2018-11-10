package drives

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/golang/glog"
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

// MountDriveBySize mounts a drive with the given size
func MountDriveBySize(size string, mountpoint string) error {
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

	for _, blockDevice := range parsed.BlockDevices {
		for _, child := range blockDevice.Children {
			if child.Size == size {
				glog.Infof("Found drive at %s", child.Name)
				if child.Mountpoint != nil && *child.Mountpoint == mountpoint {
					glog.Infof("Child is already mounted at %s", mountpoint)
					return nil
				} else if child.Mountpoint != nil {
					glog.Errorf("Drive is already mounted somewhere else: %s", child.Mountpoint)
					return fmt.Errorf("Drive is already mounted somewhere else: %s", child.Mountpoint)
				} else {
					// mount that drive
					glog.Infof("Mounting child at %s", mountpoint)
					cmd := "mount"
					args := []string{fmt.Sprintf("/dev/%s", child.Name), mountpoint}
					if err := exec.Command(cmd, args...).Run(); err != nil {
						glog.Errorf("Error mounting drive: %s", err)
						return err
					}
					return nil
				}
			}
		}
	}

	glog.Errorf("Failed to find a drive with size: %s", size)
	return fmt.Errorf("Failed to find a drive with size: %s", size)
}
