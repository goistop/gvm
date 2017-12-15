package file

import (
	"os"
	"bytes"
	"os/exec"
	"os/user"
	"archive/zip"
	"strings"
	"runtime"
	"errors"
	"log"
	"fmt"
	"io"
	"bufio"
	"path/filepath"
	"net/http"
	"time"
	"strconv"
	"gopkg.in/cheggaaa/pb.v1"
	"archive/tar"
	"compress/gzip"
	"path"
)

var (
	ErrInvalid    = errors.New("invalid argument")
	ErrPermission = errors.New("permission denied")
	ErrExist      = errors.New("file already exists")
	ErrNotExist   = errors.New("file does not exist")
)

var (
	InstallDir = "/.gvm/install/"
	PackageDir = "/.gvm/packages/"
	ConfigFile = "/.gvm/config.toml"
)

var defaultConfig = `
# gvm manager config, You should not edit it
GolangAddress = "https://redirector.gvt1.com/edgedl/go/"
GithubAddress = "https://github.com/golang/go/archive/"
ReleaseHistoryAddress = "https://golang.org/doc/devel/release.html"

[Version]
ReleaseVersion = []
InstalledVersion = []
UsedVersion = ""
`

func init()  {
	homeDir,err :=Home()
	if err !=nil{
		log.Fatal("can not get HOME path")
		return
	}
	InstallDir = homeDir+ InstallDir
	PackageDir = homeDir+ PackageDir
	ConfigFile = homeDir+ ConfigFile
	if PathExists(InstallDir) == false {
		err :=os.MkdirAll(InstallDir,os.ModePerm)
		if err != nil{
			 fmt.Println(ErrPermission)
		}
	}
	if PathExists(PackageDir) == false {
		err :=os.MkdirAll(PackageDir,os.ModePerm)
		if err != nil{
			fmt.Println(ErrPermission)
		}
	}
	 _, err = os.Stat(ConfigFile)
	if err !=nil{
		f,err :=os.Create(ConfigFile)
		defer f.Close()
		if err != nil{
			errors.New("create config failed")
		}
		f.WriteString(defaultConfig)
		f.Sync()
	}
}


func Home() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support
	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}


func PathExists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	if err != nil{
		return false
	}
	return false
}



func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var dir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				dir = fpath[:lastIndex]
			}

			err = os.MkdirAll(dir, f.Mode())
			if err != nil {
				log.Fatal(err)
				return err
			}
			f1, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			_, err = io.Copy(f1, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func Unpack() (string)  {
	return InstallDir
}

func DownloadFile(url string,destinationFile string) (err error) {

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		log.Printf("[INFO] Dowoload to: [%s]", destinationFile)
		fmt.Print("\n")
		var source io.Reader
		downFile, err := os.Create(destinationFile)
		if err != nil {
			fmt.Printf("Can't create %s: %v\n", downFile, err)
			return err
		}
		defer downFile.Close()

		i, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
		sourceSiz := int64(i)
		bar := pb.New(int(sourceSiz))
		bar.SetRefreshRate(time.Microsecond * 10)
		bar.Start()
		source = resp.Body
		reader :=bar.NewProxyReader(source)
		code,err := io.Copy(downFile, reader)
		if err != nil{
			return  err
		}
		fmt.Println(code)
		bar.Finish()

		fmt.Print("\n")
		log.Printf("[INFO] [%s]Download Success.", destinationFile)
	} else {
		fmt.Print("\n")
		log.Printf("[ERROR] [%s]Download Failed,%s.", destinationFile, resp.Status)
	}
	return nil
}

func Unpackrelease(packageName string,version string) error  {
	//packageName := download.GetTarFileName(version)
	src := PackageDir + packageName
	dist := InstallDir  + version
	err := UnTarGz(src,dist)
	if err !=nil{
		return err
	}
	return nil
}


func UnTarGz(srcFilePath string, destDirPath string) error {
	fmt.Println("UnTarGzing [" + srcFilePath + "]  to   ["+destDirPath+"]...")
	// Create destination directory
	os.Mkdir(destDirPath, os.ModePerm)

	fr, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer fr.Close()

	// Gzip reader
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return err
	}
	// Tar reader
	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		// drop go dir file list
		filePaths := strings.Join(strings.Split(hdr.Name,"/")[1:],"/")
		if hdr.Typeflag != tar.TypeDir {
			os.MkdirAll(destDirPath+"/"+path.Dir(filePaths), os.ModePerm)
			// Write data to file
			fw, err := os.Create(destDirPath + "/" + filePaths)
			if err != nil {
				return err
			}
			_, err = io.Copy(fw, tr)
			if err != nil {
				return err
			}
		}
	}
	return nil
}