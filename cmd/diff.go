package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/davidhuangdw/goJsonDiff/pkg/diff"
	"github.com/davidhuangdw/goJsonDiff/pkg/render"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	fFormat := flag.String("f", "console", "Format: console, html, patch, ...")
	fDebug := flag.Bool("debug", false, "debug bool: show timing & extra info")

	flag.Usage = func() {
		info(fmt.Sprintf("Usage: %s [OPTIONS] {from_json_file_name} {to_json_file_name} \n", os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 2 {
		fatal(fmt.Sprintf("Wrong number of Args %#v: want exactly two json-file-names\n", flag.Args()))
	}
	var fromFile, toFile string
	fromFile = flag.Arg(0)
	toFile = flag.Arg(1)
	info("==== Compare two json files: " + fromFile + " , " + toFile)

	// read file
	var from, to []byte
	from = readFile(fromFile)
	to = readFile(toFile)

	// diff
	start := time.Now()
	differ := &diff.JsonDiffer{}
	delta, err := differ.DiffJsonStr(string(from), string(to))
	if err != nil {
		fatal(fmt.Errorf("DiffJsonStr() failed: %w", err).Error())
	}

	duration := time.Since(start)
	if *fDebug {
		info(fmt.Sprintf("----- diff duration: %v", duration))
		bs, _ := json.MarshalIndent(delta, "", "    ")
		info("===== debug: write delta into delta.json")
		os.WriteFile("delta.json", bs, 0644)
	}

	// render
	var view render.DeltaView
	switch *fFormat {
	case "html":
		view = render.NewHtmlView()
	case "console":
		view = render.NewConsoleView()
	case "patch":
		view = render.NewJsonPatchView()
	default:
		fatal(fmt.Sprintf("invalid input: unknown formatterType: %s", *fFormat))
	}

	start = time.Now()
	output, err := view.Render(delta)
	if err != nil {
		fatal(fmt.Errorf("render() failed: %w", err).Error())
	}

	duration = time.Since(start)
	if *fDebug {
		info(fmt.Sprintf("----- render duration: %v", duration))
	}

	fmt.Println(output)
}

func readFile(fname string) []byte {
	bytes, err := os.ReadFile(fname)
	if err != nil {
		fatal(fmt.Sprintf("Failed to read %s: %+v\n", fname, err))
	}
	return bytes
}

func fatal(msg string) {
	println(aurora.Red(msg).String())
	flag.Usage()
	os.Exit(1)
}

func info(msg string) {
	println(aurora.Cyan(msg).String())
}
