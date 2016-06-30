package commands

import (
	"gopkg.in/urfave/cli.v2"
	"github.com/davecgh/go-spew/spew"
	"github.com/Jeffail/gabs"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
)

var (
	usr, _ = user.Current()
	configDir = usr.HomeDir + "/qbyco_bkp/"
	configFilePath = configDir + "config.json"
	configJson *gabs.Container
)

func init ()  {

}

//Command constructor
func BackupCommand() *cli.Command {

	command := &cli.Command {
		Name: "backup",
		Usage: "Manage backups",
		Aliases: []string{"bk"},
		Subcommands: []*cli.Command{
			{
				Name:  "add",
				Usage: "Add a new path to backup monitor",
				Action: func(c *cli.Context) error {
					spew.Println("New command not implemented yet")
					//c.Args().First()
					return nil
				},
			},
			{
				Name:  "remove",
				Aliases: []string{"rm"},
				Usage: "Remove an existing path from backup monitor",
				Action: func(c *cli.Context) error {
					spew.Println("Remove command not implemented yet")
					//c.Args().First()
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			response := runBackup()

			if "no_config" == response {
				return cli.Exit("Backup not init. Run:: 'qbkp init' first", 2001)
			}

			if "no_monitor" == response {
				return cli.Exit("No folders monitored yet. Run :: qbkp backup help", 2002)
			}

			if "archive_error" == response {
				return cli.Exit("Unable to create archives", 2003)
			}

			if "rsync_error" == response {
				return cli.Exit("Unable to sync", 2003)
			}

			return nil
		},
	}

	return command
}

func runBackup () string {

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return "no_config"
	}

	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}

	configJson, err = gabs.ParseJSON(configFile)

	if err != nil {
		panic(err)
	}

	tar := runTar()

	if "no_monitor" == tar || "archive_error" == tar {
		return "archive_error"
	}

	runEncrypt()

	rsync := runRsync()
	if "rsync_error" == rsync {
		return "rsync_error"
	}

	return "ok"
}

func runTar() string {

	backupPath := configJson.S("backupPath").Data().(string)
	rootPath := configJson.S("rootPath").Data().(string)
	foldersList, _ := configJson.S("monitor").Children()

	if nil == foldersList {
		return "no_monitor"
	}

	for _, folder := range foldersList {

		folderName := folder.Data().(string)
		archiveName := folderName + ".tar.gz"

		tarCmd := "tar"
		tarArgs := []string{
			"-zcvf",
			backupPath + archiveName,
			rootPath + folderName,
		}

		if err := exec.Command(tarCmd, tarArgs...).Run(); err != nil {
			return "archive_error"
		}



		spew.Printf("Archive:: %s Complete!\n", archiveName)
	}

	return "tar_complete"

}

func runEncrypt() {
	backupPath := configJson.S("backupPath").Data().(string)
	secret := configJson.S("enckey").Data().(string)

	//Execute encryption
	files, _ := ioutil.ReadDir(backupPath)
	for _, file := range files {

		mcryptCmd := "mcrypt"
		mcryptArgs := []string{
			backupPath + file.Name(),
			"--key",
			secret,
			"--unlink",
			"--force",
		}

		if err := exec.Command(mcryptCmd, mcryptArgs...).Run(); err != nil {
			spew.Print("Unable to encrypt", file.Name())
		}

		spew.Printf("Encrypt:: %s Success!\n", file.Name())
	}
}

func runRsync () string {
	backupPath := configJson.S("backupPath").Data().(string)
	remote := configJson.S("remote").Data().(string)

	rsyncCmd := "rsync"
	rsyncArgs := []string{
		"-av",
		"--delete",
		"-e ssh",
		backupPath,
		remote,

	}

	if err := exec.Command(rsyncCmd, rsyncArgs...).Run(); err != nil {
		return "rsync_error"
	}

	files, _ := ioutil.ReadDir(backupPath)
	for _, file := range files {
		spew.Printf("Sync File:: %s\n", file.Name())
		os.Remove(backupPath + file.Name())
	}
	spew.Print("Sync Complete!\n")

	return "ok"
}