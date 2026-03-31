package appointment

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewAppointment(t *testing.T) {
	userID := uuid.New()
	date := time.Date(2026, 4, 15, 10, 0, 0, 0, time.UTC)
	appt := NewAppointment(userID, "Buddy", "dog", date, "10:00-11:00", "456 Oak Ave", "Annual checkup")

	if appt.UserID != userID {
		t.Errorf("Expected user ID %v, got %v", userID, appt.UserID)
	}
	if appt.PetName != "Buddy" {
		t.Errorf("Expected pet name Buddy, got %s", appt.PetName)
	}
	if appt.PetType != "dog" {
		t.Errorf("Expected pet type dog, got %s", appt.PetType)
	}
	if appt.TimeSlot != "10:00-11:00" {
		t.Errorf("Expected time slot 10:00-11:00, got %s", appt.TimeSlot)
	}
	if appt.Status != StatusPending {
		t.Errorf("Expected status pending, got %s", appt.Status)
	}
}

func TestAppointment_StatusTransitions(t *testing.T) {
	appt := NewAppointment(uuid.New(), "Max", "cat", time.Now(), "09:00-10:00", "789 Elm St", "")

	if !appt.IsCancellable() {
		t.Error("New appointment should be cancellable")
	}

	appt.Confirm()
	if appt.Status != StatusConfirmed {
		t.Errorf("Expected status confirmed, got %s", appt.Status)
	}
	if !appt.IsCancellable() {
		t.Error("Confirmed appointment should still be cancellable")
	}

	appt.Complete()
	if appt.Status != StatusCompleted {
		t.Errorf("Expected status completed, got %s", appt.Status)
	}
	if appt.IsCancellable() {
		t.Error("Completed appointment should not be cancellable")
	}
}

func TestAppointment_Cancel(t *testing.T) {
	appt := NewAppointment(uuid.New(), "Luna", "bird", time.Now(), "14:00-15:00", "321 Pine Rd", "")

	appt.Cancel()
	if appt.Status != StatusCancelled {
		t.Errorf("Expected status cancelled, got %s", appt.Status)
	}
	if appt.IsCancellable() {
		t.Error("Cancelled appointment should not be cancellable")
	}
}
