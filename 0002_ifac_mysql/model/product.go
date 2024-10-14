package model

type Product struct {
	ID           int    `json:"id"`
	ProductName  string `json:"productName"`
	ProductPrice int    `json:"productPrice"`

	Categories []Category
}
