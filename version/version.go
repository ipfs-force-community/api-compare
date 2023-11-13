package version

var (
	CurrentCommit string

	Version = "0.0.0"
)

func UserVersion() string {
	return Version + CurrentCommit
}
