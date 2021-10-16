package utils

import (
	"embed"
	"io/fs"
)

func GetFrontendFiles(files embed.FS, baseDir string) fs.FS {
	fsys := fs.FS(files)
	frontendFiles, _ := fs.Sub(fsys, baseDir)
	return frontendFiles
}