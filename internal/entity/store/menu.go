/* Package store contains store object models and define interfaces */
package store

// Product an item or a set of items that the store sells
type Product struct {
	Name        string    `json:"name" bson:"name" validate:"required,max=50"`
	Description string    `json:"description" bson:"desc" validate:"max=1000"`
	Category    string    `json:"category" bson:"cat" validate:"required, max=50"`
	Price       int       `json:"price,omitempty" bson:"price" validate:"required_without:Variants"`
	Variants    []Variant `json:"variants,omitempty" bson:"var,omitempty" validate:"required_without:Price"`
}

// Variant variant of a product, for example: different sizes, cold/warm
type Variant struct {
	Name  string `json:"name" bson:"name" validate:"required,max=5"`
	Price int    `json:"price" bson:"price" validate:"required"`
}
