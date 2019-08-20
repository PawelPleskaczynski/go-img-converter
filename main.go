package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getDigits(input string, regex *regexp.Regexp) (int, error) {
	match := regex.FindStringSubmatch(input)
	height, err := strconv.Atoi(match[1])
	return height, err
}

func getDimensions(label *os.File) (int, int, int, string) {
	var width, height, bitdepth int
	var missionName string
	var err error

	// endRegexWhitespace := regexp.MustCompile(`\s{50,}`)
	linesRegex := regexp.MustCompile(`(?:LINES)(?:\s+)?(?:=)(?:\s+)?(\d+)`)
	lineSamplesRegex := regexp.MustCompile(`(?:LINE_SAMPLES)(?:\s+)?(?:=)(?:\s+)?(\d+)`)
	sampleBitsRegex := regexp.MustCompile(`(?:BITS)(?:\s+)?(?:=)(?:\s+)?(\d+)`)
	missionNameRegex := regexp.MustCompile(`(?:MISSION_NAME)(?:\s+)?(?:=)(?:\s+)(?:")?(.+)(?:")`)

	scanner := bufio.NewScanner(label)
	buf := make([]byte, 0, 64*1024*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		if linesRegex.MatchString(scanner.Text()) {
			heightStr := scanner.Text()
			height, err = getDigits(heightStr, linesRegex)
			if err != nil {
				log.Fatalln("No LINES string in LBL file!")
				panic(err)
			}
		}

		if lineSamplesRegex.MatchString(scanner.Text()) {
			widthStr := scanner.Text()
			width, err = getDigits(widthStr, lineSamplesRegex)
			if err != nil {
				log.Fatalln("No LINE_SAMPLES string in LBL file!")
				panic(err)
			}
		}

		if sampleBitsRegex.MatchString(scanner.Text()) {
			bitdepthStr := scanner.Text()
			bitdepth, err = getDigits(bitdepthStr, sampleBitsRegex)
			if err != nil {
				log.Fatalln("No SAMPLE_BITS string in LBL file!")
				panic(err)
			}
		}

		if missionNameRegex.MatchString(scanner.Text()) {
			missionNameStr := scanner.Text()
			missionName = missionNameRegex.FindStringSubmatch(missionNameStr)[1]
		}

		// if endRegexWhitespace.MatchString(scanner.Text()) {
		// 	whitespaceLength := len(endRegexWhitespace.FindStringSubmatch(scanner.Text())[0])
		// 	endbyte += whitespaceLength
		// 	println(endRegexWhitespace.FindStringSubmatch(scanner.Text())[0])
		// 	println(endbyte)
		// 	broke = true
		// 	break
		// } else {
		// 	endbyte += len(scanner.Text())
		// }
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// if broke == false {
	// 	endbyte = 0
	// }

	return width, height, bitdepth, missionName
}

func main() {
	filename := flag.String("input", "", "input file (.IMG)")
	labelname := flag.String("label", "", "label input file (.LBL)")
	flag.Parse()

	length := len(*filename)

	var width, height, bitdepth int
	var missionName string

	buf := bytes.NewBuffer(nil)
	file, err := os.Open(*filename)
	check(err)
	io.Copy(buf, file)
	fileslice := []byte(buf.Bytes())
	file.Close()

	if *labelname != "" {
		label, err := os.Open(*labelname)
		check(err)
		width, height, bitdepth, missionName = getDimensions(label)
	}

	if height == 0 || width == 0 || bitdepth == 0 {
		label, err := os.Open(*filename)
		check(err)
		width, height, bitdepth, missionName = getDimensions(label)
		if height == 0 || width == 0 || bitdepth == 0 {
			log.Fatalln("Can't get all needed image attributes, exiting")
		}
	}

	if missionName != "" {
		fmt.Printf("Mission: %s\n", missionName)
	}
	fmt.Printf("Width: %d, height: %d, bit depth: %d\n", width, height, bitdepth)

	if strings.EqualFold(missionName, "CASSINI-HUYGENS") {
		fileslice = fileslice[4192:]
		j, count512, count12 := 0, 0, 0
		array512 := make([]byte, len(fileslice))
		for i := 0; i < len(fileslice); i++ {
			if count512 < 1024 {
				array512[j] = fileslice[i]
				j++
				count512++
			} else if count12 < 23 && count512 == 1024 {
				count12++
			} else if count12 == 23 {
				count512 = 0
				count12 = 0
			}
		}
		fileslice = array512
	}

	if bitdepth == 8 {
		image := image.NewGray(image.Rect(0, 0, width, height))
		image.Pix = fileslice
		out, err := os.Create((*filename)[:length-4] + ".png")
		check(err)
		defer out.Close()
		png.Encode(out, image)
	} else if bitdepth == 12 || bitdepth == 16 || bitdepth == 32 {
		image := image.NewGray16(image.Rect(0, 0, width, height))
		image.Pix = fileslice
		out, err := os.Create((*filename)[:length-4] + ".png")
		check(err)
		defer out.Close()
		png.Encode(out, image)
	}
}
