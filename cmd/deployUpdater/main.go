package main

import (
	deployUpdater "deployUpdater/pkg/deployUpdater/v0"
)

func main() {
	// fmt.Println("cmd")
	d := deployUpdater.NewDeployUpdater()
	d.Run()
}
