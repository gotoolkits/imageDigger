package main

import (
	"flag"
	"os"
	"strings"

	"github.com/gotoolkits/imageDigger/image"
	"github.com/gotoolkits/imageDigger/layer"
	log "github.com/sirupsen/logrus"
)

func main() {

	types := flag.String("t", "image", "type:image/container")
	id := flag.String("id", "", "id")

	flag.Parse()

	if *id == "" {
		log.Infoln("Please set the id")
		flag.Usage()
		os.Exit(1)
	}

	if strings.Contains(*types, "image") {
		image.DisplayImagesInfoFormImagesId(*id)
		image.GetImagesLayerDatasFromCacheId()
	}

	if strings.Contains(*types, "container") {
		layer.DisplayImagesInfoFormContainerId(*id)
	}

}
