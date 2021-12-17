package main

import (
	"os"

	"github.com/mholt/archiver/v3"
)

func AttemptPassword(archive_file string, target_file string, password string) bool {
	rar := archiver.Rar{
		MkdirAll:               true,
		ContinueOnError:        false,
		OverwriteExisting:      false,
		ImplicitTopLevelFolder: false,
		Password:               password,
	}
	err := rar.Extract(archive_file, target_file, os.TempDir())
	return err == nil
}

func GetSmallestFile(archive_file string) string {

	rar := archiver.Rar{
		MkdirAll:               true,
		ContinueOnError:        false,
		OverwriteExisting:      false,
		ImplicitTopLevelFolder: false,
		Password:               "x",
	}

	var lowest int64 = -1
	var smallestFile archiver.File

	rar.Walk(archive_file, func(f archiver.File) error {
		if lowest == -1 || f.Size() < lowest {
			lowest = f.Size()
			smallestFile = f
		}
		return nil
	})
	return smallestFile.Name()
}
