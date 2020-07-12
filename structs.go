package rcmd

// CmdResult stores information about the executed cmd
type CmdResult struct {
	Stdout   string `json:"stdout,omitempty"`
	Stderr   string `json:"stderr,omitempty"`
	ExitCode int    `json:"exit_code,omitempty"`
}

// ExecSettings controls settings related to R execution
type ExecSettings struct {
	WorkDir string `json:"work_dir,omitempty"`
}

// RVersion contains information about the R version
type RVersion struct {
	Major int
	Minor int
	Patch int
}

// RSettings controls settings related to managing libraries
type RSettings struct {
	Version  RVersion `json:"r_version,omitempty"`
	LibPaths []string `json:"lib_paths,omitempty"`
	RPath    string   `json:"rpath,omitempty"`
	EnvVars  NvpList  `json:"env_vars,omitempty"`
	Platform string   `json:"platform,omitempty"`
}

// InstallArgs represents the installation arguments R CMD INSTALL can consume
type InstallArgs struct {
	Clean          bool `rcmd:"clean"`
	Preclean       bool `rcmd:"preclean"`
	Debug          bool `rcmd:"debug"`
	NoConfigure    bool `rcmd:"no-configure"`
	Example        bool `rcmd:"example"`
	Fake           bool `rcmd:"fake"`
	Build          bool `rcmd:"build"`
	InstallTests   bool `rcmd:"install-tests"`
	NoMultiarch    bool `rcmd:"no-multiarch"`
	WithKeepSource bool `rcmd:"with-keep.source"`
	ByteCompile    bool `rcmd:"byte-compile"`
	NoTestLoad     bool `rcmd:"no-test-load"`
	NoCleanOnError bool `rcmd:"no-clean-on-error"`
	//set
	Library string `rcmd:"library=%s,fmt"`
}

// PackageCache provides metadata about the package cache
// Each repository should be a subfolder from the BaseDir
// with separate folders for binary and source packages
type PackageCache struct {
	BaseDir string
}

// InstallRequest provides information about the installation request
type InstallRequest struct {
	Package      string
	Path         string
	IsBinary     bool
	Cache        PackageCache
	Args         InstallArgs
	ExecSettings ExecSettings
	RSettings    RSettings
}

// InstallResult provides information about the Job in the queue
type InstallResult struct {
	Result  CmdResult
	Package string
}

// Nvp name-value pair, each of type string
type Nvp struct {
	Name  string `json:"global_env_vars_name,omitempty"`
	Value string `json:"global_env_vars_value,omitempty"`
}

// NvpList is a slice of Nvp. The slice maintains consistent ordering of the Nvp objects
type NvpList struct {
	Pairs []Nvp `json:"pairs,omitempty"`
}
