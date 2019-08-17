package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	filename := flag.String("input", "", "input file (.IMG)")
	labelname := flag.String("label", "", "label input file (.LBL)")
	flag.Parse()

	length := len(*filename)

	var width, height, bitdepth int
	buf := bytes.NewBuffer(nil)
	file, err := os.Open(*filename)
	check(err)
	io.Copy(buf, file)
	fileslice := []byte(buf.Bytes())
	file.Close()

	label, err := os.Open(*labelname)
	check(err)
	scanner := bufio.NewScanner(label)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "LINES") {
			height, err = strconv.Atoi(strings.TrimSpace(scanner.Text()[32:]))
			if err != nil {
				fmt.Fprintln(os.Stderr, "No LINES string in LBL file!")
				panic(err)
			}
		}
		if strings.Contains(scanner.Text(), "LINE_SAMPLES") {
			width, err = strconv.Atoi(strings.TrimSpace(scanner.Text()[32:]))
			if err != nil {
				fmt.Fprintln(os.Stderr, "No LINE_SAMPLES string in LBL file!")
				panic(err)
			}
		}
		if strings.Contains(scanner.Text(), "SAMPLE_BITS") {
			bitdepth, err = strconv.Atoi(strings.TrimSpace(scanner.Text()[32:]))
			if err != nil {
				fmt.Fprintln(os.Stderr, "No SAMPLE_BITS string in LBL file!")
				panic(err)
			} else if bitdepth != 8 && bitdepth != 16 {
				fmt.Fprintln(os.Stderr, "Invalid bit depth")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	if bitdepth == 8 {
		image := image.NewGray(image.Rect(0, 0, width, height))
		image.Pix = fileslice
		out, err := os.Create((*filename)[:length-4] + ".png")
		check(err)
		defer out.Close()
		png.Encode(out, image)
	} else if bitdepth == 16 {
		image := image.NewGray16(image.Rect(0, 0, width, height))
		image.Pix = fileslice
		out, err := os.Create((*filename)[:length-4] + ".png")
		check(err)
		defer out.Close()
		png.Encode(out, image)
	}
}
