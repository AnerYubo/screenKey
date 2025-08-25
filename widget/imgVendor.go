package widget

import (
	"embed"
	"fmt"
	"path"
)

type SysImg struct{}

//go:embed images/*
var ImgFs embed.FS

func NewImgVendor() *SysImg {
	return &SysImg{}
}

func (s *SysImg) GetImg(fileName string) ([]byte, error) {
	trayIcon, err := ImgFs.ReadFile(path.Join("images", fileName))
	if err != nil {
		fmt.Println("找不到图片文件：", err.Error())
		return nil, err
	}
	return trayIcon, nil
}
