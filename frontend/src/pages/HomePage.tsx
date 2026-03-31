import { Link } from 'react-router-dom'

/**
 * HomePage — Landing page for the Haoma Vet mobile veterinary care service.
 * Highlights core features: booking, inventory, and the mobile van concept.
 */
function HomePage() {
  return (
    <div>
      <div className="page-header">
        <h1>🐾 Mobile Veterinary Care, At Your Doorstep</h1>
        <p>
          Haoma Vet brings a fully equipped veterinary clinic to your home.
          Book an appointment, and our van with a sterile operating room arrives
          at your address.
        </p>
      </div>

      <div className="feature-grid">
        <div className="card">
          <h3>📅 Easy Booking</h3>
          <p>
            Select an available time slot, provide your address, and a licensed
            veterinarian will come to you in our mobile clinic.
          </p>
          <Link to="/booking" className="btn btn-primary" style={{ marginTop: '1rem' }}>
            Book Now
          </Link>
        </div>

        <div className="card">
          <h3>🚐 Mobile Clinic</h3>
          <p>
            Our van is equipped with a sterile operating room, diagnostic tools,
            and everything your pet needs — no stressful trips to the clinic.
          </p>
        </div>

        <div className="card">
          <h3>💊 Pet Supplies</h3>
          <p>
            Browse and order medicines, premium pet food, and supplements
            directly through our platform. Delivered with your appointment
            or separately.
          </p>
        </div>
      </div>
    </div>
  )
}

export default HomePage
