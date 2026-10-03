// Package create provides functionality for creating products.
package create

// CreatorProductCommand contains the data required to create a product.
type CreatorProductCommand struct {
	Name  string  `json:"name" validate:"required,min=2,max=100"`
	Price float64 `json:"price" validate:"required,gt=0"`
}
