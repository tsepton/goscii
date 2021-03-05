package main

import (
	"image"
	"os"

	"image/color"
	_ "image/png"

	"github.com/nfnt/resize"
)

func main() {
	image, err := readFile("./gopher.png")
	if err != nil {
		panic(err.Error())
	}
	image = resize.Resize(100, 25, image, resize.NearestNeighbor)

	writeFile(asciiArt(image), "output.txt")
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
	greyscale := " .:-=+*#%@"

	min := image.Bounds().Min
	max := image.Bounds().Max
	var ascii string
	for col := min.Y; col < max.Y; col++ {
		for row := min.X; row < max.X; row++ {
			grey := color.GrayModel.Convert(image.At(row, col)).(color.Gray).Y
			key := int((float64(grey) / 255) * float64(len(greyscale)-1))
			ascii = ascii + string(greyscale[int(key)])
		}
		ascii = ascii + "\n"
	}
	return ascii
}
