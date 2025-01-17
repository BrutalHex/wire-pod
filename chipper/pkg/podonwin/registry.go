package podonwin

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows/registry"
)

// `Software\Microsoft\Windows\CurrentVersion\Uninstall\wire-pod`

// DisplayIcon string (path)
// DisplayVersion string (v1.0.0)
// Publisher string (github.com/BrutalHex/wire-pod)
// UninstallString (is.Where + uninstall.exe)
// InstallLocation (is.Where)
var Win_UninstallKeyKey = registry.LOCAL_MACHINE
var Win_UninstallKeyPath = `Software\Microsoft\Windows\CurrentVersion\Uninstall\wire-pod`
var Win_UninstallInstallerPerms uint32 = registry.ALL_ACCESS

// shouldn't ever need to access uninstall path in pod software, but go off i guess
var Win_UninstallPodPerms uint32 = registry.QUERY_VALUE

// InstallPath string (is.Where)
// PodVersion string (v1.0.0)
// LastRunningPID int (for wire-pod runtime, installer shouldn't touch this)
var Win_SoftwareKeyKey = registry.CURRENT_USER
var Win_SoftwarePodPerms uint32 = registry.READ | registry.WRITE
var Win_SoftwareInstallerPerms uint32 = registry.ALL_ACCESS
var Win_SoftwareKeyPath = `Software\wire-pod`

// key, _ := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
var Win_RunAtStartupKeyKey = registry.CURRENT_USER
var Win_RunAtStartupPerms uint32 = registry.READ | registry.WRITE
var Win_RunAtStartupKeyPath = `Software\Microsoft\Windows\CurrentVersion\Run`

var nonInitedError = "you must run podonwin.Init()"

type KeyInfo struct {
	Key     registry.Key
	Perms   uint32
	KeyPath string
}

var SoftwareKey KeyInfo
var UninstallKey KeyInfo
var StartupRunKey KeyInfo
var Inited bool
var IsInstaller bool

func Init() {
	if CheckIfElevated() {
		IsInstaller = true
	} else {
		IsInstaller = false
	}
	SoftwareKey = KeyInfo{
		Key:     Win_SoftwareKeyKey,
		KeyPath: Win_SoftwareKeyPath,
	}
	UninstallKey = KeyInfo{
		Key:     Win_UninstallKeyKey,
		KeyPath: Win_UninstallKeyPath,
	}
	StartupRunKey = KeyInfo{
		Key:     Win_RunAtStartupKeyKey,
		KeyPath: Win_RunAtStartupKeyPath,
		Perms:   Win_RunAtStartupPerms,
	}
	if IsInstaller {
		SoftwareKey.Perms = Win_SoftwareInstallerPerms
		UninstallKey.Perms = Win_UninstallInstallerPerms
	} else {
		SoftwareKey.Perms = Win_SoftwarePodPerms
		UninstallKey.Perms = Win_UninstallPodPerms
	}
	Inited = true
}

func DeleteEverythingFromRegistry() error {
	if !Inited {
		return fmt.Errorf(nonInitedError)
	}
	if !IsInstaller {
		return fmt.Errorf("must be run from installer")
	}
	DeleteRegistryKey(SoftwareKey)
	DeleteRegistryKey(UninstallKey)
	DeleteRegistryValue(StartupRunKey, "wire-pod")
	return nil
}

func DeleteRegistryKey(keyInfo KeyInfo) error {
	if !Inited {
		return fmt.Errorf(nonInitedError)
	}
	err := registry.DeleteKey(keyInfo.Key, keyInfo.KeyPath)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRegistryValue(keyInfo KeyInfo, key string) error {
	if !Inited {
		return fmt.Errorf(nonInitedError)
	}
	k, err := registry.OpenKey(keyInfo.Key, keyInfo.KeyPath, keyInfo.Perms)
	if err != nil {
		return err
	}
	defer k.Close()
	err = k.DeleteValue(key)
	if err != nil {
		return err
	}
	return nil
}

func UpdateRegistryValueString(keyInfo KeyInfo, key string, value string) error {
	if !Inited {
		return fmt.Errorf(nonInitedError)
	}
	k, _, err := registry.CreateKey(keyInfo.Key, keyInfo.KeyPath, keyInfo.Perms)
	if err != nil {
		return err
	}
	defer k.Close()
	err = k.SetStringValue(key, value)
	if err != nil {
		return err
	}
	return nil
}

func GetRegistryValueString(keyInfo KeyInfo, key string) (string, error) {
	if !Inited {
		return "", fmt.Errorf(nonInitedError)
	}
	k, err := registry.OpenKey(keyInfo.Key, keyInfo.KeyPath, keyInfo.Perms)
	if err != nil {
		return "", err
	}
	defer k.Close()
	val, _, err := k.GetStringValue(key)
	if err != nil {
		return "", err
	}
	return val, nil
}

func UpdateRegistryValueInt(keyInfo KeyInfo, key string, value int) error {
	if !Inited {
		return fmt.Errorf(nonInitedError)
	}
	k, _, err := registry.CreateKey(keyInfo.Key, keyInfo.KeyPath, keyInfo.Perms)
	if err != nil {
		return err
	}
	defer k.Close()
	err = k.SetQWordValue(key, uint64(value))
	if err != nil {
		return err
	}
	return nil
}

func GetRegistryValueInt(keyInfo KeyInfo, key string) (int, error) {
	if !Inited {
		return 0, fmt.Errorf(nonInitedError)
	}
	k, err := registry.OpenKey(keyInfo.Key, keyInfo.KeyPath, keyInfo.Perms)
	if err != nil {
		return 0, err
	}
	defer k.Close()
	val, _, err := k.GetIntegerValue(key)
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

func CheckIfElevated() bool {
	drv, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		drv.Close()
		return false
	}
	drv.Close()
	return true
}
