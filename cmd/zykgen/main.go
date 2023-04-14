package main

import (
	"fmt"
	"github.com/docopt/docopt.go"
	"zykgen"
	"os"
	"strconv"
	"time"
	"github.com/schollz/progressbar/v3"
)

const usage = `Zyxel VMG8823-B50B WPA Keygen

Usage:
  zykgen --pass (-m|-n|-c) [-o <file>] [-l <length>] <serial>
  zykgen --dump (-m|-n|-c) -o <file> [-l <length>] <RouterRange>
  zykgen -h | --help

Options:
  <RouterRange>	  "homelife" or "infostrada" or "tiscali"
  -o <file>       Output file
  -l <length>     Output key length [default: 10].
  -h --help       Show this screen.`

func main() {
	var args struct {
		Serial       string `docopt:"<serial>"`
		RouterRange  string `docopt:"<RouterRange>"`
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

	// Get cocktail (letterlist)
	var cocktail zykgen.Cocktail
	if args.Mojito {
		cocktail = zykgen.Mojito
	} else if args.Negroni {
		cocktail = zykgen.Negroni
	} else if args.Cosmopolitan {
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
		// Get routerRange
		var routerRange RouterRange
		if args.RouterRange == "homelife" {
			routerRange = RangeHomeLife
		} else if args.RouterRange == "infostrada" {
			routerRange = RangeInfostrada
		} else if args.RouterRange == "tiscali" {
			routerRange = RangeTiscali
		} else {
			fmt.Println("RouterRange is not valid! Pick something between 'homelife', 'infostrada', 'tiscali'")
			os.Exit(-1)
		}
		fmt.Println("Writing passwords to file...")
		passwordRangeToFile(args.Length, cocktail, routerRange, args.File)
	}
}

// Writes a string to a file
func writeToFile(filename string, content string){
	file, err := os.Create(filename)
    if err != nil {
        return
    }
    defer file.Close()
    file.WriteString(content)
}

// Defines a Router range, for example from R182V to S192V
type RouterRange struct {
    first, second string
}
var (
    RangeHomeLife = RouterRange{"S182V", "S192V"}
    RangeInfostrada = RouterRange{"S172V", "S182V"}
	RangeTiscali = RouterRange{"S172V", "S182V"}
)

// Writes a password range to a file
func passwordRangeToFile(length int, cocktail zykgen.Cocktail, routerRange RouterRange, filename string){
	rangeNumEnd := 99999999
	pbar := progressbar.Default(int64(rangeNumEnd*2))

	// Write output to file
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()
	
	// Start progressbar updater thread
	i:=0
	x:=0
	go func() {
		for {
			pbar.Set(i+x)
			time.Sleep(time.Second)
		}
	}()

	// Iterate through first range i.e for homelife "S182V"
	for i = 0; i <= rangeNumEnd; i++ {
        serial := fmt.Sprintf("%s%08s", routerRange.first, strconv.Itoa(i))	// Get generated serial
		password := zykgen.Wpa(serial, length, cocktail)	// Get password from serial
		file.WriteString(fmt.Sprintf("%s\n", password))  // Write to file
    }

	// Iterate through second range i.e for homelife "S192V"
	for x = 0; x <= rangeNumEnd; x++ {
		serial := fmt.Sprintf("%s%08s", routerRange.second, strconv.Itoa(x))	// Get generated serial
		password := zykgen.Wpa(serial, length, cocktail)	// Get password from serial
		file.WriteString(fmt.Sprintf("%s\n", password))  // Write to file
	}

	file.Close()
	return
}