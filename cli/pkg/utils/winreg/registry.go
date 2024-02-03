// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

// TODO(ramizpolic): This is an MVP and will be heavily changed

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"www.velocidex.com/golang/regparser"
)

type Registry struct {
	registryPath string
	registry     *regparser.Registry
}

func NewRegistry(path string) (*Registry, error) {
	regFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open registry file: %w", err)
	}

	registry, err := regparser.NewRegistry(regFile)
	if err != nil {
		return nil, fmt.Errorf("cannot create registry reader: %w", err)
	}

	return &Registry{
		registryPath: path,
		registry:     registry,
	}, nil
}

func (r *Registry) GetPlatform() (map[string]string, error) {
	osRegistry := r.registry.OpenKey("Microsoft/Windows NT/CurrentVersion")
	if osRegistry == nil {
		return nil, fmt.Errorf("key not found")
	}

	data := map[string]string{}
	for _, prop := range osRegistry.Values() {
		data[prop.ValueName()] = toString(prop.ValueData())
	}

	// remove secrets
	delete(data, "DigitalProductId")
	delete(data, "DigitalProductId4")

	// inject identifiers
	data["Type"] = "Operating System"
	data["OsType"] = "Windows"

	return data, nil
}

func (r *Registry) GetUpdates() (map[string]string, error) {
	packageReg := r.registry.OpenKey("Microsoft/Windows/CurrentVersion/Component Based Servicing/Packages")
	if packageReg == nil {
		return nil, fmt.Errorf("key not found")
	}

	updateKeys := map[string]*regparser.CM_KEY_NODE{}
	for _, pkgKey := range packageReg.Subkeys() {
		for _, prop := range pkgKey.Values() {
			switch propName := prop.ValueName(); propName {
			case "InstallClient", "UpdateAgentLCU":
				updateKeys[pkgKey.Name()] = pkgKey
			}
		}
	}

	updates := map[string]string{}
	kbRegex := regexp.MustCompile("KB[0-9]{7,}")
	for _, updateKey := range updateKeys {
		for _, prop := range updateKey.Values() {
			propName := prop.ValueName()
			propValue := toString(prop.ValueData())

			if propName == "InstallLocation" {
				if kb := kbRegex.FindString(propValue); kb != "" {
					updates[kb] = kb
				}
			}
			if propName == "CurrentState" && propValue == "112" {
				if kb := kbRegex.FindString(updateKey.Name()); kb != "" {
					updates[kb] = kb
				}
			}
		}
	}

	return updates, nil
}

func (r *Registry) GetUserProfiles() (map[string]string, error) {
	profileReg := r.registry.OpenKey("Microsoft/Windows NT/CurrentVersion/ProfileList")
	if profileReg == nil {
		return nil, fmt.Errorf("key not found")
	}

	profiles := map[string]string{}
	for _, profileKey := range profileReg.Subkeys() {
		for _, prop := range profileKey.Values() {
			propName := prop.ValueName()
			if propName != "ProfileImagePath" { // only interested in this key
				continue
			}

			propValue := toString(prop.ValueData())
			propValue = strings.ReplaceAll(propValue, "\\", "/")
			if idx := strings.Index(propValue, "/Users/"); idx >= 0 {
				propValue = path.Join(r.getMountPath(), propValue[idx:])
				profiles[propValue] = propValue
			}
		}
	}

	return profiles, nil
}

// TODO: simplify the key insertions

func (r *Registry) GetUserApps() ([]map[string]string, error) {
	profiles, err := r.GetUserProfiles()
	if err != nil {
		return nil, err
	}

	apps := []map[string]string{}
	for profile := range profiles {
		profileRegPath := path.Join(profile, "NTUSER.DAT")

		profileFile, err := os.Open(profileRegPath)
		if err != nil {
			// check and log error
			continue
		}

		profileReg, err := regparser.NewRegistry(profileFile)
		if err != nil {
			// check and log error
			continue
		}

		appReg := profileReg.OpenKey("SOFTWARE/Microsoft/Windows/CurrentVersion/Uninstall")
		if appReg == nil {
			// reg key does not exist, log
			continue
		}

		for _, appKey := range appReg.Subkeys() {
			currApp := map[string]string{}
			for _, prop := range appKey.Values() {
				propName := prop.ValueName()
				propValue := toString(prop.ValueData())

				switch propName {
				case "DisplayName", "DisplayVersion", "VersionMajor", "VersionMinor":
					currApp[propName] = propValue
				}
			}

			// add to output
			// TODO: clean this up
			if _, ok := currApp["DisplayName"]; ok {
				apps = append(apps, currApp)
			}
		}
	}

	return apps, nil
}

func (r *Registry) GetSystemApps() ([]map[string]string, error) {
	apps := []map[string]string{}
	for _, regKey := range []string{
		"Microsoft/Windows/CurrentVersion/Uninstall",
		"Wow6432Node/Microsoft/Windows/CurrentVersion/Uninstall", // case sensitive, OS dependant
		"WOW6432Node/Microsoft/Windows/CurrentVersion/Uninstall", // case sensitive, OS dependant
	} {
		appReg := r.registry.OpenKey(regKey)
		if appReg == nil {
			// reg key does not exist, log
			continue
		}

		for _, appKey := range appReg.Subkeys() {
			currApp := map[string]string{}
			for _, prop := range appKey.Values() {
				propName := prop.ValueName()
				propValue := toString(prop.ValueData())

				switch propName {
				case "DisplayName", "DisplayVersion", "VersionMajor", "VersionMinor":
					currApp[propName] = propValue
				}
			}

			// add to output
			// TODO: clean this up
			if _, ok := currApp["DisplayName"]; ok {
				apps = append(apps, currApp)
			}
		}
	}

	return apps, nil
}

func (r *Registry) GetAll() map[string]interface{} {
	osData, _ := r.GetPlatform()
	updateData, _ := r.GetUpdates()
	profileData, _ := r.GetUserProfiles()
	userApps, _ := r.GetUserApps()
	systemApps, _ := r.GetSystemApps()
	return map[string]interface{}{
		"platform":   osData,
		"kbs":        updateData,
		"users":      profileData,
		"userApps":   userApps,
		"systemApps": systemApps,
	}
}

func (r *Registry) getMountPath() string {
	// mountPath + /Windows/System32/config/SOFTWARE
	// the mount path is everything before system path, i.e. /Windows/System32/
	if idx := strings.Index(r.registryPath, "/Windows/System32/"); idx >= 0 {
		return r.registryPath[:idx]
	}
	return "/"
}

// toString converts the data to proper go string.
// TODO: Handle UTF values
func toString(data *regparser.ValueData) string {
	switch data.Type {
	case regparser.REG_SZ, regparser.REG_EXPAND_SZ:
		return strings.TrimRightFunc(data.String, func(r rune) bool {
			return r == 0 // remove null terminator
		})

	case regparser.REG_MULTI_SZ:
		return strings.Join(data.MultiSz, " ")

	case regparser.REG_DWORD, regparser.REG_DWORD_BIG_ENDIAN, regparser.REG_QWORD:
		return strconv.FormatUint(data.Uint64, 10)

	case regparser.REG_BINARY:
		// Return as hex to preserve buffer; we don't really care about this value
		return fmt.Sprintf("%X", data.Data)

	case
		regparser.REG_LINK,                       // Unicode symbolic link
		regparser.REG_RESOURCE_LIST,              // device-driver resource list
		regparser.REG_FULL_RESOURCE_DESCRIPTOR,   // hardware setting
		regparser.REG_RESOURCE_REQUIREMENTS_LIST, // hardware resource list
		regparser.REG_UNKNOWN:                    // unhandled
		fallthrough

	default:
		return ""
	}
}

func printKeysForPath(registry *regparser.Registry, keyPath string) {
	key := registry.OpenKey(keyPath)
	if key == nil {
		return
	}

	fmt.Printf("======= Subkeys:\n\n")
	for _, subkey := range key.Subkeys() {
		fmt.Printf("%-120s\n", path.Join(keyPath, subkey.Name()))
	}

	// print all keys under this path
	fmt.Printf("\n\n======= Key properties:\n\n")
	for _, prop := range key.Values() {
		fmt.Printf("%-120s : %#v\n", path.Join(keyPath, prop.ValueName()), toString(prop.ValueData()))
	}
}
