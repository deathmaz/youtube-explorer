package ui

import (
	"log"
	"math"
	"os/exec"
	"strings"
	"unicode"
)

func runcmd(cmd string, shell bool) []byte {
	if shell {
		err := exec.Command("bash", "-c", cmd).Start()
		if err != nil {
			log.Fatal(err)
			panic("some error found")
		}
	}
	out, err := exec.Command(cmd).Output()
	if err != nil {
		log.Fatal(err)
	}
	return out
}

// Round round
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

// SpaceMap remove all whitespaces from a string
func SpaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
