package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func main() {
	osType := runtime.GOOS
	fmt.Println("Detected OS Type:",osType)

	switch(osType){
	case "darwin":
		fmt.Println("MacOS is not supported")
	case "windows":
		fmt.Println("Windows is not supported")
	case "linux":
		distro, base_distro := fetchLinusDistro()
		if(distro=="unknown" && base_distro=="unknown"){fmt.Println("Failed to fetch the details of your distro")}
		fmt.Println("Distribution:", distro)
		fmt.Println("Built On:", base_distro)
	default:
		fmt.Println("OS not supported")

	}
}

func fetchLinusDistro() (string, string) {
	data, err := os.ReadFile("/etc/os-release")

	if (err != nil) {return "unknown", "unknown"}
	var distro, base_distro string

	lines:= strings.Split(string(data),"\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID="){
			distro =strings.Trim(strings.SplitN(line,"=",2)[1],`"`)
		}
		if (strings.HasPrefix(line, "ID_LIKE=")){
			base_distro = strings.Trim(strings.SplitN(line, "=",2)[1],`"`)
		}
	}

	return distro, base_distro
}