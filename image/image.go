package image

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/gotoolkits/imageDigger/common"
	log "github.com/sirupsen/logrus"
)

var RootPath = "/data/docker"
var LayerPath = RootPath + "/aufs/diff/"
var MetaPath = RootPath + "/image/aufs/imagedb/content/sha256/"
var CachePath = RootPath + "/image/aufs/layerdb/sha256/"
var imageMetaFile string
var chainIDs []string
var cacheIDs []string

func DisplayImagesInfoFormImagesId(trimid string) {

	imagesMetaFileList := common.GetFilelist(MetaPath)
	for _, v := range imagesMetaFileList {
		if strings.Contains(v, trimid) {
			imageMetaFile = v
		}
	}
	if imageMetaFile == "" {
		log.Infoln("No found the id file in", MetaPath)
	} else {

		diffIds := getDiffIds(imageMetaFile)
		getMountIdFormCalcChainId(diffIds)
		chainIdToCacheId()
	}

	log.Infof("Get Cache-Ids:\nPATH:%s\n%s", CachePath, cacheIDs)

}

func getDiffIds(file string) []string {

	ids, err := common.GetSliceFromFile(file, "rootfs.diff_ids")
	if nil != err {
		log.Errorln(err)
		return nil
	}
	return ids
}

func getMountIdFormCalcChainId(diffIds []string) {

	log.Infof("Get Diff-Ids:\n%s", diffIds)

	chainId := sha256Sum(diffIds[0], diffIds[1])

	//slice to save the chainIds
	chainIDs = append(chainIDs, diffIds[0])
	chainIDs = append(chainIDs, chainId)

	log.Infoln("Get Chain-Ids:")
	log.Infoln(chainId)

	for i, v := range diffIds {
		if i <= 1 {
			continue
		}
		chainId := sha256Sum(chainIDs[i-1], v)
		chainIDs = append(chainIDs, chainId)
		log.Infoln(chainId)
	}
}

func sha256Sum(ch, pts string) string {

	child := bytes.NewBufferString(ch).Bytes()
	parents := bytes.NewBufferString(pts).Bytes()
	chainId := chainIdToDiffId(child, parents)
	chainId = "sha256:" + chainId

	return chainId
}

func chainIdToDiffId(child, parents []byte) string {

	sep := []byte(" ")
	data := [][]byte{child, parents}

	hash := sha256.New()
	hash.Write(bytes.Join(data, sep))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)

	return mdStr
}

func GetImagesLayerDatasFromCacheId() {

	for i, v := range cacheIDs {
		path := LayerPath + v
		layerData := common.CommonandLs(path)
		log.Infof("Layer path:%s\n", path)
		log.Infof("Layer %d files:\n %s", i, layerData)
	}

}

func chainIdToCacheId() {
	for _, v := range chainIDs {
		cacheId := getCacheId(v)
		cacheIDs = append(cacheIDs, cacheId)
	}
}

func getCacheId(idShaSum string) string {
	chainId := strings.Split(idShaSum, ":")
	path := CachePath + chainId[1] + "/cache-id"

	cacheId, err := common.GetFileToStr(path)
	if err != nil {
		log.Errorln(path, err)
		return ""
	}
	return cacheId
}
