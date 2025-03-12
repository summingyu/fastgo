package version

import (
	"fmt"
	"os"
	"strconv"

	flag "github.com/spf13/pflag"
)

type versionValue int

const (
	VersionNotSet  versionValue = 0
	VersionEnabled versionValue = 1
	VersionRaw     versionValue = 2
	VersionJson    versionValue = 3
)
const (
	strRawVersion  string = "raw"
	strJsonVersion string = "json"
)

func (v *versionValue) IsBoolFlag() bool {
	return true
}

func (v *versionValue) Get() any {
	return *v
}

func (v *versionValue) Set(s string) error {
	if s == strRawVersion {
		*v = VersionRaw
		return nil
	} else if s == strJsonVersion {
		*v = VersionJson
		return nil
	}
	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = VersionEnabled
	} else {
		*v = VersionNotSet
	}
	return err
}

func (v *versionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion
	} else if *v == VersionJson {
		return strJsonVersion
	}
	return fmt.Sprintf("%v", bool(*v == VersionEnabled))
}

func (v *versionValue) Type() string {
	return "version"
}

func VersionVar(p *versionValue, name string, value versionValue, usage string) {
	*p = value
	flag.Var(p, name, usage)
	flag.Lookup(name).NoOptDefVal = "true"
}

func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)
	VersionVar(p, name, value, usage)
	return p
}

const versionFlagName = "version"

var versionFlag = Version(versionFlagName, VersionNotSet, "Print version information and quit")

func AddFlags(fs *flag.FlagSet) {
	fs.AddFlag(flag.Lookup(versionFlagName))
}

func PrintAndExitIfRequested() {
	// 检查版本标志的值并打印相应的信息
	if *versionFlag == VersionRaw {
		fmt.Printf("%s\n", Get().Text())
		os.Exit(0)
	} else if *versionFlag == VersionJson {
		fmt.Printf("%s\n", Get().ToJSON())
		os.Exit(0)
	} else if *versionFlag == VersionEnabled {
		fmt.Printf("%s\n", Get().String())
		os.Exit(0)
	}
}
