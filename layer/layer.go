package layer

import (
	"os"
	"os/exec"

	"github.com/gotoolkits/imageDigger/common"
	log "github.com/sirupsen/logrus"
)

var RootPath = "/data/docker"
var MountPath = RootPath + "/image/aufs/layerdb/mounts/"
var LayersPath = RootPath + "/aufs/layers/"
var MntsPath = RootPath + "/aufs/mnt/"

func DisplayImagesInfoFormContainerId(trimid string) {

	cmd := exec.Command("docker", "inspect", trimid)
	d, err := common.Execute(cmd)

	if err != nil {
		log.Errorln("execute the 'docker inspect' command failed!")
		os.Exit(1)
	}

	id, err := common.GetValueFromBytes(d, "Id")

	if err != nil {
		log.Errorln("get json Id is failed")
		os.Exit(1)
	}

	mountID := MountPath + id + "/" + "mount-id"
	log.Infof("Path:%s\nContainer id is:%s\n", mountID, id)

	mid, _ := common.GetFileToStr(mountID)

	log.Infof("mount-id is:\n%s", mid)

	getParentsLayers(mid)
	getImagesLayerDatasFromParentId(mid)

}

func getParentsLayers(id string) []string {
	path := LayersPath + id

	parents, err := common.GetFileToArry(path)

	if err != nil {
		log.Errorln(err)
	}

	log.Infof("Get Parents layers is:\n%s", parents)
	return parents

}

func getImagesLayerDatasFromParentId(id string) {
	path := MntsPath + id
	layerData := common.CommonandLs(path)
	log.Infof("Layer path:%s\n", path)
	log.Infof("Layer files:\n %s", layerData)

}
