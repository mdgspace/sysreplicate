package system

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"github.com/mdgspace/sysreplicate/system/output"
	"github.com/mdgspace/sysreplicate/system/utils"
)

// Run is the entry point for the system orchestrator.
func Run() {
	osType := runtime.GOOS
	fmt.Println("Detected OS Type:", osType)

	switch osType {
	case "darwin":
		fmt.Println("MacOS is not supported")
	case "windows":
		fmt.Println("Windows is not supported")
	case "linux":
		distro, baseDistro := utils.DetectDistro()
		if distro == "unknown" && baseDistro == "unknown" {
			log.Println("Failed to fetch the details of your distro")
		}
		fmt.Println("Distribution:", distro)
		fmt.Println("Built On:", baseDistro)
		packages := utils.FetchPackages(baseDistro)
		jsonObj, err := output.BuildSystemJSON(osType, distro, baseDistro, packages)
		if err != nil {
			log.Println("Error marshalling JSON:", err)
		} else {
			if err := os.MkdirAll(outputSysDir, 0744); err != nil {
				log.Println("Error creating sys output directory:", err)
			}
			if err := os.WriteFile(jsonOutputPath, jsonObj, 0644); err != nil {
				log.Println("Error writing JSON output:", err)
			}
			if err := os.MkdirAll(outputScriptsDir, 0744); err != nil {
				log.Println("Error creating scripts output directory:", err)
			}
			if err := output.GenerateInstallScript(baseDistro, packages, scriptOutputPath); err != nil {
				log.Println("Error generating install script:", err)
			} else {
				fmt.Println("Script generated successfully at:", scriptOutputPath)
			}
		}
	default:
		fmt.Println("OS not supported")
	}
}