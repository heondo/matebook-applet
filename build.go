// Copyright (C) 2019 Evgeny Kuznetsov (evgeny@kuznetsov.md)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along tihe this program. If not, see <https://www.gnu.org/licenses/>.

// +build ignore

package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

type packFile struct {
	src string
	dst string
	mod os.FileMode
}

var (
	filename  string
	packFiles = []packFile{
		{src: "LICENSE", dst: "LICENSE", mod: 0644},
		{src: "README.md", dst: "README.md", mod: 0644},
		{src: "huawei-wmi.rules", dst: "huawei-wmi.rules", mod: 0644},
	}
)

func main() {
	version := getVersion()
	btime := buildTime()
	fmt.Printf("Building version %s\n", version)
	fmt.Println("Building as of", time.Unix(btime, 0))
	buildAssets(btime)
	buildBinary(btime)
	filename = "matebook-applet-amd64" + "-" + version
	signFile("matebook-applet", "466F4F38E60211B0")
	buildTar()
	fmt.Println("archive", filename, "created")
}

func buildBinary(t int64) {
	cmd := exec.Command("go", "build")
	if err := cmd.Run(); err != nil {
		log.Fatalln("failed to build binary")
	}
	setFileTime("matebook-applet", t)
	packFiles = append(packFiles, packFile{"matebook-applet", "matebook-applet", 0755})
}

func buildAssets(t int64) {
	cmd := exec.Command("go", "run", "assets_generate.go")
	if err := cmd.Run(); err != nil {
		log.Fatalln("failed to rebuild assets")
	}
	setFileTime("assets.go", t)
}

func buildTar() {
	for i := range packFiles {
		packFiles[i].dst = filename + "/" + packFiles[i].dst
	}
	filename = filename + ".tar.gz"
	fd, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	gw, err := gzip.NewWriterLevel(fd, gzip.BestCompression)
	if err != nil {
		log.Fatal(err)
	}
	tw := tar.NewWriter(gw)
	for _, f := range packFiles {
		sf, err := os.Open(f.src)
		if err != nil {
			log.Fatal(err)
		}
		info, err := sf.Stat()
		if err != nil {
			log.Fatal(err)
		}
		h := &tar.Header{
			Name:    f.dst,
			Size:    info.Size(),
			Mode:    int64(f.mod),
			ModTime: info.ModTime(),
		}
		err = tw.WriteHeader(h)
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(tw, sf)
		if err != nil {
			log.Fatal(err)
		}
		sf.Close()
	}
	err = tw.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = gw.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = fd.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func setFileTime(f string, t int64) {
	cmd := exec.Command("touch", "-t", fmt.Sprint(time.Unix(t, 0).Format("200601021504.05")), f)
	if err := cmd.Run(); err != nil {
		log.Fatalln("failed to set time on", f)
	}
}

func signFile(f string, k string) {
	cmd := exec.Command("gpg", "--detach-sign", "--yes", "-u", k, f)
	if err := cmd.Run(); err != nil {
		fmt.Println("signing", f, "failed")
		filename = filename + "-unsigned"
	} else {
		fmt.Println(f, "successfully signed with key", k)
		packFiles = append(packFiles, packFile{"matebook-applet.sig", "matebook-applet.sig", 0644})
	}
}

func getVersion() string {
	s, err := getString("git", "describe", "--always", "--dirty")
	versionRe := regexp.MustCompile(`^v[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}(\-[0-9]{1,3}\-g[0-9a-f]{5,15})?`)
	if err == nil {
		if versionRe.MatchString(s) {
			return s
		}
	}
	return "unknown"
}

func getString(c string, a ...string) (string, error) {
	cmd := exec.Command(c, a...)
	b, err := cmd.CombinedOutput()
	return string(bytes.TrimSpace(b)), err
}

func buildTime() int64 {
	s, err := getString("git", "show", "-s", "--format=%ct")
	if err == nil {
		if i, e := strconv.ParseInt(s, 10, 64); e == nil {
			return i
		}
	}
	return time.Now().Unix()
}