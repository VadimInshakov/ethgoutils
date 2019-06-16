package contract

import (
	"fmt"
	"regexp"
	"os/exec"
	"log"
)

func Compile(contractPath string) {

	r, _ := regexp.Compile("([/\\w-/]+/)([\\w-]+\\.)")
	arr := r.FindStringSubmatch(contractPath)

	contractName := arr[2]
	contractName = string(contractName[:len(contractName)-1])

	// create abi
	abi := exec.Command("solc", "--abi", contractPath, "-o", "build")
	output, err := abi.CombinedOutput()
	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + string(output))
		return
	}
	fmt.Println(string(output))
	abi.Wait()

	// create go package
	gopack := exec.Command("abigen", fmt.Sprintf("--abi=./build/%s.abi", contractName), fmt.Sprintf("--pkg=%s", contractName), fmt.Sprintf("--out=./build/%s.go", contractName))
	output, err = gopack.CombinedOutput()
	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + string(output))
		return
	}
	fmt.Println(string(output))
	gopack.Wait()

	// create bin
	bin := exec.Command("solc", "--bin", contractPath, "-o", "build")
	output, err = bin.CombinedOutput()
	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + string(output))
		return
	}
	fmt.Println(string(output))
	bin.Wait()

	// create go package with deploy method
	godeploy := exec.Command("abigen", fmt.Sprintf("--bin=./build/%s.bin", contractName), fmt.Sprintf("--abi=./build/%s.abi", contractName), fmt.Sprintf("--pkg=%s", contractName), fmt.Sprintf("--out=./build/%s.go", contractName))
	output, err = godeploy.CombinedOutput()
	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + string(output))
		return
	}
	fmt.Println(string(output))
	godeploy.Wait()
}