package resizer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

func New(localDir string, imgSizes []int) *Resizer {
	return &Resizer{
		localDir: localDir,
		resizes:  imgSizes,
	}
}

type Resizer struct {
	localDir string
	resizes  []int
}

func (rez *Resizer) GetResizedImagesPath(srcImgPath string) ([]string, error) {
	resizedImgPaths := make([]string, 0, len(rez.resizes))

	for _, imgSize := range rez.resizes {
		resizeImgPath, err := resize(rez.localDir, srcImgPath, imgSize)
		if err != nil {
			return nil, err
		}
		resizedImgPaths = append(resizedImgPaths, resizeImgPath)
	}
	return resizedImgPaths, nil
}
func (rez *Resizer) RemoveImages() error {
	files, err := filepath.Glob(filepath.Join(rez.localDir, "*"))
	if err != nil {
		log.Fatalf("cannot join and append filepath : %v", err)
		return err
	}
	for _, file := range files {
		err := os.RemoveAll(file)
		if err != nil {
			log.Fatalf("cannot remove recursively dir : %v", err)
			return err
		}
	}
	return nil
}

func resize(localDir string, srcImgPath string, imgSize int) (string, error) {
	img, err := imaging.Open(srcImgPath, imaging.AutoOrientation(true))
	if err != nil {
		log.Fatalf("fail opening path, %s", err.Error())
		return "", err
	}
	resizedImg := imaging.Resize(img, imgSize, imgSize, imaging.Lanczos)
	resizedImgName := getResizedImageName(srcImgPath, imgSize)
	imgPath := localDir + "/" + resizedImgName
	err = imaging.Save(resizedImg, imgPath)
	if err != nil {
		log.Fatalf("resized img is not saved, %v", err)
		return "", err
	}
	return resizedImgName, nil

}
func getResizedImageName(srcImgPath string, imgWH int) string {
	srcImgWithExt := filepath.Base(srcImgPath)
	extImgSrc := filepath.Ext(srcImgWithExt)
	srcImgWithoutExt := strings.TrimSuffix(srcImgWithExt, extImgSrc)

	return fmt.Sprintf("%s_%d%s", srcImgWithoutExt, imgWH, extImgSrc)
}
