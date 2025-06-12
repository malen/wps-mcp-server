package main

import (
	"aoisoft/wps-mcp-server/internal/mcps"
)

func main() {
	s := mcps.New()

	s.Start()
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Failed to start the wps mcp server:%v \n", err)
	// 	os.Exit(1)
	// }
}
