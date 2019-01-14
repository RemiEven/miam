package model

// Product is a product with id/name
type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// EditableProduct is an editable product
type EditableProduct struct {
	Name string `json:"name"`
}
