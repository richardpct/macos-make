// make package
package main

import (
	"flag"
	"fmt"
	"github.com/richardpct/pkgsrc"
	"log"
	"os"
	"os/exec"
	"path"
)

var destdir = flag.String("destdir", "", "directory installation")
var pkg pkgsrc.Pkg

const (
	name     = "make"
	vers     = "4.2"
	ext      = "tar.gz"
	url      = "http://ftp.gnu.org/gnu/make"
	hashType = "sha256"
	hash     = "e968ce3c57ad39a593a92339e23eb148af6296b9f40aa453a9a9202c99d34436"
)

func checkArgs() error {
	if *destdir == "" {
		return fmt.Errorf("Argument destdir is missing")
	}
	return nil
}

func configure() {
	fmt.Println("Waiting while configuring ...")
	cmd := exec.Command("./configure",
		"--prefix="+*destdir)
	path := os.Getenv("PATH")
	cmd.Env = append(os.Environ(), "PATH="+path+":"+*destdir+"/bin")
	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", out)
	}
}

func build() {
	fmt.Println("Waiting while compiling ...")
	cmd := exec.Command("make", "-j"+pkgsrc.Ncpu)
	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", out)
	}
}

func install() {
	fmt.Println("Waiting while installing ...")
	cmd := exec.Command("make", "install")
	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", out)
	}
}

func main() {
	flag.Parse()
	if err := checkArgs(); err != nil {
		log.Fatal(err)
	}

	pkg.Init(name, vers, ext, url, hashType, hash)
	pkg.CleanWorkdir()
	if !pkg.CheckSum() {
		pkg.DownloadPkg()
	}
	if !pkg.CheckSum() {
		log.Fatal("Package is corrupted")
	}

	pkg.Unpack()
	wdPkgName := path.Join(pkgsrc.Workdir, pkg.PkgName)
	if err := os.Chdir(wdPkgName); err != nil {
		log.Fatal(err)
	}
	configure()
	build()
	install()
}
