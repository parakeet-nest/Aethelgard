package main

import (
	"encoding/json"
	"strconv"
	"github.com/extism/go-pdk"
)

type Arguments struct {
	From int `json:"from"`
	To int `json:"to"`
}

//export move
func move() {
	arguments := pdk.InputString()

	var args Arguments
	json.Unmarshal([]byte(arguments), &args)
	
	pdk.OutputString("👣 from " + strconv.Itoa(args.From) + " to " + strconv.Itoa(args.To))

}

func main() {}
