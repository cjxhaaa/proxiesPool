package config

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	SECTCRE *regexp.Regexp = regexp.MustCompile(`\[(?P<header>[^]]+)]`)
	OPTCRE *regexp.Regexp = regexp.MustCompile(`(?P<option>.*?)\s*[=:]\s*(?P<value>.*)$`)
	//NONSPACECRE *regexp.Regexp = regexp.MustCompile(`\S`)
	BOOLEAN_STATES map[string]bool = map[string]bool{"1":true,"yes":true,"true":true,"on":true,
													 "0":false,"no":false,"false":false,"off":false}
)


type configObj struct {
		sections      map[string]map[string]string
		sectionBefore string
	}

func reCompileGroup(rexp *regexp.Regexp,matchStr string) map[string]string {
	matchs := rexp.FindStringSubmatch(matchStr)
	if len(matchs) == 0 {
		return nil
	}

	groupNames := rexp.SubexpNames()
	result := make(map[string]string)
	for i,name := range groupNames {
		if i != 0 && name != "" {
			result[name] = matchs[i]
		}
	}
	return result
}

func New(filename string) *configObj {
	var configInfo configObj
	configInfo.sections = make(map[string]map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		lineValue := strings.TrimSpace(scanner.Text())
		if lineValue == "" {
			continue
		}

		sectname := reCompileGroup(SECTCRE,lineValue)
		if sectname != nil {
			header := sectname["header"]

			configInfo.sectionBefore = header

		} else {
			optname := reCompileGroup(OPTCRE,lineValue)
			if optname != nil {
				if configInfo.sections[configInfo.sectionBefore] == nil {
					configInfo.sections[configInfo.sectionBefore] = make(map[string]string)
				}
				configInfo.sections[configInfo.sectionBefore][optname["option"]] = optname["value"]
			}
		}

	}
	return &configInfo
}

func (co *configObj) Set(section string, option string, value string) {
	if co.sections[section] == nil {
		co.sections[section] = make(map[string]string)
	}
	co.sections[section][option] = value
	co.sectionBefore = section
}

func (co *configObj) GetSection(section string) map[string]string {
	if v,ok := co.sections[section];ok {
		return v
	}
	return nil
}

func (co *configObj) HasSection(section string) bool {
	if _,ok := co.sections[section];ok {
		return true
	}
	return false
}

func (co *configObj) GetOptionValue(section string, option string) string {
	options := co.GetSection(section)
	if options != nil {
		return options[option]
	}
	return ""
}

func (co *configObj) GetOptionBool(section string, option string) bool {
	options := co.GetSection(section)
	if options != nil {
		return BOOLEAN_STATES[options[option]]
	}
	return false
}
