package main

import (
	"fmt"
	"os"
	"os/exec"
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
		fmt.Println(fetchPackages(base_distro))
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

func fetchPackages(base_distro string) []string {
	var cmd *exec.Cmd
	switch(base_distro) {
	case "debian":
		cmd = exec.Command("dpkg", "--get-selections")
	case "arch":
		cmd = exec.Command("pacman", "-Q")
	case "rhel", "fedora":
		cmd = exec.Command("rpm", "-qa")
	case "void":
		cmd = exec.Command("xbps-query", "-l")
	default :
		fmt.Println("Your distro is unsupported, cannot identify package manager !")
		return []string{"unknown"}
	}

	output, err := cmd.CombinedOutput()
	if (err != nil) {
		fmt.Println("Error in retrieving packages: ", err)
	}
	fmt.Println("Installed Packages: ")
	packages:=strings.Split(string(output), "\n")
	return packages
}	