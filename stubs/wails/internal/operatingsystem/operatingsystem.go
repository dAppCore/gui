package operatingsystem

// OS is the minimal host metadata exposed by the Wails environment API.
type OS struct {
	Name    string
	Version string
	Build   string
}
