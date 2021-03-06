package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8yaml "k8s.io/apimachinery/pkg/util/yaml"
)

func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func listYamlFiles(root string) ([]string, error) {
	var files []string

	path, err := filepath.Abs(root)
	if err != nil {
		return files, err
	}

	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		filename := file.Name()
		if filepath.Ext(filename) == ".yaml" || filepath.Ext(filename) == ".yml" {
			files = append(files, filepath.Join(path, filename))
		}
	}

	return files, nil
}

func readFilesAsManifests(paths []string) (result []unstructured.Unstructured, errs []error) {
	for _, path := range paths {
		rawdata, err := ioutil.ReadFile(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("could not read YAML: %s from disk: %s", path, err))
		}
		manifest, err := readManifestData(bytes.NewReader(rawdata))
		if err != nil {
			errs = append(errs, fmt.Errorf("could not read YAML: %s from disk: %s", path, err))
		}
		result = append(result, manifest...)
	}

	return result, errs
}

func readManifestData(yamlData io.Reader) ([]unstructured.Unstructured, error) {
	decoder := k8yaml.NewYAMLToJSONDecoder(yamlData)

	var manifests []unstructured.Unstructured
	for {
		nxtManifest := unstructured.Unstructured{}
		err := decoder.Decode(&nxtManifest)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		manifests = append(manifests, nxtManifest)
	}

	return manifests, nil
}
