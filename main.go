package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

// Creating container from scratch using go.

// Program Synopsis
// docker run image <cmd> <params>

// go program

// go run main.go run <cmds> <params>

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	case "default":
		panic("Invalid Command")
	}

}

func run() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	cmd.Run()

}

// Run function for arguments and running commands
func child() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cg()

	syscall.Sethostname([]byte("container"))
	syscall.Chroot("/")
	syscall.Chdir("/")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

}

func cg() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	os.Mkdir(filepath.Join(pids, "virus333"), 0755)
	must(ioutil.WriteFile(filepath.Join(pids, "virus333/pids.max"), []byte("20"), 0700))
	// Removes the new cgroup in place after the container exits
	must(ioutil.WriteFile(filepath.Join(pids, "virus333/notify_on_release"), []byte("1"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids, "virus333/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
