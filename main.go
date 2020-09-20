package main

import (
	"github.com/lainio/findcfg/findcfg"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(findcfg.Analyzer) }
