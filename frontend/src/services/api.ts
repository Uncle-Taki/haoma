/**
 * API service — Centralized HTTP client for communicating with the Go backend.
 * All vet service API calls are routed through /api/v1/vet/.
 */

const API_BASE = import.meta.env.VITE_API_BASE_URL || '/api/v1/vet'

/** Generic fetch wrapper with JSON parsing and error handling. */
async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const token = localStorage.getItem('token')

  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...options.headers,
  }

  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers,
  })

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Request failed' }))
    throw new Error(error.error || `HTTP ${response.status}`)
  }

  return response.json()
}

// --- Auth ---

export interface LoginPayload {
  email: string
  password: string
}

export interface SignupPayload {
  name: string
  email: string
  password: string
  phone?: string
  address?: string
}

export const authService = {
  signup: (data: SignupPayload) =>
    request('/auth/signup', { method: 'POST', body: JSON.stringify(data) }),

  login: (data: LoginPayload) =>
    request('/auth/login', { method: 'POST', body: JSON.stringify(data) }),

  getProfile: () => request('/auth/profile'),
}

// --- Appointments ---

export interface BookAppointmentPayload {
  pet_name: string
  pet_type: string
  date: string
  time_slot: string
  address: string
  notes?: string
}

export const appointmentService = {
  book: (data: BookAppointmentPayload) =>
    request('/appointments', { method: 'POST', body: JSON.stringify(data) }),

  getMyAppointments: () => request('/appointments'),

  getTodayAppointments: () => request('/appointments/today'),

  cancel: (id: string) =>
    request(`/appointments/${id}/cancel`, { method: 'POST' }),
}

// --- Inventory ---

export const inventoryService = {
  getAll: (category?: string) => {
    const query = category ? `?category=${category}` : ''
    return request(`/inventory${query}`)
  },

  getById: (id: string) => request(`/inventory/${id}`),
}
