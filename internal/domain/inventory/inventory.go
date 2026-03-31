package inventory

import (
	"time"

	"github.com/google/uuid"
)

// Category represents the type of inventory item.
type Category string

const (
	CategoryMedicine   Category = "medicine"
	CategoryFood       Category = "food"
	CategorySupplement Category = "supplement"
)

// InventoryItem represents a product available in the mobile vet van.
type InventoryItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Category    Category  `gorm:"not null;index" json:"category"`
	Price       float64   `gorm:"not null" json:"price"`
	Quantity    int       `gorm:"not null;default:0" json:"quantity"`
	Unit        string    `gorm:"not null" json:"unit"` // e.g., "tablet", "kg", "bottle"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewInventoryItem creates a new inventory item.
func NewInventoryItem(name, description string, category Category, price float64, quantity int, unit string) *InventoryItem {
	return &InventoryItem{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Category:    category,
		Price:       price,
		Quantity:    quantity,
		Unit:        unit,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// IsInStock returns true if the item has available quantity.
func (i *InventoryItem) IsInStock() bool {
	return i.Quantity > 0
}

// Deduct reduces the quantity by the given amount. Returns false if insufficient stock.
func (i *InventoryItem) Deduct(amount int) bool {
	if i.Quantity < amount {
		return false
	}
	i.Quantity -= amount
	i.UpdatedAt = time.Now()
	return true
}

// Restock adds the given amount to the current quantity.
func (i *InventoryItem) Restock(amount int) {
	i.Quantity += amount
	i.UpdatedAt = time.Now()
}

// ValidCategory checks if a string is a valid inventory category.
func ValidCategory(cat string) bool {
	switch Category(cat) {
	case CategoryMedicine, CategoryFood, CategorySupplement:
		return true
	}
	return false
}
