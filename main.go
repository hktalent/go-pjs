package main

import (
	"github.com/hktalent/go-pjs/pkg"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	os.Args = []string{"", "/Users/51pwn/MyWork/TestPoc/CVE_2021_2394_Payload.dat"}
	if data, err := ioutil.ReadFile(os.Args[1]); nil == err {
		if c, err := pkg.ParseSerializedObject(data); nil == err {
			log.Println(c)
		} else {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}

}
