package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("C:\\Python\\python.exe", "C:\\Learning\\Py_Code\\Go_use\\1.py", "1", "2")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
	}
	fmt.Println(string(output))
}
