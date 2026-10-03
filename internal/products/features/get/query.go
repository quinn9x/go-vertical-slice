// Package get defines the query structure for retrieving product information.
package get

// GetterProductQuery represents the query structure for retrieving a product by its ID.
type GetterProductQuery struct {
	ID string `json:"id"`
}
