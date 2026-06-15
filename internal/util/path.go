package util

import (
	"os"
	"path/filepath"
	//"runtime"
)

func GetExeDirPath() string {
	exe, _ := os.Executable()
	exe, _ = filepath.EvalSymlinks(exe)
	dir := filepath.Dir(exe)
	return dir
}

func RelativeToAbsolute(rpath ...string) string {
	dir := GetExeDirPath()
	apath := filepath.Join(append([]string{dir}, rpath...)...)
	return apath
}

/*
func GetConfigDir() string {
	app := "bitcoin-cli"
	home, _ := os.UserHomeDir()

	switch runtime.GOOS {
	case "windows":
		// AppData\Roaming
		return filepath.Join(os.Getenv("APPDATA"), app)

	case "darwin":
		// macOS
		return filepath.Join(home, "Library", "Application Support", app)

	default:
		// Linux
		return filepath.Join(home, ".config", app)
	}
}
*/
