package main

import (
	"flag"
	"fmt"
	"github.com/MasenkoHa/s3discover"
)

func main() {
	flag.Parse()
	domains := flag.Args()
	fmt.Println(s3discover.S3discover(domains))

}