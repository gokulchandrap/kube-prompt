package kube

import (
	"io/ioutil"
	"log"
	"strings"

	"path/filepath"

	"github.com/c-bata/go-prompt"
)

func init() {
	fileListCache = map[string][]prompt.Suggest{}
}

func getPreviousOption(d prompt.Document) (option string, found bool) {
	args := strings.Split(d.TextBeforeCursor(), " ")
	l := len(args)
	if l >= 2 {
		option = args[l-2]
	}
	if strings.HasPrefix(option, "-") {
		return option, true
	}
	return "", false
}

func completeOptionArguments(d prompt.Document) ([]prompt.Suggest, bool) {
	option, found := getPreviousOption(d)
	if !found {
		return []prompt.Suggest{}, false
	}

	switch option {
	case "-f", "--filename":
		return fileCompleter(d), true
	default:
		return []prompt.Suggest{}, false
	}
}

/* file list */

var fileListCache map[string][]prompt.Suggest

func fileCompleter(d prompt.Document) []prompt.Suggest {
	path := d.GetWordBeforeCursor()
	dir := filepath.Dir(path)
	if cached, ok := fileListCache[dir]; ok {
		return cached
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Print("[ERROR] catch error " + err.Error())
		return []prompt.Suggest{}
	}
	suggests := make([]prompt.Suggest, 0, len(files))
	for _, f := range files {
		if !f.IsDir() &&
			!strings.HasSuffix(f.Name(), ".yml") &&
			!strings.HasSuffix(f.Name(), ".yaml") {
			continue
		}
		suggests = append(suggests, prompt.Suggest{Text: filepath.Join(dir, f.Name())})
	}
	return prompt.FilterHasPrefix(suggests, path, false)
}
