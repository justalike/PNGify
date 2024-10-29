package main

import (
	"flag"
	"fmt"
	"github.com/tofl/pngify/image"
	"os"
	"time"
)

func printHelp() {
	fmt.Println(`Usage:
  pngify [command] [options]

Commands:
  encode    Encode text or file data into a PNG image
  decode    Decode a PNG image to retrieve the original data
  help      Show this help message

Options:

Encode:
  -t "text"   Text to encode into the PNG image
  -f "file"   File to encode into the PNG image

Decode:
  -p "path"   Path to the PNG image to decode

Examples:
  Encode text:
    pngify encode -t "Your text here"
  
  Encode a file:
    pngify encode -f /path/to/file

  Decode a file:
    pngify decode -p /path/to/image.png

Output:
  An image named output.png will be created in the current directory for encoding.`)
}

func main() {
	encodeCmd := flag.NewFlagSet("encode", flag.ExitOnError)
	encodeText := encodeCmd.String("t", "", "Text to encode")
	encodeFile := encodeCmd.String("f", "", "File to encode")

	decodeCmd := flag.NewFlagSet("decode", flag.ExitOnError)
	filePath := decodeCmd.String("p", "", "Path to the file to decode")

	if len(os.Args) < 2 {
		fmt.Println("A command must be specified.")
		printHelp()
		os.Exit(1)
	}

	executionTime := time.Now()

	switch os.Args[1] {
	case "encode":
		_ = encodeCmd.Parse(os.Args[2:])
		var img *image.Image

		if *encodeText != "" {
			img = image.NewImage([]byte(*encodeText))
		} else if *encodeFile != "" {
			content, err := os.ReadFile(*encodeFile)
			if err != nil {
				fmt.Println("Couldn't read file", *encodeFile)
				os.Exit(1)
			}

			fmt.Println("Encoding file...")
			img = image.NewImage(content)
			img.MakeText([]byte("filename"), []byte(*encodeFile))
		}

		img.MakeImage()

	case "decode":
		_ = decodeCmd.Parse(os.Args[2:])
		f, err := os.Open(*filePath)
		defer f.Close()

		if err != nil {
			fmt.Println("Couldn't open the file")
			os.Exit(1)
		}

		fmt.Println("Decoding image...")
		data, fileName := image.Decode(f)

		if fileName == "" {
			fmt.Println(data)
		} else {
			f, err := os.Create(fileName)
			defer f.Close()

			if err != nil {
				fmt.Println("Couldn't create the output file.")
				os.Exit(1)
			}

			_, err = f.Write([]byte(data))
			if err != nil {
				fmt.Println("Couldn't write data to the output file.")
				os.Exit(1)
			}
		}

	case "help", "--help", "-help":
		printHelp()

	default:
		fmt.Println("Command not recognised")
		printHelp()
		os.Exit(1)
	}

	elapsed := time.Since(executionTime)
	fmt.Println("Done in", elapsed)
}
