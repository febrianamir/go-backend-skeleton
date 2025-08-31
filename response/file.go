package response

type UploadFile struct {
	Filename string `json:"filename"`
	Filepath string `json:"filepath"`
	URL      string `json:"url"`
}
