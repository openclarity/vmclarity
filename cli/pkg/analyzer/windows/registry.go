// Copyright © 2024 Cisco Systems, Inc. and its affiliates.
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

package windows

import (
	"encoding/xml"
	"fmt"
	cdx "github.com/CycloneDX/cyclonedx-go"
	"github.com/google/uuid"
	"github.com/openclarity/vmclarity/cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"www.velocidex.com/golang/regparser"
)

// Documentation about registries, its structure and further details can be found at:
// - https://en.wikipedia.org/wiki/Windows_Registry
// - https://techdirectarchive.com/2020/02/07/how-to-check-if-windows-updates-were-installed-on-your-device-via-the-registry-editor/
// - https://jgmes.com/webstart/library/qr_windowsxp.htm#:~:text=In%20Windows%20XP%2C%20the%20registry,corresponding%20location%20of%20each%20hive.

type Registry struct {
	drivePath   string              // root path to Windows drive
	softwareReg *regparser.Registry // HKEY_LOCAL_MACHINE/SOFTWARE registry
	cleanup     func() error
	logger      *log.Entry
}

func NewRegistry(drivePath string, logger *log.Entry) (*Registry, error) {
	// Windows XP has a different location for the registries "/Windows/system32/config/",
	// while Vista and upwards share the same location at "/Windows/System32/config/".
	// The registry key structure is almost identical for all Windows NT distributions.
	// Check: https://en.wikipedia.org/wiki/Windows_Registry#File_locations
	//
	// TODO(ramizpolic): We can check which registry file to open via loop.
	//  If needed to run on Windows, convert to valid path.
	registryPath := path.Join(drivePath, "/Windows/System32/config/SOFTWARE")
	registryFile, err := os.Open(registryPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open registry file %s: %w", registryPath, err)
	}

	// Registry file must remain open as it is read on-the-fly
	softwareReg, err := regparser.NewRegistry(registryFile)
	if err != nil {
		return nil, fmt.Errorf("cannot create registry reader: %w", err)
	}

	return &Registry{
		drivePath:   drivePath,
		softwareReg: softwareReg,
		cleanup:     registryFile.Close,
		logger:      logger,
	}, nil
}

// Close needs to be called when done to free up resources. Registry is not
// usable once closed.
func (r *Registry) Close() error {
	if err := r.cleanup(); err != nil {
		return fmt.Errorf("failed to close registry: %w", err)
	}
	return nil
}

// GetPlatform returns OS-specific data from the registry.
func (r *Registry) GetPlatform() (map[string]string, error) {
	// Open key to fetch operating system version and configuration data
	platformKey, err := openKey(r.softwareReg, "Microsoft/Windows NT/CurrentVersion")
	if err != nil {
		return nil, err
	}

	// Extract all platform data from the registry and strip all secrets
	platform := getValuesMap(platformKey)
	delete(platform, "DigitalProductId")
	delete(platform, "DigitalProductId4")

	return platform, nil
}

// GetUpdates returns a slice of all installed system updates from the registry.
func (r *Registry) GetUpdates() ([]string, error) {
	// Open key to fetch CBS data about packages (updates and components)
	packagesKey, err := openKey(r.softwareReg, "Microsoft/Windows/CurrentVersion/Component Based Servicing/Packages")
	if err != nil {
		return nil, err
	}

	// Extract all updates from installed packages
	updates := make(map[string]struct{})
	updateRegex := regexp.MustCompile("KB[0-9]{7,}")
	for _, pkgKey := range packagesKey.Subkeys() {
		pkgName := pkgKey.Name()
		pkgValues := getValuesMap(pkgKey)

		// Ignore packages that were not installed as system components or via updates
		_, isComponent := pkgValues["InstallClient"]
		_, isUpdate := pkgValues["UpdateAgentLCU"]
		if !isComponent && !isUpdate {
			continue
		}

		// Install location value for a given package can contain update identifier such
		// as "C:\Windows\CbsTemp\31075171_2144217839\Windows10.0-KB5032189-x64.cab\"
		if location, ok := pkgValues["InstallLocation"]; ok {
			if kb := updateRegex.FindString(location); kb != "" {
				updates[kb] = struct{}{}
			}
		}

		// If the installed package contains state value, it indicates a potential system
		// update. We are only curious about "112" state code which translates to
		// successful package installation. When this is the case, package registry key
		// contains update identifier such as "Package_10_for_KB5011048..."
		if state, ok := pkgValues["CurrentState"]; ok && state == "112" {
			if kb := updateRegex.FindString(pkgName); kb != "" {
				updates[kb] = struct{}{}
			}
		}
	}

	return utils.GetMapKeys(updates), nil
}

// GetUsersApps returns installed apps from all users
func (r *Registry) GetUsersApps() ([]map[string]string, error) {
	// Open key to fetch system user profiles in order to get their mount paths
	profilesKey, err := openKey(r.softwareReg, "Microsoft/Windows NT/CurrentVersion/ProfileList")
	if err != nil {
		return nil, err
	}

	// Extract all installed applications for each user
	apps := []map[string]string{}
	for _, profileKey := range profilesKey.Subkeys() {
		// Run in a function to allow cleanup
		func(profileValues map[string]string) {
			// Extract profile path from the registry key values. The path is
			// Windows-specific, but the mount path must be Unix-specific.
			// TODO(ramizpolic): If needed to run on Windows, convert to valid path.
			profileLocation, ok := profileValues["ProfileImagePath"]
			if !ok {
				return // silent skip, not a user profile
			}
			profileLocation = strings.ReplaceAll(profileLocation, "\\", "/")

			// The actual user location in the registry is specified as "C:/Users/...".
			// However, due to the actual mount location, the actual path could be
			// "/var/mounts/Users/...". Strip everything before the "/Users/" to construct a
			// valid mount path.
			if prefixIdx := strings.Index(profileLocation, "/Users/"); prefixIdx >= 0 {
				baseProfileLocation := profileLocation[prefixIdx:]
				profileLocation = path.Join(r.drivePath, baseProfileLocation)
			} else {
				return // silent skip, not a user profile
			}

			// Open profile registry file to access profile-specific registry
			profileRegPath := path.Join(profileLocation, "NTUSER.DAT")
			profileRegFile, err := os.Open(profileRegPath)
			if err != nil {
				r.logger.Warnf("failed to open user profile: %v", err)
				return
			}
			defer profileRegFile.Close()

			profileReg, err := regparser.NewRegistry(profileRegFile)
			if err != nil {
				r.logger.Warnf("failed to create user registry reader: %v", err)
				return
			}

			// Open key to fetch installed profile apps
			profileAppsKey, err := openKey(profileReg, "SOFTWARE/Microsoft/Windows/CurrentVersion/Uninstall")
			if err != nil {
				r.logger.Warnf("failed to open key: %v", err)
				return
			}

			// Extract all apps from user registry key. When the application registry key
			// values contain application name, add them to the result.
			for _, appKey := range profileAppsKey.Subkeys() {
				appValues := getValuesMap(appKey)
				if _, ok := appValues["DisplayName"]; ok {
					apps = append(apps, appValues)
				}
			}
		}(getValuesMap(profileKey))
	}

	return apps, nil
}

func (r *Registry) GetSystemApps() ([]map[string]string, error) {
	// Try multiple keys to fetch installed system apps
	apps := []map[string]string{}
	for _, appsKey := range []string{
		"Microsoft/Windows/CurrentVersion/Uninstall",             // for newer Windows NT
		"Wow6432Node/Microsoft/Windows/CurrentVersion/Uninstall", // store for 32-bit apps on 64-bit systems
		"WOW6432Node/Microsoft/Windows/CurrentVersion/Uninstall", // same as before, resolves compatibility issues
	} {
		appsKey, err := openKey(r.softwareReg, appsKey)
		if err != nil {
			r.logger.Warnf("failed to get installed system apps: %v", err)
			continue
		}

		// Extract all apps from system registry. When the application registry key
		// values contain application name, add them to the result.
		for _, appKey := range appsKey.Subkeys() {
			appValues := getValuesMap(appKey)
			if _, ok := appValues["DisplayName"]; ok {
				apps = append(apps, appValues)
			}
		}
	}

	return apps, nil
}

func (r *Registry) GetAll() map[string]interface{} {
	// Fetch required data from the registry
	platformData, _ := r.GetPlatform()
	updateData, _ := r.GetUpdates()
	usersApps, _ := r.GetUsersApps()
	systemApps, _ := r.GetSystemApps()

	// Create SBOM
	// TODO(ramizpolic): Convert to SBOM, check https://github.com/MartinStengard/rust-sbom-windows as reference
	_ = &cdx.BOM{
		XMLName:      xml.Name{},
		XMLNS:        "",
		JSONSchema:   "",
		BOMFormat:    cdx.BOMFormat,
		SpecVersion:  cdx.SpecVersion1_5,
		SerialNumber: uuid.New().URN(),
		Version:      1,
		Metadata: &cdx.Metadata{
			Timestamp:  "",
			Lifecycles: nil,
			Tools:      nil,
			Authors:    nil,
			Component: &cdx.Component{
				BOMRef:             "",
				MIMEType:           "",
				Type:               cdx.ComponentTypeOS,
				Supplier:           nil,
				Author:             platformData["SoftwareType"],
				Publisher:          "",
				Group:              "",
				Name:               platformData["ProductName"],
				Version:            "",
				Description:        "",
				Scope:              "",
				Hashes:             nil,
				Licenses:           nil,
				Copyright:          "",
				CPE:                "",
				PackageURL:         "",
				SWID:               nil,
				Modified:           nil,
				Pedigree:           nil,
				ExternalReferences: nil,
				Properties:         nil,
				Components:         nil,
				Evidence:           nil,
				ReleaseNotes:       nil,
				ModelCard:          nil,
				Data:               nil,
			},
			Manufacture: &cdx.OrganizationalEntity{
				Name:    "",
				URL:     nil,
				Contact: nil,
			},
			Supplier:   nil,
			Licenses:   nil,
			Properties: &[]cdx.Property{},
		},
		Components:         nil,
		Services:           nil,
		ExternalReferences: nil,
		Dependencies:       nil,
		Compositions:       nil,
		Properties:         nil,
		Vulnerabilities:    nil,
		Annotations:        nil,
		Formulation:        nil,
	}

	// TODO: to be replaced
	return map[string]interface{}{
		"platform":   platformData,
		"updates":    updateData,
		"usersApps":  usersApps,
		"systemApps": systemApps,
	}
}

// openKey opens a given registry key from the given registry or returns an error.
// Returned key can have multiple sub-keys and values specified.
func openKey(registry *regparser.Registry, key string) (*regparser.CM_KEY_NODE, error) {
	keyNode := registry.OpenKey(key)
	if keyNode == nil {
		return nil, fmt.Errorf("cannot open key %s", key)
	}
	return keyNode, nil
}

// getValuesMap returns all registry key values for a given registry key as a map
// of value name and its data.
func getValuesMap(key *regparser.CM_KEY_NODE) map[string]string {
	valuesMap := make(map[string]string)
	for _, keyValue := range key.Values() {
		valuesMap[keyValue.ValueName()] = convertKVData(keyValue.ValueData())
	}
	return valuesMap
}

// convertKVData returns the registry key value data as a valid string
func convertKVData(value *regparser.ValueData) string {
	switch value.Type {
	case regparser.REG_SZ, regparser.REG_EXPAND_SZ: // null-terminated string
		return strings.TrimRightFunc(value.String, func(r rune) bool {
			return r == 0
		})

	case regparser.REG_MULTI_SZ: // multi-part string
		return strings.Join(value.MultiSz, " ")

	case regparser.REG_DWORD, regparser.REG_DWORD_BIG_ENDIAN, regparser.REG_QWORD: // unsigned 32/64-bit value
		return strconv.FormatUint(value.Uint64, 10)

	case regparser.REG_BINARY: // non-stringable binary value
		// Return as hex to preserve buffer; we don't really care about this value
		return fmt.Sprintf("%X", value.Data)

	case
		regparser.REG_LINK,                       // unicode symbolic link
		regparser.REG_RESOURCE_LIST,              // device-driver resource list
		regparser.REG_FULL_RESOURCE_DESCRIPTOR,   // hardware setting
		regparser.REG_RESOURCE_REQUIREMENTS_LIST, // hardware resource list
		regparser.REG_UNKNOWN:                    // no-type
		fallthrough

	default:
		return ""
	}
}
