package account

type Account struct {
	ID         int64   `json:"id"`
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	Balance    float64 `json:"balance"`
	CardNumber string  `json:"cardNumber"`
	Password   string  `json:"password"`
}

type Info struct {
	ID         int64   `json:"id"`
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	Balance    float64 `json:"balance"`
	CardNumber string  `json:"cardNumber"`
}

type Login struct {
	CardNumber string `json:"cardNumber" validate:"required" `
	Password   string `json:"password" validate:"required" `
}

type Registration struct {
	FirstName  string `json:"firstName" validate:"required" `
	LastName   string `json:"lastName" validate:"required" `
	CardNumber string `json:"cardNumber" validate:"required" `
	Password   string `json:"password" validate:"required" `
}

type ResponseBalance struct {
	Balance float64 `json:"amount" validate:"required" `
}
