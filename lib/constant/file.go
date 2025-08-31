package constant

type UploadFileProps struct {
	// FormValue is request from form which contain file.
	FormValue string

	// AllowedExtensions is list allowed extensions in lower case.
	AllowedExtensions []string

	// MaxFilesize is maximum filesize. in MB.
	MaxFilesize int64
}

var MapUploadFileProps = map[string]UploadFileProps{
	"image": {
		FormValue:         "file",
		AllowedExtensions: []string{"jpg", "jpeg", "png", "webp"},
		MaxFilesize:       2,
	},
}
