package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}

func main() {
	args := os.Args[1:]
	if len(args) <= 2 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Println("Generate collapsed yaml based on app.yaml and role")
		os.Exit(0)
	}
	if len(args) != 2 {
		print("Usage: roleout <role> <output_path>")
		os.Exit(0)
	}

	appyaml := "conf/app.yaml"
	role := args[0]
	outputFile := args[1]

	if ok, err := exists(appyaml); !ok {
		fmt.Println("No file found at conf/app.yaml. Make sure you are in the context of a deployable application.")
		if err != nil {
			fmt.Println("Additionally, the following error occurred:")
			fmt.Println(err.Error())
		}
		os.Exit(0);
	}

	file, err := ioutil.ReadFile(appyaml)
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
