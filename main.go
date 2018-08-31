package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"

	"github.com/golang/glog"
)

func init() {
	flag.Parse()
}

func main() {
	cmdName := "rclone"
	cmdArgs := []string{"--config=./rclone.conf", "lsd", "backblaze:"}

	var (
		cmdOut []byte
		err    error
	)
	cmdStr := fmt.Sprintf("%s %s", cmdName, strings.Join(cmdArgs, " "))
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		glog.Fatalf("There was an error running command \"%s\": %s", cmdStr, err)
	}

	glog.Infof("Successfully ran \"%s\"", cmdStr)
	glog.Infof(string(cmdOut))
}
