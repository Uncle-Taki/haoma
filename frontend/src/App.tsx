import { Routes, Route, Link } from 'react-router-dom'
import HomePage from './pages/HomePage'
import BookingDashboard from './pages/BookingDashboard'
import AdminDashboard from './pages/AdminDashboard'
import './App.css'

function App() {
  return (
    <div className="app">
      <nav className="navbar">
        <Link to="/" className="nav-brand">🐾 Haoma Vet</Link>
        <div className="nav-links">
          <Link to="/">Home</Link>
          <Link to="/booking">Book Appointment</Link>
          <Link to="/admin">Vet Dashboard</Link>
        </div>
      </nav>

      <main className="main-content">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/booking" element={<BookingDashboard />} />
          <Route path="/admin" element={<AdminDashboard />} />
        </Routes>
      </main>

      <footer className="footer">
        <p>&copy; {new Date().getFullYear()} Haoma Vet — Mobile Veterinary Care</p>
      </footer>
    </div>
  )
}

export default App
