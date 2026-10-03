// Package update provides the update functionality for products.
package update

// UpdaterProductCommand represents the command to update a product's details.
type UpdaterProductCommand struct {
	ID      string  `json:"-"`
	Name    string  `json:"name" validate:"required,min=2,max=100"`
	Price   float64 `json:"price" validate:"required,gt=0"`
	Version int64   `json:"version" validate:"required,min=1"`
}
