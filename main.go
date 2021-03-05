package main

import (
	"errors"
	"fmt"
	"image"
	"os"
	"strconv"

	"github.com/nfnt/resize"

	_ "image/jpeg"
	_ "image/png"
)

func main() {
	input, output, divider, err := parseArgs()
	if err != nil {
		fmt.Println(err.Error())
		showHelp()
		os.Exit(2)
	}

	image, err := readFile(input)
	if err != nil {
		panic(err.Error())
	}
	image = resize.Resize(
		uint(image.Bounds().Dx()/divider),
		uint(image.Bounds().Dy()/(int(float64(divider)*2.5))),
		image,
		resize.NearestNeighbor)

	writeFile(asciiArt(image), output)
}

func parseArgs() (input string, output string, divider int, err error) {
	divider = 1 // default value is 0, change it to 1...
	if len(os.Args) > 2 {
		if len(os.Args[1:])%2 != 0 {
			err = errors.New("Missing value for a specified flag.")
			return
		}
		for index, value := 1, 2; value < len(os.Args); index, value = index+2, value+2 {
			switch os.Args[index] {
			case "-i":
				input = os.Args[value]
			case "-o":
				output = os.Args[value]
			case "-d":
				divider, err = parseDivider(os.Args[value])
			default:
				err = errors.New("Unknown argument specified.")
				return
			}
		}
	} else {
		err = errors.New("No arguments were specified.")
		return
	}
	if output == "" || input == "" {
		err = errors.New("Missing input and/or output files.")
		return
	}
	return
}

func showHelp() {
	fmt.Println("[-f file] [-o output] [-d int]")
	fmt.Println("	-i input  | jpg or png input ")
	fmt.Println("	-o output | output text file ")
	fmt.Println("	-d int    | (optional) int by which original file ratio will be divided ")
}

func parseDivider(args string) (int, error) {
	if divider, err := strconv.ParseInt(args, 10, 0); err != nil {
		return 0,
			errors.New("Warning: argument specified incorrect, it must be an integer")
	} else if divider < 1 {
		return 0,
			errors.New("Warning: argument specified incorrect, it must be superior to 0")
	} else {
		return int(divider), nil
	}
}

func readFile(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	image, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return image, nil
}

func writeFile(ascii string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(ascii)
	return err
}

func asciiArt(image image.Image) (ascii string) {
	greyscale := "@%#*+=-:. "

	min := image.Bounds().Min
	max := image.Bounds().Max
	for col := min.Y; col < max.Y; col++ {
		for row := min.X; row < max.X; row++ {
			r, g, b, alpha := image.At(row, col).RGBA()
			if alpha == 0 {
				ascii = ascii + string(' ')
				continue
			}
			grey := (19595*r + 38470*g + 7471*b + 1<<15) >> 24
			key := int(float64(grey) / 255.0 * float64(len(greyscale)-1))
			ascii = ascii + string(greyscale[key])
		}
		ascii = ascii + "\n"
	}
	return
}
