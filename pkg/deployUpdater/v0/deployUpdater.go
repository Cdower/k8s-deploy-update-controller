package deployupdater

import (
	"log"
	"net/http"
	"net/url"
	"path"

	kubeV0 "deployUpdater/pkg/kube/v0"
	"deployUpdater/pkg/utils"
)

type DeployUpdater struct {
	settings *utils.Specification
	kube     *kubeV0.KubeClient
}

func NewDeployUpdater() *DeployUpdater {
	d := &DeployUpdater{}
	// log.Println("new deploy updater")
	d.settings = utils.NewSettings()
	d.kube = kubeV0.NewKubeClient(d.settings)
	return d
}

func (d *DeployUpdater) Run() {
	// log.Println("Run")
	d.kube.PrintDeployments()
	latest := d.getLatestRelease()
	deployed_version := d.kube.GetDeploymentVersion()
	if latest != deployed_version {
		log.Printf("Need to update deployment with latest release. From %s to %s", deployed_version, latest)
		err := d.kube.UpdateDeploymentVersion(latest)
		if err != nil {
			panic(err)
		}
		log.Printf("Deploy successfully applied with version %s.\n", latest)
	}
}

func (d *DeployUpdater) getLatestRelease() string {
	// Do not follow redirect
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	// log.Printf("Checking channel url: %s\n", d.settings.Channel)
	resp, err := client.Get(d.settings.Channel)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode == 302 {
		// log.Println(resp.Header.Get("location"))
		loc, err := url.Parse(resp.Header.Get("location"))
		if err != nil {
			panic(err.Error())
		}
		return path.Base(loc.Path)
	} else {
		return ""
	}
}
