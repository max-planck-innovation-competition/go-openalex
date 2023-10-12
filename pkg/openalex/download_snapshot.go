package openalex

import (
	"log/slog"
	"os"
	"os/exec"
)

// Sync downloads the latest snapshot from openalex
// "AWS CLI" installation required
// https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html
func Sync(destPath string) (err error) {
	logger := slog.With("destPath", destPath)
	source := "s3://openalex"
	arg := "--no-sign-request"
	argDelete := "--delete"
	dest := destPath

	// aws sync copies new or modified files to the destination, but does not delete old files
	downloadCmd := exec.Command("aws", "s3", "sync", source, dest, arg)
	downloadCmd.Stdout = os.Stdout
	err = downloadCmd.Run()
	if err != nil {
		logger.With("err", err).Error("error while downloading snapshot")
		return err
	}

	// delete outdated data that exist in the destination but not in the source
	deleteCmd := exec.Command("aws", "s3", "sync", source, dest, arg, argDelete)
	deleteCmd.Stdout = os.Stdout
	err = deleteCmd.Run()
	if err != nil {
		logger.With("err", err).Error("error while deleting outdated data")
		return err
	}

	return err
}
