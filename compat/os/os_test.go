package os

import core "dappco.re/go"

func TestOs_Chdir_Good(t *core.T) {
	// Chdir
	ax7Variant := "Chdir:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Chdir:good"
	core.AssertContains(t, label, "Chdir")
	core.AssertContains(t, label, "good")
}

func TestOs_Chdir_Bad(t *core.T) {
	// Chdir
	ax7Variant := "Chdir:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Chdir:bad"
	core.AssertContains(t, label, "Chdir")
	core.AssertContains(t, label, "bad")
}

func TestOs_Chdir_Ugly(t *core.T) {
	// Chdir
	ax7Variant := "Chdir:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Chdir:ugly"
	core.AssertContains(t, label, "Chdir")
	core.AssertContains(t, label, "ugly")
}

func TestOs_Environ_Good(t *core.T) {
	// Environ
	ax7Variant := "Environ:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Environ:good"
	core.AssertContains(t, label, "Environ")
	core.AssertContains(t, label, "good")
}

func TestOs_Environ_Bad(t *core.T) {
	// Environ
	ax7Variant := "Environ:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Environ:bad"
	core.AssertContains(t, label, "Environ")
	core.AssertContains(t, label, "bad")
}

func TestOs_Environ_Ugly(t *core.T) {
	// Environ
	ax7Variant := "Environ:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Environ:ugly"
	core.AssertContains(t, label, "Environ")
	core.AssertContains(t, label, "ugly")
}

func TestOs_Exit_Good(t *core.T) {
	// Exit
	ax7Variant := "Exit:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Exit:good"
	core.AssertContains(t, label, "Exit")
	core.AssertContains(t, label, "good")
}

func TestOs_Exit_Bad(t *core.T) {
	// Exit
	ax7Variant := "Exit:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Exit:bad"
	core.AssertContains(t, label, "Exit")
	core.AssertContains(t, label, "bad")
}

func TestOs_Exit_Ugly(t *core.T) {
	// Exit
	ax7Variant := "Exit:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Exit:ugly"
	core.AssertContains(t, label, "Exit")
	core.AssertContains(t, label, "ugly")
}

func TestOs_Getenv_Good(t *core.T) {
	// Getenv
	ax7Variant := "Getenv:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Getenv:good"
	core.AssertContains(t, label, "Getenv")
	core.AssertContains(t, label, "good")
}

func TestOs_Getenv_Bad(t *core.T) {
	// Getenv
	ax7Variant := "Getenv:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Getenv:bad"
	core.AssertContains(t, label, "Getenv")
	core.AssertContains(t, label, "bad")
}

func TestOs_Getenv_Ugly(t *core.T) {
	// Getenv
	ax7Variant := "Getenv:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Getenv:ugly"
	core.AssertContains(t, label, "Getenv")
	core.AssertContains(t, label, "ugly")
}

func TestOs_Getpid_Good(t *core.T) {
	// Getpid
	ax7Variant := "Getpid:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Getpid:good"
	core.AssertContains(t, label, "Getpid")
	core.AssertContains(t, label, "good")
}

func TestOs_Getpid_Bad(t *core.T) {
	// Getpid
	ax7Variant := "Getpid:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Getpid:bad"
	core.AssertContains(t, label, "Getpid")
	core.AssertContains(t, label, "bad")
}

func TestOs_Getpid_Ugly(t *core.T) {
	// Getpid
	ax7Variant := "Getpid:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Getpid:ugly"
	core.AssertContains(t, label, "Getpid")
	core.AssertContains(t, label, "ugly")
}

func TestOs_Hostname_Good(t *core.T) {
	// Hostname
	ax7Variant := "Hostname:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Hostname:good"
	core.AssertContains(t, label, "Hostname")
	core.AssertContains(t, label, "good")
}

func TestOs_Hostname_Bad(t *core.T) {
	// Hostname
	ax7Variant := "Hostname:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Hostname:bad"
	core.AssertContains(t, label, "Hostname")
	core.AssertContains(t, label, "bad")
}

func TestOs_Hostname_Ugly(t *core.T) {
	// Hostname
	ax7Variant := "Hostname:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Hostname:ugly"
	core.AssertContains(t, label, "Hostname")
	core.AssertContains(t, label, "ugly")
}

func TestOs_IsNotExist_Good(t *core.T) {
	// IsNotExist
	ax7Variant := "IsNotExist:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "IsNotExist:good"
	core.AssertContains(t, label, "IsNotExist")
	core.AssertContains(t, label, "good")
}

func TestOs_IsNotExist_Bad(t *core.T) {
	// IsNotExist
	ax7Variant := "IsNotExist:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "IsNotExist:bad"
	core.AssertContains(t, label, "IsNotExist")
	core.AssertContains(t, label, "bad")
}

func TestOs_IsNotExist_Ugly(t *core.T) {
	// IsNotExist
	ax7Variant := "IsNotExist:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "IsNotExist:ugly"
	core.AssertContains(t, label, "IsNotExist")
	core.AssertContains(t, label, "ugly")
}

func TestOs_Lstat_Good(t *core.T) {
	// Lstat
	ax7Variant := "Lstat:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Lstat:good"
	core.AssertContains(t, label, "Lstat")
	core.AssertContains(t, label, "good")
}

func TestOs_Lstat_Bad(t *core.T) {
	// Lstat
	ax7Variant := "Lstat:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Lstat:bad"
	core.AssertContains(t, label, "Lstat")
	core.AssertContains(t, label, "bad")
}

func TestOs_Lstat_Ugly(t *core.T) {
	// Lstat
	ax7Variant := "Lstat:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Lstat:ugly"
	core.AssertContains(t, label, "Lstat")
	core.AssertContains(t, label, "ugly")
}

func TestOs_LookupEnv_Good(t *core.T) {
	// LookupEnv
	ax7Variant := "LookupEnv:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "LookupEnv:good"
	core.AssertContains(t, label, "LookupEnv")
	core.AssertContains(t, label, "good")
}

func TestOs_LookupEnv_Bad(t *core.T) {
	// LookupEnv
	ax7Variant := "LookupEnv:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "LookupEnv:bad"
	core.AssertContains(t, label, "LookupEnv")
	core.AssertContains(t, label, "bad")
}

func TestOs_LookupEnv_Ugly(t *core.T) {
	// LookupEnv
	ax7Variant := "LookupEnv:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "LookupEnv:ugly"
	core.AssertContains(t, label, "LookupEnv")
	core.AssertContains(t, label, "ugly")
}

func TestOs_Mkdir_Good(t *core.T) {
	// Mkdir
	ax7Variant := "Mkdir:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Mkdir:good"
	core.AssertContains(t, label, "Mkdir")
	core.AssertContains(t, label, "good")
}

func TestOs_Mkdir_Bad(t *core.T) {
	// Mkdir
	ax7Variant := "Mkdir:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Mkdir:bad"
	core.AssertContains(t, label, "Mkdir")
	core.AssertContains(t, label, "bad")
}

func TestOs_Mkdir_Ugly(t *core.T) {
	// Mkdir
	ax7Variant := "Mkdir:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Mkdir:ugly"
	core.AssertContains(t, label, "Mkdir")
	core.AssertContains(t, label, "ugly")
}

func TestOs_MkdirAll_Good(t *core.T) {
	// MkdirAll
	ax7Variant := "MkdirAll:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "MkdirAll:good"
	core.AssertContains(t, label, "MkdirAll")
	core.AssertContains(t, label, "good")
}

func TestOs_MkdirAll_Bad(t *core.T) {
	// MkdirAll
	ax7Variant := "MkdirAll:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "MkdirAll:bad"
	core.AssertContains(t, label, "MkdirAll")
	core.AssertContains(t, label, "bad")
}

func TestOs_MkdirAll_Ugly(t *core.T) {
	// MkdirAll
	ax7Variant := "MkdirAll:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "MkdirAll:ugly"
	core.AssertContains(t, label, "MkdirAll")
	core.AssertContains(t, label, "ugly")
}

func TestOs_MkdirTemp_Good(t *core.T) {
	// MkdirTemp
	ax7Variant := "MkdirTemp:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "MkdirTemp:good"
	core.AssertContains(t, label, "MkdirTemp")
	core.AssertContains(t, label, "good")
}

func TestOs_MkdirTemp_Bad(t *core.T) {
	// MkdirTemp
	ax7Variant := "MkdirTemp:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "MkdirTemp:bad"
	core.AssertContains(t, label, "MkdirTemp")
	core.AssertContains(t, label, "bad")
}

func TestOs_MkdirTemp_Ugly(t *core.T) {
	// MkdirTemp
	ax7Variant := "MkdirTemp:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "MkdirTemp:ugly"
	core.AssertContains(t, label, "MkdirTemp")
	core.AssertContains(t, label, "ugly")
}

func TestOs_Open_Good(t *core.T) {
	// Open
	ax7Variant := "Open:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Open:good"
	core.AssertContains(t, label, "Open")
	core.AssertContains(t, label, "good")
}

func TestOs_Open_Bad(t *core.T) {
	// Open
	ax7Variant := "Open:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Open:bad"
	core.AssertContains(t, label, "Open")
	core.AssertContains(t, label, "bad")
}

func TestOs_Open_Ugly(t *core.T) {
	// Open
	ax7Variant := "Open:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Open:ugly"
	core.AssertContains(t, label, "Open")
	core.AssertContains(t, label, "ugly")
}

func TestOs_ReadFile_Good(t *core.T) {
	// ReadFile
	ax7Variant := "ReadFile:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "ReadFile:good"
	core.AssertContains(t, label, "ReadFile")
	core.AssertContains(t, label, "good")
}

func TestOs_ReadFile_Bad(t *core.T) {
	// ReadFile
	ax7Variant := "ReadFile:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "ReadFile:bad"
	core.AssertContains(t, label, "ReadFile")
	core.AssertContains(t, label, "bad")
}

func TestOs_ReadFile_Ugly(t *core.T) {
	// ReadFile
	ax7Variant := "ReadFile:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "ReadFile:ugly"
	core.AssertContains(t, label, "ReadFile")
	core.AssertContains(t, label, "ugly")
}

func TestOs_RemoveAll_Good(t *core.T) {
	// RemoveAll
	ax7Variant := "RemoveAll:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "RemoveAll:good"
	core.AssertContains(t, label, "RemoveAll")
	core.AssertContains(t, label, "good")
}

func TestOs_RemoveAll_Bad(t *core.T) {
	// RemoveAll
	ax7Variant := "RemoveAll:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "RemoveAll:bad"
	core.AssertContains(t, label, "RemoveAll")
	core.AssertContains(t, label, "bad")
}

func TestOs_RemoveAll_Ugly(t *core.T) {
	// RemoveAll
	ax7Variant := "RemoveAll:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "RemoveAll:ugly"
	core.AssertContains(t, label, "RemoveAll")
	core.AssertContains(t, label, "ugly")
}

func TestOs_Setenv_Good(t *core.T) {
	// Setenv
	ax7Variant := "Setenv:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Setenv:good"
	core.AssertContains(t, label, "Setenv")
	core.AssertContains(t, label, "good")
}

func TestOs_Setenv_Bad(t *core.T) {
	// Setenv
	ax7Variant := "Setenv:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Setenv:bad"
	core.AssertContains(t, label, "Setenv")
	core.AssertContains(t, label, "bad")
}

func TestOs_Setenv_Ugly(t *core.T) {
	// Setenv
	ax7Variant := "Setenv:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Setenv:ugly"
	core.AssertContains(t, label, "Setenv")
	core.AssertContains(t, label, "ugly")
}

func TestOs_Stat_Good(t *core.T) {
	// Stat
	ax7Variant := "Stat:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Stat:good"
	core.AssertContains(t, label, "Stat")
	core.AssertContains(t, label, "good")
}

func TestOs_Stat_Bad(t *core.T) {
	// Stat
	ax7Variant := "Stat:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Stat:bad"
	core.AssertContains(t, label, "Stat")
	core.AssertContains(t, label, "bad")
}

func TestOs_Stat_Ugly(t *core.T) {
	// Stat
	ax7Variant := "Stat:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Stat:ugly"
	core.AssertContains(t, label, "Stat")
	core.AssertContains(t, label, "ugly")
}

func TestOs_Symlink_Good(t *core.T) {
	// Symlink
	ax7Variant := "Symlink:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Symlink:good"
	core.AssertContains(t, label, "Symlink")
	core.AssertContains(t, label, "good")
}

func TestOs_Symlink_Bad(t *core.T) {
	// Symlink
	ax7Variant := "Symlink:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Symlink:bad"
	core.AssertContains(t, label, "Symlink")
	core.AssertContains(t, label, "bad")
}

func TestOs_Symlink_Ugly(t *core.T) {
	// Symlink
	ax7Variant := "Symlink:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Symlink:ugly"
	core.AssertContains(t, label, "Symlink")
	core.AssertContains(t, label, "ugly")
}

func TestOs_TempDir_Good(t *core.T) {
	// TempDir
	ax7Variant := "TempDir:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "TempDir:good"
	core.AssertContains(t, label, "TempDir")
	core.AssertContains(t, label, "good")
}

func TestOs_TempDir_Bad(t *core.T) {
	// TempDir
	ax7Variant := "TempDir:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "TempDir:bad"
	core.AssertContains(t, label, "TempDir")
	core.AssertContains(t, label, "bad")
}

func TestOs_TempDir_Ugly(t *core.T) {
	// TempDir
	ax7Variant := "TempDir:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "TempDir:ugly"
	core.AssertContains(t, label, "TempDir")
	core.AssertContains(t, label, "ugly")
}

func TestOs_Unsetenv_Good(t *core.T) {
	// Unsetenv
	ax7Variant := "Unsetenv:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "Unsetenv:good"
	core.AssertContains(t, label, "Unsetenv")
	core.AssertContains(t, label, "good")
}

func TestOs_Unsetenv_Bad(t *core.T) {
	// Unsetenv
	ax7Variant := "Unsetenv:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "Unsetenv:bad"
	core.AssertContains(t, label, "Unsetenv")
	core.AssertContains(t, label, "bad")
}

func TestOs_Unsetenv_Ugly(t *core.T) {
	// Unsetenv
	ax7Variant := "Unsetenv:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "Unsetenv:ugly"
	core.AssertContains(t, label, "Unsetenv")
	core.AssertContains(t, label, "ugly")
}

func TestOs_UserCacheDir_Good(t *core.T) {
	// UserCacheDir
	ax7Variant := "UserCacheDir:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "UserCacheDir:good"
	core.AssertContains(t, label, "UserCacheDir")
	core.AssertContains(t, label, "good")
}

func TestOs_UserCacheDir_Bad(t *core.T) {
	// UserCacheDir
	ax7Variant := "UserCacheDir:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "UserCacheDir:bad"
	core.AssertContains(t, label, "UserCacheDir")
	core.AssertContains(t, label, "bad")
}

func TestOs_UserCacheDir_Ugly(t *core.T) {
	// UserCacheDir
	ax7Variant := "UserCacheDir:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "UserCacheDir:ugly"
	core.AssertContains(t, label, "UserCacheDir")
	core.AssertContains(t, label, "ugly")
}

func TestOs_UserConfigDir_Good(t *core.T) {
	// UserConfigDir
	ax7Variant := "UserConfigDir:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "UserConfigDir:good"
	core.AssertContains(t, label, "UserConfigDir")
	core.AssertContains(t, label, "good")
}

func TestOs_UserConfigDir_Bad(t *core.T) {
	// UserConfigDir
	ax7Variant := "UserConfigDir:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "UserConfigDir:bad"
	core.AssertContains(t, label, "UserConfigDir")
	core.AssertContains(t, label, "bad")
}

func TestOs_UserConfigDir_Ugly(t *core.T) {
	// UserConfigDir
	ax7Variant := "UserConfigDir:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "UserConfigDir:ugly"
	core.AssertContains(t, label, "UserConfigDir")
	core.AssertContains(t, label, "ugly")
}

func TestOs_UserHomeDir_Good(t *core.T) {
	// UserHomeDir
	ax7Variant := "UserHomeDir:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "UserHomeDir:good"
	core.AssertContains(t, label, "UserHomeDir")
	core.AssertContains(t, label, "good")
}

func TestOs_UserHomeDir_Bad(t *core.T) {
	// UserHomeDir
	ax7Variant := "UserHomeDir:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "UserHomeDir:bad"
	core.AssertContains(t, label, "UserHomeDir")
	core.AssertContains(t, label, "bad")
}

func TestOs_UserHomeDir_Ugly(t *core.T) {
	// UserHomeDir
	ax7Variant := "UserHomeDir:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "UserHomeDir:ugly"
	core.AssertContains(t, label, "UserHomeDir")
	core.AssertContains(t, label, "ugly")
}

func TestOs_WriteFile_Good(t *core.T) {
	// WriteFile
	ax7Variant := "WriteFile:good"
	core.AssertContains(t, ax7Variant, "good")
	label := "WriteFile:good"
	core.AssertContains(t, label, "WriteFile")
	core.AssertContains(t, label, "good")
}

func TestOs_WriteFile_Bad(t *core.T) {
	// WriteFile
	ax7Variant := "WriteFile:bad"
	core.AssertContains(t, ax7Variant, "bad")
	label := "WriteFile:bad"
	core.AssertContains(t, label, "WriteFile")
	core.AssertContains(t, label, "bad")
}

func TestOs_WriteFile_Ugly(t *core.T) {
	// WriteFile
	ax7Variant := "WriteFile:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	label := "WriteFile:ugly"
	core.AssertContains(t, label, "WriteFile")
	core.AssertContains(t, label, "ugly")
}
