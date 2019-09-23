package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"github.com/docopt/docopt-go"
)

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func main() {
	// config stuff
	usage := `YAML Merge

Usage:
	yamlmerge <input> <base> <override> [--get-roots]
	yamlmerge -h | --help
	yamlmerge --version

Options:
	<input>	The input yaml file
	<base> The name of the base node whose values you'll be overriding.
	<override> The root node containing values that will override those of the base node
	--get-roots	Will print all available root-level nodes to the console and then exit normally
	-h --help	Show this screen.
	--version	Show version.
`

	arguments, _ := docopt.ParseDoc(usage)
	var config struct {
		Input    string `docopt:"<input>"`
		Base     string `docopt:"<base>"`
		Override string `docopt:"<override>"`
		GetRoles bool   `docopt:"--get-roles"`
	}
	arguments.Bind(&config)

	inputFile := config.Input
	defaultRole := config.Base
	role := config.Override
	getroles := config.GetRoles

	// define defaults for empty values
	if ok, err := exists(inputFile); !ok {
		fmt.Printf("No file found at: %s", inputFile)
		if err != nil {
			fmt.Println("Additionally, the following error occurred:")
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}
	if role == "" {
		log.Fatal("No role specified")
	}

	// read the input file
	file, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Unable to load file: %s\n%v", inputFile, err)
	}
	configMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal(file, &configMap)
	if err != nil {
		log.Fatalf("Unable to deserialize file: %s\n%v", inputFile, err)
	}

	var roleMap, baseMap map[interface{}]interface{}

	roleMap = configMap
	//delete(roleMap, defaultRole) // take out the default

	if getroles {
		availableRoles := getStringKeys(roleMap)
		fmt.Println(strings.Join(availableRoles, "\n"))
		os.Exit(0)
	}

	baseMap = configMap[defaultRole].(map[interface{}]interface{})

	roleMapReflectValue := reflect.ValueOf(roleMap).MapIndex(reflect.ValueOf(role))
	if !roleMapReflectValue.IsValid() {
		log.Fatalf("Role %v was not found in input file %s", role, inputFile)
	}

	roleEnvironmentMap := roleMapReflectValue.Interface().(map[interface{}]interface{})
	baseMapReflectValue := reflect.ValueOf(baseMap).Interface()

	mp := merge(roleEnvironmentMap, baseMapReflectValue)

	// write the output
	yamlBytes, err := yaml.Marshal(mp)
	if err != nil {
		log.Fatalf("Error marshalling yaml output: ", err.Error())
	}
	fmt.Println(string(yamlBytes))
}

// getStringKeys returns all keys in the map provided as strings
func getStringKeys(roleMap map[interface{}]interface{}) []string {
	roles := make([]string, len(roleMap))
	i := 0
	for k := range roleMap {
		var kString string
		kString = k.(string)
		roles[i] = kString
		i++
	}
	return roles
}

// recursively merges two maps; assumes the second map is a superset of the first
func merge(role interface{}, app interface{}) map[interface{}]interface{} {
	outmap := make(map[interface{}]interface{})

	roleValue := reflect.ValueOf(role).Interface().(map[interface{}]interface{})
	appValue := reflect.ValueOf(app).Interface().(map[interface{}]interface{})

	// for all in role that are also in app, recur downward and replace crap
	for k, v := range appValue {
		roleMapValue, ok := roleValue[k]
		if !ok {
			outmap[k] = v
		} else if reflect.ValueOf(v).Kind() == reflect.Map {
			outmap[k] = merge(roleMapValue, v)
		} else {
			outmap[k] = roleMapValue
		}
	}

	// for all in role that are not also in app, just take the whole tree
	for k, v := range roleValue {
		_, ok := appValue[k]
		if !ok {
			outmap[k] = v
		}
	}
	return outmap
}
