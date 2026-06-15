package util

import (
	"os"
	"path/filepath"
	//"runtime"
)

func Get_exe_dir_path() string {
	exe, _ := os.Executable()
	exe, _ = filepath.EvalSymlinks(exe)
	dir := filepath.Dir(exe)
	return dir
}

func Relative_to_absolute(rpath ...string) string {
	dir := Get_exe_dir_path()
	apath := filepath.Join(append([]string{dir}, rpath...)...)
	return apath
}

/*
func Get_config_dir() string {
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
