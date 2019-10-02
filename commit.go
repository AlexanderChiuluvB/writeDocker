package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os/exec"
)

/*
通过commitContainer函数来实现把容器的文件系统打包成镜像文件
 */
func commitContainer(imageName string) {
	mntURL := "/opt/test2/mnt"
	imageTar := "/opt/test2/" + imageName + ".tar"
	fmt.Printf("%s",imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		log.Errorf("Tar folder %s error %v", mntURL, err)
	}
}


