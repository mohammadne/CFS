//go:build linux
// +build linux

package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func main() {
	// Hidden child mode: handled only internally
	if os.Getenv("DOCKER_CHILD") == "1" {
		child()
		return
	}

	// Normal mode
	if len(os.Args) < 2 {
		panic("usage: pull | run")
	}

	switch os.Args[1] {
	case "pull":
		if len(os.Args) < 3 {
			panic("usage: pull <name>")
		}
		pull(os.Args[2])

	case "run":
		if len(os.Args) < 4 {
			panic("usage: run <image> <command> [args...]")
		}
		run()

	default:
		panic("unknown command")
	}
}

func pull(image string) {
	url := "https://dl-cdn.alpinelinux.org/alpine/v3.23/releases/x86_64/alpine-minirootfs-3.23.0-x86_64.tar.gz"

	log.Printf("Pulling %s from %s\n", image, url)

	resp, err := http.Get(url)
	must(err)
	defer resp.Body.Close()

	// Prepare image directory
	rootfsPath := filepath.Join("images", image)
	must(os.MkdirAll(rootfsPath, 0755))

	// Extract tar.gz
	gzReader, err := gzip.NewReader(resp.Body)
	must(err)

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		must(err)

		target := filepath.Join(rootfsPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			must(os.MkdirAll(target, os.FileMode(header.Mode)))

		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			must(err)
			_, err = io.Copy(f, tarReader)
			f.Close()
			must(err)

		case tar.TypeSymlink:
			must(os.Symlink(header.Linkname, target))
		}
	}

	log.Printf("Image %s pulled to %s\n", image, rootfsPath)
}

func run() {
	image := os.Args[2]
	rootfs := filepath.Join("images", image)

	log.Printf("Running container using image %s\n", image)

	// Prepare re-exec command
	cmd := exec.Command(
		"/proc/self/exe",
		append([]string{rootfs}, os.Args[3:]...)..., // new argv: [rootfs, cmd, args...]
	)

	// Hide child mode behind env var
	cmd.Env = append(os.Environ(), "DOCKER_CHILD=1")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	must(cmd.Run())
}

func child() {
	rootfs := os.Args[1]
	entrypoint := os.Args[2]
	args := os.Args[3:]

	log.Printf("Child: starting container with rootfs=%s command=%s", rootfs, entrypoint)

	must(syscall.Sethostname([]byte("container")))

	// Change root filesystem
	must(syscall.Chroot(rootfs))
	must(syscall.Chdir("/"))

	// Mount proc
	must(os.MkdirAll("/proc", 0555))
	must(syscall.Mount("proc", "/proc", "proc", 0, ""))

	// Run actual user command
	cmd := exec.Command(entrypoint, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	// Unmount proc
	_ = syscall.Unmount("/proc", 0)

	must(err)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
