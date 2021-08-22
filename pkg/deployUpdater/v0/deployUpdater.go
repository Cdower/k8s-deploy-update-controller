package deployupdater

import (
	"fmt"
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
	// fmt.Println("new deploy updater")
	d.settings = utils.NewSettings()
	d.kube = kubeV0.NewKubeClient(d.settings)
	return d
}

func (d *DeployUpdater) Run() {
	// fmt.Println("Run")
	d.kube.GetDeployments()
	latest := d.getLatestRelease()
	if latest != d.kube.GetDeploymentVersion() {
		fmt.Println("Need to update deployment with latest release.")
		d.kube.UpdateDeploymentVersion(latest)
	}
}

func (d *DeployUpdater) getLatestRelease() string {
	// Do not follow redirect
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	// fmt.Printf("Checking channel url: %s\n", d.settings.Channel)
	resp, err := client.Get(d.settings.Channel)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode == 302 {
		// fmt.Println(resp.Header.Get("location"))
		loc, err := url.Parse(resp.Header.Get("location"))
		if err != nil {
			panic(err.Error())
		}
		return path.Base(loc.Path)
	} else {
		return ""
	}
}
