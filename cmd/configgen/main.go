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
)

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}

func main() {
	var configFile, outputFile, role string
	var getroles bool
	flag.StringVar(&configFile, "input", "conf/app.yaml", "Path to the config file")
	flag.StringVar(&outputFile, "output", "conf/", "Output file directory")
	flag.StringVar(&role, "role", "", "The role to generate a config for")
	flag.BoolVar(&getroles, "getroles", false, "If true, will print available roles and exit")
	flag.Parse()

	if ok, err := exists(configFile); !ok {
		fmt.Printf("No file found at %v. Make sure you are in the context of a deployable application.\n", configFile)
		if err != nil {
			fmt.Println("Additionally, the following error occurred:")
			fmt.Println(err.Error())
		}
		os.Exit(0);
	}

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Unable to load app.yaml: ", err.Error())
	}
	configMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal(file, &configMap)
	if err != nil {
		log.Fatalf("Unable to deserialize app.yaml: ", err.Error())
	}

	var roleMap, baseMap map[interface{}]interface{}

	roleMap = configMap["envs"].(map[interface{}]interface{})

	if getroles {
		availableRoles := getRoles(roleMap)
		fmt.Println(strings.Join(availableRoles, "\n"))
		os.Exit(0)
	}

	baseMap = configMap["base-app-config"].(map[interface{}]interface{})

	roleMapReflectValue := reflect.ValueOf(roleMap).MapIndex(reflect.ValueOf(role))
	if !roleMapReflectValue.IsValid() {
		log.Fatalf("Role %v was not found in app.yaml", role)
	}

	roleEnvironmentMap := roleMapReflectValue.Interface().(map[interface{}]interface{})["config-override"]
	baseMapReflectValue := reflect.ValueOf(baseMap).Interface()

	mp := merge(roleEnvironmentMap, baseMapReflectValue)

	outputFilename := outputFile + role + ".yaml"
	yamlBytes, err := yaml.Marshal(mp)
	if err != nil {
		log.Fatalf("Error marshalling yaml output: ", err.Error())
	}
	outF, err := os.Create(outputFilename)
	if err != nil {
		panic(err)
	}
	outF.Write(yamlBytes)
}

func getRoles(roleMap map[interface{}]interface{}) []string {
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
	return outmap
}
