package commands

import (
	"gopkg.in/urfave/cli.v2"
	"github.com/davecgh/go-spew/spew"
	"github.com/Jeffail/gabs"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"time"
)

var (
	usr, _ = user.Current()
	configDir = usr.HomeDir + "/"
	configFilePath = configDir + "qbkp_config.json"
	configJson *gabs.Container
)


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

	defer timeTrack(time.Now(), "[PROCESS][BACKUP]:: Total time")

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

	//runEncrypt()

	rsync := runRsync()
	if "rsync_error" == rsync {
		return "rsync_error"
	}

	return "ok"
}

func runTar() string {
	defer timeTrack(time.Now(), "[PROCESS][ZIP]:: Time")

	backupPath := configJson.S("backupPath").Data().(string)
	rootPath := configJson.S("rootPath").Data().(string)
	foldersList, _ := configJson.S("monitor").Children()
	secret := configJson.S("enckey").Data().(string)

	if nil == foldersList {
		return "no_monitor"
	}

	for _, folder := range foldersList {

		backupRoot := false
		folderName := folder.Data().(string)
		archiveName := folderName + ".zip"

		if "*" == folderName {
			backupRoot = true
			archiveName = usr.Username + ".zip"
		}

		zipCmd := "zip"
		zipArgs := []string{
			"-9",
			"-P",
			secret,
			"-s",
			"1024m",
			"-r",
			spew.Sprintf("%s", backupPath + archiveName),
		}

		if true == backupRoot {
			zipArgs = append(zipArgs, rootPath)

		} else {
			zipArgs = append(zipArgs, spew.Sprintf("%s/", rootPath + folderName))

		}

		runCmd := exec.Command(zipCmd, zipArgs...)
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		runCmd.Run()
	}

	return "tar_complete"

}

func runRsync () string {
	defer timeTrack(time.Now(), "[PROCESS][RSYNC]:: Time")

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

	runCmd := exec.Command(rsyncCmd, rsyncArgs...)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	runCmd.Run()

	files, _ := ioutil.ReadDir(backupPath)
	for _, file := range files {
		os.Remove(backupPath + file.Name())
	}

	return "ok"
}

// Helper
// Time tracking
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	spew.Printf("%s took %s\n\n", name, elapsed)
}