package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)

	if len(os.Args) <= 1 {
		fmt.Println("please enter a valid subcommand.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		runCmd.Parse(os.Args[2:])
		arg := runCmd.Args()
		run(arg)
	case "child":
		runCmd.Parse(os.Args[2:])
		arg := runCmd.Args()
		child(arg)
	case "build":
		tag := buildCmd.String("tag", "", "Name of container image")
		path := buildCmd.String("path", "", "Path to ContainerFile")
		buildCmd.Parse(os.Args[2:])
		build(*tag, *path)
	default:
		fmt.Printf("invalid subcommand %s", os.Args[1])
		os.Exit(1)
	}
}

func run(args []string) {
	fmt.Printf("Running %v \n", args)

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func child(args []string) {
	fmt.Printf("Running from proc in namespace %v \n", args)

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := syscall.Sethostname([]byte("container"))
	if err != nil {
		panic(err)
	}

	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

func build(tag, path string) {
	fmt.Printf("running build with tag: %s and %s", tag, path)
}
