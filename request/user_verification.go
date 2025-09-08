package request

type GetUserVerification struct {
	Type     string
	UserID   uint
	Code     string
	IsUsed   *bool
	Preloads []string
}
