package rcmd

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
