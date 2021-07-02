package structures

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// SaveTriangleScene saves a triangle scene to a json file with the specified name.
func SaveTriangleScene(trs *TriangleScene, name string) {
	data, err := json.Marshal(trs)

	if err != nil {
		fmt.Println("Marshalling failed")
		return
	}

	file, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	file.Write(data)

	fmt.Println("Scene saved")
}

// LoadTriangleScene loads a triangle scene from a json file given its name.
func LoadTriangleScene(name string) *TriangleScene {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var trs TriangleScene
	json.Unmarshal(data, &trs)

	return &trs
}
