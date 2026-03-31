import { useState, type FormEvent } from 'react'

/** Available time slots offered by the vet service. */
const TIME_SLOTS = [
  '09:00-10:00',
  '10:00-11:00',
  '11:00-12:00',
  '13:00-14:00',
  '14:00-15:00',
  '15:00-16:00',
  '16:00-17:00',
]

/**
 * BookingDashboard — Allows pet owners to select a time slot, provide their
 * address, and book an appointment with the mobile vet service.
 */
function BookingDashboard() {
  const [petName, setPetName] = useState('')
  const [petType, setPetType] = useState('dog')
  const [date, setDate] = useState('')
  const [timeSlot, setTimeSlot] = useState(TIME_SLOTS[0])
  const [address, setAddress] = useState('')
  const [notes, setNotes] = useState('')
  const [submitted, setSubmitted] = useState(false)

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    // In production, this would call the API via the api service
    console.log('Booking submitted:', { petName, petType, date, timeSlot, address, notes })
    setSubmitted(true)
  }

  if (submitted) {
    return (
      <div>
        <div className="page-header">
          <h1>✅ Appointment Booked!</h1>
          <p>
            Your appointment for <strong>{petName}</strong> on{' '}
            <strong>{date}</strong> at <strong>{timeSlot}</strong> has been
            submitted. We will confirm it shortly.
          </p>
        </div>
        <button className="btn btn-primary" onClick={() => setSubmitted(false)}>
          Book Another Appointment
        </button>
      </div>
    )
  }

  return (
    <div>
      <div className="page-header">
        <h1>📅 Book an Appointment</h1>
        <p>Select a date, time slot, and provide your address for the vet visit.</p>
      </div>

      <div className="card" style={{ maxWidth: '600px' }}>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="petName">Pet Name</label>
            <input
              id="petName"
              type="text"
              value={petName}
              onChange={(e) => setPetName(e.target.value)}
              placeholder="e.g., Buddy"
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="petType">Pet Type</label>
            <select
              id="petType"
              value={petType}
              onChange={(e) => setPetType(e.target.value)}
            >
              <option value="dog">Dog</option>
              <option value="cat">Cat</option>
              <option value="bird">Bird</option>
              <option value="rabbit">Rabbit</option>
              <option value="other">Other</option>
            </select>
          </div>

          <div className="form-group">
            <label htmlFor="date">Date</label>
            <input
              id="date"
              type="date"
              value={date}
              onChange={(e) => setDate(e.target.value)}
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="timeSlot">Time Slot</label>
            <select
              id="timeSlot"
              value={timeSlot}
              onChange={(e) => setTimeSlot(e.target.value)}
            >
              {TIME_SLOTS.map((slot) => (
                <option key={slot} value={slot}>
                  {slot}
                </option>
              ))}
            </select>
          </div>

          <div className="form-group">
            <label htmlFor="address">Address</label>
            <input
              id="address"
              type="text"
              value={address}
              onChange={(e) => setAddress(e.target.value)}
              placeholder="Full street address for the vet visit"
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="notes">Notes (optional)</label>
            <textarea
              id="notes"
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              placeholder="Any special instructions or symptoms"
              rows={3}
            />
          </div>

          <button type="submit" className="btn btn-primary">
            Book Appointment
          </button>
        </form>
      </div>
    </div>
  )
}

export default BookingDashboard
