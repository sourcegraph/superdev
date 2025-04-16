package superdev

import (
	"fmt"
	"os"
)

// Execute initializes and runs the CLI application
func Execute() {
	// Initialize all commands
	InitCommands()
	
	// Execute the root command
	if err := GetRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}