package request

type GetUserVerification struct {
	UserID   uint
	Preloads []string
}
