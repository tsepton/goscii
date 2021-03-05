package main

import (
	"fmt"
	"image"
	"os"
	"strconv"

	"github.com/nfnt/resize"

	_ "image/png"
)

func main() {
	image, err := readFile("./gopher.png")
	if err != nil {
		panic(err.Error())
	}

	divider := parseArgs()
	image = resize.Resize(
		uint(image.Bounds().Dx()/divider),
		uint(image.Bounds().Dy()/(divider*2)),
		image,
		resize.NearestNeighbor)

	writeFile(asciiArt(image), "output.txt")
}

func parseArgs() int {
	if len(os.Args) > 1 {
		if arg, err := strconv.ParseInt(os.Args[1], 10, 0); err != nil {
			fmt.Println("Warning: argument specified incorrect, it must be an integer")
		} else if arg < 1 {
			fmt.Println("Warning: argument specified incorrect, it must be superior to 0")
		} else {
			fmt.Println("Dividing proportion by ", arg)
			return int(arg)
		}
	}
	fmt.Println("Keeping original proportion")
	return 1
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

func asciiArt(image image.Image) string {
	greyscale := "@%#*+=-:. "

	min := image.Bounds().Min
	max := image.Bounds().Max
	var ascii string
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
	return ascii
}
