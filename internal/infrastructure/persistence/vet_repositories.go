package persistence

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"haoma/internal/domain/appointment"
	"haoma/internal/domain/inventory"
	"haoma/internal/domain/user"
)

// VetUserRepository implements persistence for vet service users.
type VetUserRepository struct {
	db *gorm.DB
}

func NewVetUserRepository(db *gorm.DB) *VetUserRepository {
	return &VetUserRepository{db: db}
}

func (r *VetUserRepository) Save(u *user.User) error {
	return r.db.Create(u).Error
}

func (r *VetUserRepository) FindByID(id uuid.UUID) (*user.User, error) {
	var u user.User
	err := r.db.First(&u, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &u, err
}

func (r *VetUserRepository) FindByEmail(email string) (*user.User, error) {
	var u user.User
	err := r.db.Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &u, err
}

// AppointmentRepository implements persistence for appointments.
type AppointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func (r *AppointmentRepository) Save(appt *appointment.Appointment) error {
	return r.db.Create(appt).Error
}

func (r *AppointmentRepository) FindByID(id uuid.UUID) (*appointment.Appointment, error) {
	var appt appointment.Appointment
	err := r.db.First(&appt, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("appointment not found")
	}
	return &appt, err
}

func (r *AppointmentRepository) FindByUserID(userID uuid.UUID) ([]appointment.Appointment, error) {
	var appointments []appointment.Appointment
	err := r.db.Where("user_id = ?", userID).Order("date DESC").Find(&appointments).Error
	return appointments, err
}

func (r *AppointmentRepository) FindByDate(date time.Time) ([]appointment.Appointment, error) {
	var appointments []appointment.Appointment
	startOfDay := date.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)
	err := r.db.Where("date >= ? AND date < ?", startOfDay, endOfDay).
		Order("time_slot ASC").
		Find(&appointments).Error
	return appointments, err
}

func (r *AppointmentRepository) Update(appt *appointment.Appointment) error {
	return r.db.Save(appt).Error
}

// InventoryRepository implements persistence for inventory items.
type InventoryItemRepository struct {
	db *gorm.DB
}

func NewInventoryItemRepository(db *gorm.DB) *InventoryItemRepository {
	return &InventoryItemRepository{db: db}
}

func (r *InventoryItemRepository) Save(item *inventory.InventoryItem) error {
	return r.db.Create(item).Error
}

func (r *InventoryItemRepository) FindByID(id uuid.UUID) (*inventory.InventoryItem, error) {
	var item inventory.InventoryItem
	err := r.db.First(&item, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("inventory item not found")
	}
	return &item, err
}

func (r *InventoryItemRepository) FindAll() ([]inventory.InventoryItem, error) {
	var items []inventory.InventoryItem
	err := r.db.Order("category, name").Find(&items).Error
	return items, err
}

func (r *InventoryItemRepository) FindByCategory(category inventory.Category) ([]inventory.InventoryItem, error) {
	var items []inventory.InventoryItem
	err := r.db.Where("category = ?", category).Order("name").Find(&items).Error
	return items, err
}

func (r *InventoryItemRepository) Update(item *inventory.InventoryItem) error {
	return r.db.Save(item).Error
}
