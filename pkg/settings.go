package pkg

import "github.com/pelletier/go-toml"

var (
	Version string
	ReadSymLink = false
	Verbose = false
	HashAlgorithm = "sha256"
	RepoPath = ""
	BufferSize = 4*1024*1024
)

type settings struct {
	Version       string `toml:version`
	HashAlgorithm string `toml:hashAlgorithm`
	IsLocked      bool   `toml:IsLocked`
}

func generateSettingsToml(isLocked bool) ([]byte, error) {
	return toml.Marshal(settings{
		Version:       Version,
		HashAlgorithm: HashAlgorithm,
		IsLocked:      isLocked,
	})
}
