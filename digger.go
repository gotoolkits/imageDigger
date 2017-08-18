package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gotoolkits/backend/config"
	log "github.com/sirupsen/logrus"
)

var RootPath = "/data/docker"
var MountPath = RootPath + "/image/aufs/layerdb/mounts/"
var LayerPath = RootPath + "/"

func main() {

	cmd := exec.Command("docker", "inspect", "c4ec0fe067a5")
	d, err := Execute(cmd)

	fmt.Println(string(d))
	if err != nil {
		log.Errorln("execute the 'docker inspect' command failed!")
		os.Exit(1)
	}

	id, err := GetValueFromBytes(d, "Id")

	if err != nil {
		log.Errorln("get json Id is failed")
		os.Exit(1)
	}

	fmt.Println("docker id is:", id)

	mountID := MountPath + id + "mount-id"

	mid, _ := GetFileToStr(mountID)

	fmt.Println(mid)

}

func getFilelist(path string) []string {
	var paths []string
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		paths = append(paths, path)
		return nil
	})

	if err != nil {
		log.Errorf("filepath returned %v\n", err)
		return nil
	}

	return paths
}

func Execute(c *exec.Cmd) ([]byte, error) {
	var err error
	var d []byte

	stdout, _ := c.StdoutPipe()
	err = c.Start()
	if err != nil {
		return nil, err
	}

	d, err = ioutil.ReadAll(stdout)

	end := len(d) - 2
	js := d[1:end]

	if err != nil {
		return nil, err
	}

	err = c.Wait()
	if err != nil {
		return nil, err
	}

	return js, nil

}

func GetValueFromBytes(bs []byte, jkey string) (string, error) {
	jc, err := rrconfig.LoadJsonConfigFromBytes(bs)
	if err != nil {
		return "", err
	}

	v, err := jc.GetString(jkey)

	if err != nil {
		return "", err
	}

	return v, nil
}

func GetValueFromFile(path string, jkey string) (string, error) {

	jc, err := rrconfig.LoadJsonConfigFromFile(path)

	if err != nil {
		return "", err
	}

	v, err := jc.GetString(jkey)

	if err != nil {
		return "", err
	}

	return v, nil

}

func GetFileToArry(path string) ([]string, error) {

	fb, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	strs := strings.Split(string(fb), "\n")

	return strs, nil

}

func GetFileToStr(path string) (string, error) {

	fb, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(fb), nil

}
