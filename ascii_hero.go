package main

import (
	"strings"
)

// Bubble-style ASCII (~46 cols per line) for a strong curl hero.
const asciiNameAyush = `    █████╗ ██╗   ██╗██╗   ██╗███████╗██╗  ██╗
   ██╔══██╗╚██╗ ██╔╝██║   ██║██╔════╝██║  ██║
   ███████║ ╚████╔╝ ██║   ██║███████╗███████║
   ██╔══██║  ╚██╔╝  ██║   ██║╚════██║██╔══██║
   ██║  ██║   ██║   ╚██████╔╝███████║██║  ██║
   ╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚══════╝╚═╝  ╚═╝`

   const asciiNameKashyap = `   ██╗  ██╗ █████╗ ███████╗██╗  ██╗██╗   ██╗ █████╗ ██████╗ 
   ██║ ██╔╝██╔══██╗██╔════╝██║  ██║██║   ██║██╔══██╗██╔══██╗
   █████╔╝ ███████║███████╗███████║██║   ██║███████║██████╔╝
   ██╔═██╗ ██╔══██║╚════██║██╔══██║██║   ██║██╔══██║██╔═══╝ 
   ██║  ██╗██║  ██║███████║██║  ██║╚██████╔╝██║  ██║██║     
   ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝╚═╝     `

// Slim “terminal” panel on the left (~26 cols wide). Kept for optional reuse.
const ASCIILeftPanel = `   +--------------+
   |##############|
   |#  .------.  #|
   |# (        ) #|
   |#  '------'  #|
   |#    |  |    #|
   |##############|
   +--------------+
   ___|________|___`

// Block “HELLO” in # (5 lines, monospace).
const ASCIIHelloLetters = `
#   # ##### #     #     # 
#   # #     #     #     # 
##### ##### #     #     # 
#   # #     #     #     # 
#   # ##### ##### ##### #####`

const ASCIIWorldLine = `        w  o  r  l  d`

// Small ASCII cat (under “world”).
const CatASCII = `
       /\_/\  
      ( o.o ) 
       > ^ <
      /|   |\
     (_|   |_)`

func splitLines(s string) []string {
	s = strings.Trim(s, "\n")
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

func runeLen(s string) int {
	return len([]rune(s))
}
