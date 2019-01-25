package restore

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"strconv"

	utils "github.com/maorfr/helm-plugin-utils/pkg"
)

// Restore performs a restore of a release
func Restore(releaseName, tillerNamespace, label string) error {
	listOptions := utils.ListOptions{
		ReleaseName:     releaseName,
		TillerNamespace: tillerNamespace,
		TillerLabel:     label,
	}
	releases, err := utils.ListReleases(listOptions)
	if err != nil {
		return err
	}
	if len(releases) != 1 {
		return fmt.Errorf("%s has no deployed releases", releaseName)
	}

	fileName := "/tmp/helm-restore-manifests."+strconv.Itoa(os.Getpid())+".yaml"
	os.Remove(fileName)
	if err := ioutil.WriteFile(fileName, []byte(releases[0].Manifest), 0644); err != nil {
		return err
	}
	applyCmd := []string{"oc", "apply", "--namespace", releases[0].Namespace, "-f", fileName}
	output := utils.Execute(applyCmd)
	for _, line := range strings.Split((string)(output), "\n") {
		if line == "" {
			continue
		}
		log.Print(line)
	}
	//os.Remove(fileName)
	return nil
}
