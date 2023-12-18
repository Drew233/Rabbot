package scheduler

import (
	"os"
	"path/filepath"

	"rabbot/internal/log"
)

// 定时删除./data/tmp目录下的文件，默认一天删除一次
func deleteTmpFiles(dirPath string) {
	log.RabLog.Infof("Clean tmp path %s begin", dirPath)
	deleteFiles(dirPath)
}

// 定时删除./data/pic目录下的文件，默认五分钟清理一次
func deletePicFiles(dirPath string) {
	log.RabLog.Debugf("Clean pic path %s begin", dirPath)
	deleteFiles(dirPath)
}

func deleteFiles(dirPath string) {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.RabLog.Errorf("Walk file %s failed, %v", path, err)
			return err
		}
		if !info.IsDir() {
			err := os.Remove(path)
			if err != nil {
				log.RabLog.Errorf("Remove file %s failed, %v", dirPath, err)
			} else {
				log.RabLog.Infof("Remove file %s successed", dirPath)
			}
		}

		return nil
	})
	if err != nil {
		log.RabLog.Errorf("Walk dirpath %s failed, %v", dirPath, err)
		return
	}
}