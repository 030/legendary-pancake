package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"gopkg.in/yaml.v3"
)

type regexFindAndReplace struct {
	find, replace *regexp.Regexp
}

func validateRegex(b []byte, regex *regexp.Regexp) error {
	match := regex.Match(b)
	length := len(regex.FindSubmatch(b))
	if !match || length < 1 {
		return fmt.Errorf("it does not match: '%v' and '%d'", match, length)
	}
	return nil
}

func regexForFindingAndReplacing(b []byte, regex string) (regexFindAndReplace, error) {
	regexOptions := `(?m)`
	regexForFindingKeyValuesToBeSorted := regexp.MustCompile(regexOptions + "^" + regex)
	if err := validateRegex(b, regexForFindingKeyValuesToBeSorted); err != nil {
		return regexFindAndReplace{}, err
	}
	regexReplace := regexp.MustCompile(regexOptions + `(` + regex + `){1,}$`)
	if err := validateRegex(b, regexReplace); err != nil {
		return regexFindAndReplace{}, err
	}
	return regexFindAndReplace{regexForFindingKeyValuesToBeSorted, regexReplace}, nil
}

func sortElements(r *regexp.Regexp, b []byte) string {
	unsortedElements := r.FindAllStringSubmatch(string(b), -1)
	toBeSortedElements := []string{}
	for i := range unsortedElements {
		toBeSortedElements = append(toBeSortedElements, unsortedElements[i][0])
	}
	sort.Strings(toBeSortedElements)
	fmt.Println(toBeSortedElements)
	toBeSortedElementsString := ""
	for _, toBeSortedElement := range toBeSortedElements {
		toBeSortedElementsString = toBeSortedElementsString + "\n" + toBeSortedElement
	}

	return toBeSortedElementsString
}

func sortElementsInFileByKey(key, regex string) (string, error) {
	b, err := ioutil.ReadFile(filepath.Clean(filepath.Join("test", "data", key, "input.yaml")))
	if err != nil {
		return "", err
	}

	r, err := regexForFindingAndReplacing(b, regex)
	if err != nil {
		return "", err
	}

	s := r.replace.ReplaceAllString(string(b), sortElements(r.find, b))

	return s, nil
}

func writeToFile(key, s string) (errs []error) {
	f, err := os.Create(filepath.Join("test", "data", key, "actual.yaml"))
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	w := bufio.NewWriter(f)
	if _, err := w.WriteString(s); err != nil {
		errs = append(errs, err)
		return errs
	}
	if err := w.Flush(); err != nil {
		errs = append(errs, err)
		return errs
	}
	defer func() {
		if err := f.Close(); err != nil {
			errs = append(errs, err)
		}
	}()
	return nil
}

func main() {
	key := "matchLabels"
	regex := `\s+-\s` + key + `:\n\s+[a-z0-9\/\.\":\s-]+$`
	s, err := sortElementsInFileByKey(key, regex)
	if err != nil {
		log.Fatal(err)
	}
	errs := writeToFile(key, s)
	for _, err := range errs {
		if err != nil {
			log.Fatal(err)
		}
	}

	key = "resources"
	regex = `\s+-\s[a-zA-Z0-9\/\.\"*:\s-]+\s+\s` + key + `:\n\s+[a-zA-Z0-9\/\.\"*:\s-]+$`
	s, err = sortElementsInFileByKey(key, regex)
	if err != nil {
		log.Fatal(err)
	}
	errs = writeToFile(key, s)
	for _, err := range errs {
		if err != nil {
			log.Fatal(err)
		}
	}

	bla()
}

type Rule struct {
	ApiGroups []string `yaml:"apiGroups"`
	Resources []string `yaml:"resources"`
	Verbs     []string `yaml:"verbs"`
}
type T struct {
	Rules []Rule `yaml:"rules"`
}

func bla() []Rule {
	b, err := ioutil.ReadFile(filepath.Clean(filepath.Join("test", "data", "resources", "input.yaml")))
	if err != nil {
		log.Fatal(err)
	}

	t := T{}
	err = yaml.Unmarshal(b, &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	resources := []string{}
	for _, rule := range t.Rules {
		resources = append(resources, rule.Resources...)
	}
	sort.Strings(resources)
	fmt.Println(resources)
	for i, resource := range resources {
		t.Rules[i].Resources = []string{resource}
	}
	fmt.Println(t.Rules)

	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m dump:\n%s\n\n", string(d))

	errs := writeToFile("resources", "---\n"+string(d))
	for _, err := range errs {
		if err != nil {
			log.Fatal(err)
		}
	}
	return t.Rules
}
