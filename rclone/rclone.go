package rclone

import (
	"fmt"
	"strconv"

	"github.com/golang/glog"
	_ "github.com/ncw/rclone/backend/all" // needed to init the backends
	"github.com/ncw/rclone/fs"
	"github.com/ncw/rclone/fs/config"
	"github.com/ncw/rclone/fs/filter"
	"github.com/ncw/rclone/fs/sync"
)

// Sync configures rclone and executes the sync
func Sync(configPath string, bucketName string, mountpoint string, filterStrs []string, dryRun bool,
	hardDelete bool) error {
	config.ConfigPath = configPath

	fsrc, err := fs.NewFs(mountpoint)
	if err != nil {
		glog.Errorf("Error creating src fs: %s", err)
		return err
	}

	fdst, err := getB2Fs(fmt.Sprintf("backblaze:%s", bucketName), hardDelete)
	if err != nil {
		glog.Fatalf("Error creating dst fs: %s", err)
	}

	// get filters
	filters, err := getRcloneFiltersFromStrs(filterStrs)
	if err != nil {
		glog.Errorf("Error parsing filters: %s", err)
		return err
	}
	filter.Active = filters

	fs.Config.DryRun = dryRun

	return sync.Sync(fdst, fsrc)
}

func getRcloneFiltersFromStrs(filters []string) (*filter.Filter, error) {
	opt := filter.DefaultOpt
	for _, e := range filters {
		opt.ExcludeRule = append(opt.ExcludeRule, e)
	}

	return filter.NewFilter(&opt)
}

// this is basically just https://github.com/ncw/rclone/blob/9322f4baef55380a88d2a5d5ce2fa1c12f7e51f0/fs/fs.go#L1045
// but it sets some additional config on the configmap
func getB2Fs(path string, hardDelete bool) (fs.Fs, error) {
	fsInfo, configName, fsPath, config, err := fs.ConfigFs(path)
	if err != nil {
		return nil, err
	}

	// rclone does some reflection magic to generate config structs from the config map
	// the struct fields get mapped to snake case
	// https://github.com/ncw/rclone/blob/9322f4baef55380a88d2a5d5ce2fa1c12f7e51f0/fs/config/configstruct/configstruct.go#L94
	config.Set("hard_delete", strconv.FormatBool(hardDelete))

	return fsInfo.NewFs(configName, fsPath, config)
}
