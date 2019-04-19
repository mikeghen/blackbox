package apt

import (
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora"

	"github.com/go-cmd/cmd"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// CachePolicy just holds the output of `apt-cache policy x` for serialization
type CachePolicy struct {
	App struct {
		Installed string `yaml:"Installed"`
		Candidate string `yaml:"Candidate"`
	} `yaml:"blackboxd"`
}

// Package gives some simple encapsulated access
type PackageInfo struct {
	Name       string `json:"name"`
	Installed  string `json:"installed"`
	Candidate  string `json:"candidate"`
	Updateable bool   `json:"updateable"`
}

// IsInstalled tests to see if the package is installed. It should be.
func (p *PackageInfo) IsInstalled() bool {
	return p.Installed != "(none)"
}

func (p *PackageInfo) UpgradeAvailable() bool {
	return p.Installed != p.Candidate
}

// apt-get install --only-upgrade <packagename>
func Upgrade(name string) (*cmd.Status, error) {
	info, err := GetPackageInfo(name)
	if err != nil {
		return nil, err
	}

	if !info.UpgradeAvailable() {
		return nil, errors.New("no upgrade available")
	}

	var cacheCmd = cmd.NewCmdOptions(cmd.Options{Buffered: true}, "apt-get", "install", "--only-upgrade", name)
	status := <-cacheCmd.Start()
	if status.Exit != 0 {
		return nil, errors.Wrap(status.Error, "bad command")
	}

	log(status)

	return &status, nil
}

func GetPackageInfo(name string) (*PackageInfo, error) {
	updateCmd := cmd.NewCmdOptions(cmd.Options{Buffered: true}, "apt-get", "update")
	status := <-updateCmd.Start()
	if status.Exit != 0 {
		return nil, errors.Wrap(status.Error, "bad command")
	}

	log(status)

	var cacheCmd = cmd.NewCmdOptions(cmd.Options{Buffered: true}, "apt-cache", "policy", name)
	status = <-cacheCmd.Start()
	if status.Exit != 0 {
		return nil, errors.Wrap(status.Error, "bad command")
	}

	log(status)

	out := new(CachePolicy)
	err := yaml.Unmarshal([]byte(strings.Join(status.Stdout, "\n")), out)
	if err != nil {
		return nil, err
	}

	info := &PackageInfo{
		Installed:  out.App.Installed,
		Candidate:  out.App.Candidate,
		Updateable: out.App.Installed != out.App.Candidate,
	}

	return info, nil
}

func log(status cmd.Status) {
	fmt.Println(aurora.Blue(strings.Join(status.Stdout, "\n")))
}
