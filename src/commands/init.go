package commands

import (
	"gopkg.in/urfave/cli.v2"
	"os"
	"os/user"
	"github.com/Jeffail/gabs"
	"io/ioutil"
)

//Command constructor
func InitCommand() *cli.Command {
	var enckey string
	var remote string

	command := &cli.Command {
		Name: "init",
		Usage: "Initialize a backup",
		Action: func (c *cli.Context) error {

			if 0 == len(enckey) {
				return cli.Exit("You must specify a secret key", 1001)
			}

			if 0 == len(remote) {
				return cli.Exit("You must specify a remote address (ssh user@domain.com:/path_to_backup_folder)", 1002)
			}

			response := createConfig(enckey, remote)
			return cli.Exit(response, 1003)

		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "secret",
				Usage: "Setup a secret key for encrypt and decrypt",
				Destination: &enckey,
			},
			&cli.StringFlag{
				Name: "remote",
				Usage: "Specify remote path for backup (only ssh supported)",
				Destination: &remote,
			},
		},
	}

	return command
}

func createConfig (secret string, remote string) string {

	usr, _ := user.Current()
	configDir := usr.HomeDir + "/qbyco_bkp/"
	configFile := configDir + "config.json"

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.Mkdir(configDir, 0777);
	}

	if _, err := os.Stat(configDir + "bkp"); os.IsNotExist(err) {
		os.Mkdir(configDir + "bkp", 0777);
	}

	if _, err := os.Stat(configFile); err == nil {
		return "file_exist"
	}

	configJson := gabs.New()
	configJson.Set(secret, "enckey")
	configJson.Set(remote, "remote")
	configJson.Set(usr.HomeDir + "/", "rootPath")
	configJson.Set(configDir + "bkp/", "backupPath")


	err := ioutil.WriteFile(configFile, []byte(configJson.String()), 0644)
	if nil != err {
		panic(err)
	}

	return "file_created"
}