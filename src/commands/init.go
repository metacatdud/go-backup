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
	configFilePath := usr.HomeDir + "/" + "qbkp_config.json"
	backupPath := "/tmp/qbkp/"

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.Mkdir(configDir, 0777);
	}

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		os.Mkdir(backupPath, 0777);
	}

	if _, err := os.Stat(configFilePath); err == nil {
		return "file_exist"
	}

	configJson := gabs.New()
	configJson.Set(secret, "enckey")
	configJson.Set(remote, "remote")
	configJson.Set(usr.HomeDir + "/", "rootPath")
	configJson.Set(backupPath, "backupPath")


	err := ioutil.WriteFile(configFilePath, []byte(configJson.String()), 0644)
	if nil != err {
		panic(err)
	}

	return "file_created"
}