package openalex

import (
	"log"
	"os"
	"os/exec"
)

// download snapshot from openalex, "AWS CLI" installation required
// https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html
func Sync(destPath string) (err error) {
	source := "s3://openalex"
	arg := "--no-sign-request"
	argDelete := "--delete"
	dest := destPath

	// aws sync copies new or modified files to the destination, but does not delete old files
	downloadCmd := exec.Command("aws", "s3", "sync", source, dest, arg)
	downloadCmd.Stdout = os.Stdout
	if err := downloadCmd.Run(); err != nil {
		log.Fatal(err)
		return err
	}

	// delete outdated data that exist in the destination but not in the source
	deleteCmd := exec.Command("aws", "s3", "sync", source, dest, arg, argDelete)
	deleteCmd.Stdout = os.Stdout
	if err := deleteCmd.Run(); err != nil {
		log.Fatal(err)
		return err
	}

	return err
}
