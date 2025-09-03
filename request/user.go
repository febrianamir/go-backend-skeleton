package request

type GetUsers struct {
	BasePaginateRequest
	Preloads []string
}

func (query *GetUsers) GetOrderQuery() string {
	fieldMap := map[string]string{
		"name":       "name",
		"email":      "email",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}
	return buildOrderQuery(query.Sort, fieldMap)
}

type GetUser struct {
	ID       uint
	Name     string
	Email    string
	Preloads []string
}

type CreateUser struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type UpdateUser struct {
	ID          uint
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type DeleteUser struct {
	ID uint
}
