package rcmd

// CmdResult stores information about the executed cmd.
type CmdResult struct {
	Stdout   string `json:"stdout,omitempty"`
	Stderr   string `json:"stderr,omitempty"`
	ExitCode int    `json:"exit_code,omitempty"`
}

// ExecSettings controls settings related to R execution.
type ExecSettings struct {
	WorkDir string `json:"work_dir,omitempty"`
}

// RVersion contains information about the R version.
type RVersion struct {
	Major int
	Minor int
	Patch int
}

// RSettings controls settings related to managing libraries.
// If AsUser is set, R will be run as the user would launch from their normal session.
// with no interception/injection of library paths or environment variables for R_LIBS_SITE and R_LIBS_USER.
type RSettings struct {
	AsUser   bool     `json:"as_user,omitempty"`
	Version  RVersion `json:"r_version,omitempty"`
	LibPaths []string `json:"lib_paths,omitempty"`
	RPath    string   `json:"rpath,omitempty"`
	Platform string   `json:"platform,omitempty"`
}

// InstallArgs represents the installation arguments R CMD INSTALL can consume.
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
	// set
	Library string `rcmd:"library=%s,fmt"`
}

type CheckArgs struct {
	NoClean          bool `rcmd:"no-clean"`
	NoInstall        bool `rcmd:"no-install"`
	NoTests          bool `rcmd:"no-tests"`
	NoManual         bool `rcmd:"no-manual"`
	NoVignettes      bool `rcmd:"no-vignettes"`
	NoBuildVignettes bool `rcmd:"no-build-vignettes"`
	IgnoreVignettes  bool `rcmd:"ignore-vignettes"`
	InstallTests     bool `rcmd:"install-tests"`
	Multiarch        bool `rcmd:"multiarch"`
	NoMultiarch      bool `rcmd:"no-multiarch"`
	AsCran           bool `rcmd:"as-cran"`
	//  Output directory for output, default is current directory.
	//           Logfiles, R output, etc. will be placed in 'pkg.Rcheck'
	//           in this directory, where 'pkg' is the name of the
	//           checked package
	Output string `rcmd:"output=%s,fmt"`
	//  library directory used for test installation of packages (default is outdir)
	Library string `rcmd:"library=%s,fmt"`
}
