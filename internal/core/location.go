package core

// enum values for location kind
const (
	KindS3   = "s3"
	KindFTP  = "ftp"
	KindFTPS = "ftps"
	KindSFTP = "sftp"
)

type Location struct {
	// readonly (from database, after creation)
	ID int64

	Kind    string
	Info    string
	Project Project
}

func NewLocation(kind, info string, project Project) Location {
	location := Location{
		Kind:    kind,
		Info:    info,
		Project: project,
	}
	return location
}

type LocationStorage interface {
	Create(location *Location) error
	Read(id int64) (Location, error)
	Update(location Location) error
	Delete(location Location) error

	ReadAll() ([]Location, error)
}

type S3Info struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	BucketName      string `json:"bucket_name"`
}

type FTPInfo struct {
	Endpoint string `json:"endpoint"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type FTPSInfo struct {
	Endpoint string `json:"endpoint"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SFTPInfo struct {
	Endpoint   string `json:"endpoint"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
}
