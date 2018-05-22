package main

import (
	"github.com/deathmaz/my-youtube/config"
	"github.com/deathmaz/my-youtube/ui"
)

func main() {
	config.Parse()
	ui.Run()
}
