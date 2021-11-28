package common

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/xerrors"
)

const (
	ArthasUp    = ".arthasup"
	Arthas      = ".arthas"
	DownloadUrl = "https://arthas.aliyun.com/download/latest_version?mirror=aliyun"
)

func Home() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", xerrors.Errorf("get user home error, %w", err)
	}
	return filepath.Join(home, ArthasUp), nil
}

func ArthasHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", xerrors.Errorf("get user home error, %w", err)
	}
	return filepath.Join(home, Arthas), nil
}

func Download() (string, error) {
	get, err := http.Get(DownloadUrl)
	if err != nil {
		return "", xerrors.Errorf("download error, %w", err)
	}

	if get.StatusCode != 200 {
		return "", xerrors.New("download error, please make sure your network available")
	}

	newPath := get.Request.URL.String()

	home, err := Home()
	if err != nil {
		return "", xerrors.Errorf("download error, %w", err)
	}

	name := filepath.Join(home, fileName(newPath))

	if exist(name) {
		return name, xerrors.New("file exist, just return")
	}

	response, err := http.Get(newPath)

	if err != nil {
		return "", xerrors.Errorf("download error, %w", err)
	}

	if response.StatusCode != 200 {
		return "", xerrors.New("download error, please make sure your network available")
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	err = os.MkdirAll(home, os.ModePerm)
	if err != nil {
		return "", xerrors.Errorf("download error, %w", err)
	}

	create, err := os.Create(name)
	if err != nil {
		return "", xerrors.Errorf("download error, %w", err)
	}

	defer func(create *os.File) {
		_ = create.Close()
	}(create)

	_, _ = io.Copy(create, response.Body)

	return name, nil
}

func Unzip(file string) (string, error) {
	fmt.Println(file)
	version := strings.Split(filepath.Base(file), "-")[2]
	home, err := ArthasHome()
	if err != nil {
		return "", xerrors.Errorf("unzip error, %w", err)
	}
	dst := filepath.Join(home, version)
	_ = os.MkdirAll(dst, os.ModePerm)

	archive, err := zip.OpenReader(file)
	if err != nil {
		panic(err)
	}
	defer func(archive *zip.ReadCloser) {
		_ = archive.Close()
	}(archive)

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return "", xerrors.New("invalid file path")
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			_ = os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return "", xerrors.New("invalid file path")
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", xerrors.New("invalid file path")
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return "", xerrors.New("invalid file path")
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return "", xerrors.New("invalid file path")
		}

		_ = dstFile.Close()
		_ = fileInArchive.Close()
	}
	return version, nil
}

func fileName(urlPath string) string {
	p := path.Base(urlPath)
	if strings.Contains(p, "?") {
		pathSlice := strings.Split(p, "?")
		return pathSlice[0]
	} else {
		return p
	}
}

func exist(p string) bool {
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func Alias(version string) {
	os := runtime.GOOS
	//  todo
	match := strings.HasPrefix(os, "linux") || strings.HasPrefix(os, "darwin")

	if !match {
		return
	}

	p := Profile()
	if p == "" {
		return
	}

	home, _ := ArthasHome()

	boot := filepath.Join(home, version, "arthas-boot.jar")

	AppendStringToFile(p, fmt.Sprintf("alias arthas='java -jar %s'", boot))
	exec.Command("source", p).Run()
}

func Profile() string {
	c := os.ExpandEnv("$SHELL")
	home, _ := os.UserHomeDir()

	var result string

	switch {
	case strings.ContainsAny(c, "zsh"):
		result = filepath.Join(home, ".zshrc")
	case strings.ContainsAny(c, "bash"):
		result = filepath.Join(home, ".bashrc")
	default:
		result = ""
	}
	return result
}

func AppendStringToFile(path, text string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}
