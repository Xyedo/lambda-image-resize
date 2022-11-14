package appcnf

type Config struct {
	Region     string
	LocalDir   string
	BucketName string
	UploadDir  string
}

const (
	region    = "ap-southeast-1"
	localDir  = "/tmp"
	bucket    = "blindate-bucket"
	uploadDir = "profile-picture-resized"
)

func New() Config {
	return Config{
		Region:     region,
		LocalDir:   localDir,
		BucketName: bucket,
		UploadDir:  uploadDir,
	}
}
