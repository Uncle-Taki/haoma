/**
 * AdminDashboard — Vet/admin view showing today's route and appointments.
 * Displays scheduled appointments and allows managing them.
 */
function AdminDashboard() {
  // Placeholder data — in production these would come from the API
  const todayAppointments = [
    {
      id: '1',
      time: '09:00-10:00',
      petName: 'Buddy',
      petType: 'Dog',
      address: '123 Maple Street',
      status: 'confirmed',
    },
    {
      id: '2',
      time: '10:00-11:00',
      petName: 'Whiskers',
      petType: 'Cat',
      address: '456 Oak Avenue',
      status: 'pending',
    },
    {
      id: '3',
      time: '13:00-14:00',
      petName: 'Tweety',
      petType: 'Bird',
      address: '789 Pine Road',
      status: 'confirmed',
    },
  ]

  return (
    <div>
      <div className="page-header">
        <h1>🩺 Vet Dashboard</h1>
        <p>Today's scheduled appointments and route overview.</p>
      </div>

      <div className="card" style={{ marginBottom: '2rem' }}>
        <h3>📊 Today's Summary</h3>
        <p>
          <strong>{todayAppointments.length}</strong> appointments scheduled
          &nbsp;|&nbsp;
          <strong>{todayAppointments.filter((a) => a.status === 'confirmed').length}</strong> confirmed
          &nbsp;|&nbsp;
          <strong>{todayAppointments.filter((a) => a.status === 'pending').length}</strong> pending
        </p>
      </div>

      <h2 style={{ marginBottom: '1rem' }}>Appointments</h2>
      <div className="appointment-list">
        {todayAppointments.map((appt) => (
          <div key={appt.id} className="appointment-card">
            <span className="time">{appt.time}</span>
            <div className="details">
              <strong>{appt.petName}</strong> ({appt.petType})<br />
              <small>{appt.address}</small>
            </div>
            <span className={`status-badge status-${appt.status}`}>
              {appt.status}
            </span>
          </div>
        ))}
      </div>
    </div>
  )
}

export default AdminDashboard
