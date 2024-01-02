package config

import "flag"

// CLI flags are important in order to set variables when running tests in pipelines,
// for example: passing the 'profile' flag in order to load different configuration files for local / CI pipeline test runs
// but reading CLI flags in tests can apparently only be done when they are declared inside the test file itself
// a way that will work is to define them in the MainTest and pass them to the setup method, but that is extra effort for the user
// we will keep these defaults here for now until we find a solution
var baseDir = flag.String("clm-basedir", defaultBaseDir, "Base directory where to look for all files.")
var configFile = flag.String("clm-config", defaultConfigFile, "Specific configuration file.")
var activeProfile = flag.String("clm-profile", defaultProfile, "Active test profile. When specified here, it will overwrite the value from the configuration.")

func (config *Config) overwriteWithCliFlags() {

}
