package mount

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"k8s.io/utils/mount"
)

var pairsRE = regexp.MustCompile(`([A-Z]+)=(?:"(.*?)")`)

type BlockDevice struct {
	DeviceName     string
	Size           uint64
	Label          string
	UUID           string
	FilesystemType string
	MountPoint     string
}

// ListBlockDevices Taken from https://github.com/BishopFox/dufflebag
func ListBlockDevices() ([]BlockDevice, error) {
	log.Info("Listing block devices...")
	columns := []string{
		"NAME",       // name
		"SIZE",       // size
		"LABEL",      // filesystem label
		"UUID",       // filesystem UUID
		"FSTYPE",     // filesystem type
		"TYPE",       // device type
		"MOUNTPOINT", // device mountpoint
	}

	log.Info("executing lsblk...")
	output, err := exec.Command(
		"lsblk",
		"-b", // output size in bytes
		"-P", // output fields as key=value pairs
		"-o", strings.Join(columns, ","),
	).Output()
	if err != nil {
		return nil, fmt.Errorf("cannot list block devices: %v", err)
	}

	blockDeviceMap := make(map[string]BlockDevice)
	s := bufio.NewScanner(bytes.NewReader(output))
	for s.Scan() {
		pairs := pairsRE.FindAllStringSubmatch(s.Text(), -1)
		var dev BlockDevice
		var deviceType string
		for _, pair := range pairs {
			switch pair[1] {
			case "NAME":
				dev.DeviceName = pair[2]
			case "SIZE":
				size, err := strconv.ParseUint(pair[2], 10, 64)
				if err != nil {
					log.Warnf(
						"Invalid size %q from lsblk: %v", pair[2], err,
					)
				} else {
					// the number of bytes in a MiB.
					dev.Size = size / 1024 * 1024
				}
			case "LABEL":
				dev.Label = pair[2]
			case "UUID":
				dev.UUID = pair[2]
			case "FSTYPE":
				dev.FilesystemType = pair[2]
			case "TYPE":
				deviceType = pair[2]
			case "MOUNTPOINT":
				dev.MountPoint = pair[2]
			default:
				log.Warnf("unexpected field from lsblk: %q", pair[1])
			}
		}

		if deviceType == "loop" {
			continue
		}

		blockDeviceMap[dev.DeviceName] = dev
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("cannot parse lsblk output: %v", err)
	}

	blockDevices := make([]BlockDevice, 0, len(blockDeviceMap))
	for _, dev := range blockDeviceMap {
		blockDevices = append(blockDevices, dev)
	}
	return blockDevices, nil
}

func (b BlockDevice) Mount(mountPoint string) error {
	// Make a directory for the device to mount to
	if err := os.MkdirAll(mountPoint, os.ModePerm); err != nil {
		return fmt.Errorf("failed to run mkdir comand: %v", err)
	}

	// Do the mount

	mounter := mount.New(mountPoint)
	if err := mounter.Mount("/dev/"+b.DeviceName, mountPoint, b.FilesystemType, nil); err != nil {
		return fmt.Errorf("failed to run mount command: %v", err)
	}
	b.MountPoint = mountPoint

	return nil
}
