package download

import (
	"os"
	"runtime"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"strings"
	"github.com/ntfs32/gvm/src/gvm/config"
	"github.com/ntfs32/gvm/src/gvm/file"
	"os/exec"
)


var platform = runtime.GOOS

var arch = runtime.GOARCH

var PlatformEnd = map[string]string{
	"linux": ".tar.gz",
	"darwin":".tar.gz",
	"freebsd": ".tar.gz",
	"Windows": ".zip",

}

var versionArr []string

func init()  {
	
}

func GetLink(version string) (url string) {
	if os.Getenv("GVM_SOURCE") == "GITHUB" {
		url = config.Get().GithubAddress + "go"+version + PlatformEnd[platform]
	}else{
		url = config.Get().GolangAddress+version+"."+platform+"-"+arch+PlatformEnd[platform]
	}
	return
}

func GetReleaseVersion()(versionArr []string) {
	doc, err := goquery.NewDocument(config.Get().ReleaseHistoryAddress)
	if err != nil {
		fmt.Println("Get Release list failed")
		return
	}
	fmt.Println("Golang Release Version:")
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		item := s.Text()
		if strings.Count(item, "(released") >= 1 {
			versionArr = append(versionArr, item[1:strings.Index(item, "(released")])
		}
	})
	return
}

func SystemExistsVersion()(string)  {
	goroot := os.Getenv("GOROOT")
	//var cmdOut []byte
	if goroot != ""{
		if cmdOut, err := exec.Command(goroot+"/bin/go","version").Output();err !=nil{
			return ""
		}else {
			return strings.Split(string(cmdOut)," ")[3]
		}

	}
	return ""
}

func GetTarFileName(version string) string {
	return version+"."+platform+"-"+arch+PlatformEnd[platform]
}

func DownReleaseVersion(version string) error{
	fileName := file.PackageDir + GetTarFileName(version)
	if _,err := os.Stat(fileName); err !=nil{
		fmt.Println("start download...")
		return file.DownloadFile(GetLink(version),fileName)
	}else{
		fmt.Println("this version "+fileName+" file exists.")
		return nil
	}
}
