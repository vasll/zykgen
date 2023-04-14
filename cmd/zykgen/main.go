package main

import (
	"fmt"
	"github.com/docopt/docopt.go"
	"github.com/luc10/zykgen"
	"os"
)

const usage = `Zyxel VMG8823-B50B WPA Keygen

Usage:
  zykgen --pass (-m|-n|-c) [-o <file>] [-l <length>] <serial>
  zykgen --dump (-m|-n|-c) [-o <file>] [-l <length>] <serialStart> <serialEnd>
  zykgen -h | --help

Options:
  -o <file>		  Output file
  -l <length>     Output key length [default: 10].
  -h --help       Show this screen.`

func main() {
	var cocktail zykgen.Cocktail
	var args struct {
		Serial       string `docopt:"<serial>"`
		SerialStart  string `docopt:"<serialStart>"`
		SerialEnd    string `docopt:"<serialEnd>"`
		Length       int    `docopt:"-l"`
		Mojito       bool   `docopt:"-m"`
		Negroni      bool   `docopt:"-n"`
		Cosmopolitan bool   `docopt:"-c"`
		File		 string	`docopt:"-o"`
		Pass		 bool 	`docopt:"--pass"`
		Dump	     bool   `docopt:"--dump"`
	}

	opts, err := docopt.DefaultParser.ParseArgs(usage, os.Args[1:], "")
	if err != nil {
		return
	}

	opts.Bind(&args)
	if args.Mojito {
		cocktail = zykgen.Mojito
	}
	if args.Negroni {
		cocktail = zykgen.Negroni
	}
	if args.Cosmopolitan {
		cocktail = zykgen.Cosmopolitan
	}

	if args.Pass {
		password := zykgen.Wpa(args.Serial, args.Length, cocktail)
		fmt.Printf("Serial: %s\nPassword: %s", args.Serial, password)
		if len(args.File) > 0 {  // If there is a file
			writeToFile(args.File, fmt.Sprintf("Serial: %s\nPassword: %s", args.Serial, password))
			fmt.Printf("\nWritten to file '%s'", args.File)
		}
	}else if args.Dump{
		fmt.Println(args.File)
	}
	
}

func writeToFile(filename string, content string){
	file, err := os.Create(filename)
    if err != nil {
        return
    }
    defer file.Close()
    file.WriteString(content)
}