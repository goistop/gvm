package action

import (
	"regexp"
	"github.com/ntfs32/gvm/src/gvm/download"
	"github.com/ntfs32/gvm/src/gvm/file"
	"github.com/ntfs32/gvm/src/gvm/config"
	"fmt"
	"github.com/urfave/cli"
)

func Install(c *cli.Context) error {
	//fmt.Println("Install version: ")
	version := c.Args().First()
	reg := regexp.MustCompile(`^go\d{1,2}\.\d{1,2}\.\d{1,2}$`)
	if reg.MatchString(version) {
		err :=download.DownReleaseVersion(version)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		file.Unpackrelease(download.GetTarFileName(version),version)
	}else{
		fmt.Println("Version format error")
	}
	return nil
}




func ListLocal(c *cli.Context) error {
	fmt.Println("Version had installed local")
	versionConfig := config.Get()
	for _,version := range versionConfig.Version.InstalledVersion{
		fmt.Println("\t"+ version)
	}
	return nil
}


func ListRemote(c *cli.Context) error {
	versionConfig := download.GetReleaseVersion()
	for _,version := range versionConfig{
		fmt.Println("\t"+ version)
	}
	return nil
}