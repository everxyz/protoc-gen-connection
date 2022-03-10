package tools

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"strings"
)

func GetPrefix(names []string) string {
	if 0 == len(names) {
		return ""
	}
	// Find word with minimum length
	short := names[0]
	for _, s := range names {
		if len(short) >= len(s) {
			short = s
		}
	}
	// Loop over minword length: from one character prefix to full minword
	prefixArray := []string{}
	prefix := ""
	oldPrefix := ""
	for i := 0; i < len(short); i++ {
		// https://hermanschaaf.com/efficient-string-concatenation-in-go/
		prefixArray = append(prefixArray, string(short[i]))
		prefix = strings.Join(prefixArray, "")
		// Sub loop check all elements start with the prefix
		for _, s := range names {
			// https://gist.github.com/lbvf50mobile/65298c689e2f2b850aa6ad8bd7b61717
			if !strings.HasPrefix(s, prefix) {
				return oldPrefix
			}
		}
		oldPrefix = prefix
	}

	return prefix
}

type YamlFunc func(path string)

func ParseYaml(ins interface{}) YamlFunc {
	return func(path string) {
		yamlFile, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}

		err = yaml.Unmarshal(yamlFile, ins)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
	}
}
