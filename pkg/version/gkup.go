package version

// GkupVersion represents the version of the specific release of gkup.
var GkupVersion Version

// GkupVersionStr represents the string of the version of the specific release of gkup.
// It's defined at compile time.
var GkupVersionStr string

func init() {
	var err error
	GkupVersion, err = Parse(GkupVersionStr)
	if err != nil {
		panic(err.Error() + ": " + GkupVersionStr)
	}
}
