package request

type GetUserVerification struct {
	UserID   uint
	Code     string
	Preloads []string
}
