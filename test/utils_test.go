package test

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"1994.arthasup/common"
)

func TestProfile(t *testing.T) {
	s := common.Profile()
	fmt.Printf("s: %v\n", s)
}

func TestDownload(t *testing.T) {
	a := make([]string, 1)
	b := append(a, "1")
	fmt.Printf("b: %v\n", b)
}
func Test(t *testing.T) {
	name, err := common.Download()
	if err != nil && name == "" {
		panic(err)
	}
	version, _ := common.Unzip(name)
	common.Alias(version)
}

func TestAppendFile(t *testing.T) {
	home, _ := common.ArthasHome()
	user, _ := os.UserHomeDir()
	path := filepath.Join(user, ".zshrc")
	f, _ := os.Open(path)

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		str := scanner.Text()
		if !strings.HasPrefix(str, "alias arthas") {
			lines = append(lines, scanner.Text())
		}
	}

	boot := filepath.Join(home, "3.5.4", "arthas-boot.jar")
	lines = append(lines, fmt.Sprintf("alias arthas='%s -jar %s'\n", common.Java, boot))

	err := ioutil.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)

	if err != nil {
		panic(err)
	}
}

func TestPre(t *testing.T) {
	home, _ := common.Home()
	arthas, _ := common.ArthasHome()
	fmt.Printf("start to cleanup arthas dir: %s and %s", home, arthas)
	os.RemoveAll(home)
	os.RemoveAll(arthas)
}

func TestSys(t *testing.T) {
	binary, _ := exec.LookPath("java")
	a, _ := common.ArthasHome()

	str := filepath.Join(a, "3.5.4", "arthas-boot.jar")
	fmt.Println(str)
	args := []string{"java", "-jar", str}
	syscall.Exec(binary, args, os.Environ())
}
