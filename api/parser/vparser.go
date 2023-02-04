package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
	Tag   string
	Build string
	Date  string
	Error bool
}

type OSInfo struct {
	Os           string
	Architecture string
}

type LanguageInfo struct {
	Name string
	Version string
}

type ParsedInfo struct {
	Name     string
	Label    string
	Version  Version
	Os       OSInfo
	Language LanguageInfo
}

var reLanguage = regexp.MustCompile(`(?P<name>[a-zA-Z]+)?-?(?P<version>[\d+.?]+)`)

func (p *ParsedInfo) String() string {
	return fmt.Sprintf("%v (%v) %v %v", p.Name, p.Version, p.Os, p.Language)
}

func ParseVersionString(clientType, clientVersion, osType, goVersion string) *ParsedInfo {
	var output ParsedInfo

	output.Name = strings.ToLower(clientType)
	output.Label = strings.ToLower(clientType)
	output.Version = parseVersion(clientVersion)
	output.Os = parseOS(osType)
	output.Language = parseLanguage(goVersion)

	return &output
}

func parseLanguage(input string) LanguageInfo {
	var languageInfo LanguageInfo
	if input == "" {
		return languageInfo
	}
	match := reLanguage.FindStringSubmatch(input)

	if len(match) > 0 {
		languageInfo.Name = strings.ToLower(match[reLanguage.SubexpIndex("name")])
		languageInfo.Version = match[reLanguage.SubexpIndex("version")]
	}

	return languageInfo
}

func parseVersion(input string) Version {
	var vers Version
	if input == "" {
		return vers
	}

	split := strings.Split(input, "-")
	split_length := len(split)
	switch len(split) {
	case 8:
		fallthrough
	case 7:
		fallthrough
	case 6:
		fallthrough
	case 5:
		vers.Date = split[split_length - 1]
		vers.Build = split[split_length - 2]
		vers.Tag = strings.Join(split[1:split_length - 3], "")
		vers.Major, vers.Minor, vers.Patch = parseVersionNumber(split[0])
	case 4:
		// Date
		vers.Date = split[3]
		fallthrough
	case 3:
		// Build
		vers.Build = split[2]
		fallthrough
	case 2:
		// Tag
		vers.Tag = split[1]
		fallthrough
	case 1:
		// Version
		vers.Major, vers.Minor, vers.Patch = parseVersionNumber(split[0])
	}

	if vers.Major == 0 && vers.Minor == 0 && vers.Patch == 0 {
		fmt.Println("Version string is invalid:", input)
		vers.Error = true
	}
	
	return vers
}

func parseVersionNumber(input string) (int, int, int) {
	// Version
	trimmed := strings.TrimLeft(input, "v")
	vSplit := strings.Split(trimmed, ".")
	var major, minor, patch int

	switch len(vSplit) {
	case 4:
		fallthrough
	case 3:
		patch, _ = strconv.Atoi(vSplit[2])
		fallthrough
	case 2:
		minor, _ = strconv.Atoi(vSplit[1])
		fallthrough
	case 1:
		major, _ = strconv.Atoi(vSplit[0])
	}

	return major, minor, patch
}

func parseOS(input string) OSInfo {
	var osInfo OSInfo
	if input == "" {
		return osInfo
	}

	split := strings.Split(input, "-")
	switch len(split) {
	case 2:
		osInfo.Architecture = strings.ToLower(split[1])
		fallthrough
	case 1:
		osInfo.Os = strings.ToLower(split[0])
	}
	return osInfo
}
