package apparmor

import (
	"strings"

	"github.com/Sirupsen/logrus"
	aaprofile "github.com/docker/docker/profiles/apparmor"
	"github.com/opencontainers/runc/libcontainer/apparmor"
)

const (
	// defaultApparmorProfile is the name of default apparmor profile name.
	defaultApparmorProfile = "crio-default"

	// ContainerAnnotationKeyPrefix is the prefix to an annotation key specifying a container profile.
	ContainerAnnotationKeyPrefix = "container.apparmor.security.beta.kubernetes.io/"

	// ProfileRuntimeDefault is he profile specifying the runtime default.
	ProfileRuntimeDefault = "runtime/default"
	// ProfileNamePrefix is the prefix for specifying profiles loaded on the node.
	ProfileNamePrefix = "localhost/"
)

func installDefaultAppArmorProfile() {
	if apparmor.IsEnabled() {
		if err := aaprofile.InstallDefault(defaultApparmorProfile); err != nil {
			apparmorProfiles := []string{defaultApparmorProfile}

			// Allow daemon to run if loading failed, but are active
			// (possibly through another run, manually, or via system startup)
			for _, policy := range apparmorProfiles {
				if err := aaprofile.IsLoaded(policy); err != nil {
					logrus.Errorf("AppArmor enabled on system but the %s profile could not be loaded.", policy)
				}
			}
		}
	}
}

// GetAppArmorProfileName gets the profile name for the given container.
func GetAppArmorProfileName(annotations map[string]string, ctrName string) string {
	profile := GetProfileNameFromPodAnnotations(annotations, ctrName)
	if profile == "" || profile == ProfileRuntimeDefault {
		// If the value is runtime/default, then it is equivalent to not specifying a profile.
		return ""
	}

	profileName := strings.TrimPrefix(profile, ProfileNamePrefix)
	return profileName
}

// GetProfileNameFromPodAnnotations gets the name of the profile to use with container from
// pod annotations
func GetProfileNameFromPodAnnotations(annotations map[string]string, containerName string) string {
	return annotations[ContainerAnnotationKeyPrefix+containerName]
}

