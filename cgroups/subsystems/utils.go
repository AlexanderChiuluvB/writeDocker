package subsystems

import (
	"fmt"
	"strings"
	"os"
	"path"
	"bufio"
)


//找出cgroup的挂载点
func FindCgroupMountpoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		//47 33 0:42 / /sys/fs/cgroup/memory rw,nosuid,nodev,noexec,relatime shared:24 - cgroup cgroup rw,memory
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				//47 33 0:42 / /sys/fs/cgroup/memory rw,nosuid,nodev,noexec,relatime shared:24 - cgroup cgroup rw,memory
				//fields[4] = /sys/fs/cgroup/memory
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}


func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsystem)
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err == nil {
			} else {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}