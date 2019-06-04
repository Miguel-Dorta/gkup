package pkg

import "github.com/pelletier/go-toml"

var (
	Version string
	ReadSymLink = false
	Verbose = false
	HashAlgorithm = "sha256"
	RepoPath = ""
	BufferSize = 4*1024*1024
	OmitErrors = false
)

type settings struct {
	Version       string `toml:version`
	HashAlgorithm string `toml:hashAlgorithm`
}

func generateSettingsToml() (data []byte) {
	data, _ = toml.Marshal(settings{
		Version:       Version,
		HashAlgorithm: HashAlgorithm,
	})
	return
}
