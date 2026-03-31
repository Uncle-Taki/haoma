package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"haoma/internal/domain/appointment"
	"haoma/internal/domain/inventory"
	"haoma/internal/domain/user"
)

// VetService handles veterinary care business logic for bookings and inventory.
type VetService struct {
	userRepo        VetUserRepository
	appointmentRepo AppointmentRepository
	inventoryRepo   InventoryRepository
}

// VetUserRepository defines persistence operations for vet service users.
type VetUserRepository interface {
	Save(user *user.User) error
	FindByID(id uuid.UUID) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
}

// AppointmentRepository defines persistence operations for appointments.
type AppointmentRepository interface {
	Save(appt *appointment.Appointment) error
	FindByID(id uuid.UUID) (*appointment.Appointment, error)
	FindByUserID(userID uuid.UUID) ([]appointment.Appointment, error)
	FindByDate(date time.Time) ([]appointment.Appointment, error)
	Update(appt *appointment.Appointment) error
}

// InventoryRepository defines persistence operations for inventory items.
type InventoryRepository interface {
	Save(item *inventory.InventoryItem) error
	FindByID(id uuid.UUID) (*inventory.InventoryItem, error)
	FindAll() ([]inventory.InventoryItem, error)
	FindByCategory(category inventory.Category) ([]inventory.InventoryItem, error)
	Update(item *inventory.InventoryItem) error
}

// NewVetService creates a new VetService with the provided repositories.
func NewVetService(
	userRepo VetUserRepository,
	appointmentRepo AppointmentRepository,
	inventoryRepo InventoryRepository,
) *VetService {
	return &VetService{
		userRepo:        userRepo,
		appointmentRepo: appointmentRepo,
		inventoryRepo:   inventoryRepo,
	}
}

// CreateUser registers a new user for the vet service.
func (s *VetService) CreateUser(u *user.User) error {
	_, err := s.userRepo.FindByEmail(u.Email)
	if err == nil {
		return errors.New("user already exists")
	}
	return s.userRepo.Save(u)
}

// GetUserByID retrieves a user by their ID.
func (s *VetService) GetUserByID(id uuid.UUID) (*user.User, error) {
	return s.userRepo.FindByID(id)
}

// GetUserByEmail retrieves a user by their email address.
func (s *VetService) GetUserByEmail(email string) (*user.User, error) {
	return s.userRepo.FindByEmail(email)
}

// BookAppointment creates a new appointment for a user.
func (s *VetService) BookAppointment(userID uuid.UUID, petName, petType string, date time.Time, timeSlot, address, notes string) (*appointment.Appointment, error) {
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	appt := appointment.NewAppointment(userID, petName, petType, date, timeSlot, address, notes)
	if err := s.appointmentRepo.Save(appt); err != nil {
		return nil, err
	}

	return appt, nil
}

// GetUserAppointments returns all appointments for a given user.
func (s *VetService) GetUserAppointments(userID uuid.UUID) ([]appointment.Appointment, error) {
	return s.appointmentRepo.FindByUserID(userID)
}

// GetAppointmentsByDate returns all appointments for a given date (for vet dashboard).
func (s *VetService) GetAppointmentsByDate(date time.Time) ([]appointment.Appointment, error) {
	return s.appointmentRepo.FindByDate(date)
}

// CancelAppointment cancels an existing appointment.
func (s *VetService) CancelAppointment(appointmentID uuid.UUID) error {
	appt, err := s.appointmentRepo.FindByID(appointmentID)
	if err != nil {
		return errors.New("appointment not found")
	}

	if !appt.IsCancellable() {
		return errors.New("appointment cannot be cancelled")
	}

	appt.Cancel()
	return s.appointmentRepo.Update(appt)
}

// GetInventory returns all inventory items, optionally filtered by category.
func (s *VetService) GetInventory(category string) ([]inventory.InventoryItem, error) {
	if category != "" {
		if !inventory.ValidCategory(category) {
			return nil, errors.New("invalid category")
		}
		return s.inventoryRepo.FindByCategory(inventory.Category(category))
	}
	return s.inventoryRepo.FindAll()
}

// GetInventoryItem returns a single inventory item by ID.
func (s *VetService) GetInventoryItem(id uuid.UUID) (*inventory.InventoryItem, error) {
	return s.inventoryRepo.FindByID(id)
}
