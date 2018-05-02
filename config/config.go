package config

import (
	"fmt"
	"os/user"

	"github.com/BurntSushi/toml"
)

// Config configuration struct
type Config struct {
	DownloadPath string
}

// Conf configuration struct
var Conf Config

// Parse parse config file
func Parse() {
	usr, _ := user.Current()
	if _, err := toml.DecodeFile(usr.HomeDir+"/.config/my-youtube/config.toml", &Conf); err != nil {
		fmt.Println(err)
	}
}
