package inventory

import (
	"testing"
)

func TestNewInventoryItem(t *testing.T) {
	item := NewInventoryItem("Amoxicillin", "Antibiotic for pets", CategoryMedicine, 12.99, 50, "tablet")

	if item.Name != "Amoxicillin" {
		t.Errorf("Expected name Amoxicillin, got %s", item.Name)
	}
	if item.Category != CategoryMedicine {
		t.Errorf("Expected category medicine, got %s", item.Category)
	}
	if item.Price != 12.99 {
		t.Errorf("Expected price 12.99, got %f", item.Price)
	}
	if item.Quantity != 50 {
		t.Errorf("Expected quantity 50, got %d", item.Quantity)
	}
	if item.Unit != "tablet" {
		t.Errorf("Expected unit tablet, got %s", item.Unit)
	}
}

func TestInventoryItem_IsInStock(t *testing.T) {
	item := NewInventoryItem("Dog Food", "Premium kibble", CategoryFood, 25.00, 10, "kg")
	if !item.IsInStock() {
		t.Error("Item with quantity 10 should be in stock")
	}

	item.Quantity = 0
	if item.IsInStock() {
		t.Error("Item with quantity 0 should not be in stock")
	}
}

func TestInventoryItem_Deduct(t *testing.T) {
	item := NewInventoryItem("Vitamin C", "Pet supplement", CategorySupplement, 8.50, 20, "bottle")

	if !item.Deduct(5) {
		t.Error("Expected successful deduction of 5 from 20")
	}
	if item.Quantity != 15 {
		t.Errorf("Expected quantity 15 after deducting 5, got %d", item.Quantity)
	}

	if item.Deduct(100) {
		t.Error("Expected failed deduction when amount exceeds stock")
	}
	if item.Quantity != 15 {
		t.Errorf("Quantity should remain 15 after failed deduction, got %d", item.Quantity)
	}
}

func TestInventoryItem_Restock(t *testing.T) {
	item := NewInventoryItem("Flea Drops", "Anti-flea treatment", CategoryMedicine, 15.00, 5, "bottle")

	item.Restock(10)
	if item.Quantity != 15 {
		t.Errorf("Expected quantity 15 after restocking 10, got %d", item.Quantity)
	}
}

func TestValidCategory(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"medicine", true},
		{"food", true},
		{"supplement", true},
		{"electronics", false},
		{"", false},
	}

	for _, tt := range tests {
		got := ValidCategory(tt.input)
		if got != tt.want {
			t.Errorf("ValidCategory(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
