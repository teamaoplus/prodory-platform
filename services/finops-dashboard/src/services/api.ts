import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor for auth
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Dashboard API
export const fetchDashboardData = async () => {
  const response = await api.get('/dashboard')
  return response.data
}

// Cost Analysis API
export const fetchCostAnalysis = async (filters: {
  provider?: string
  service?: string
  startDate?: string
  endDate?: string
}) => {
  const response = await api.get('/costs/analysis', { params: filters })
  return response.data
}

// Recommendations API
export const fetchRecommendations = async () => {
  const response = await api.get('/recommendations')
  return response.data
}

export const applyRecommendation = async (id: string) => {
  const response = await api.post(`/recommendations/${id}/apply`)
  return response.data
}

// Budgets API
export const fetchBudgets = async () => {
  const response = await api.get('/budgets')
  return response.data
}

export const createBudget = async (budget: {
  name: string
  amount: number
  period: string
  alerts: number[]
}) => {
  const response = await api.post('/budgets', budget)
  return response.data
}

// Reports API
export const fetchReports = async () => {
  const response = await api.get('/reports')
  return response.data
}

export const generateReport = async (type: string, params: Record<string, unknown>) => {
  const response = await api.post('/reports', { type, params })
  return response.data
}

export const downloadReport = async (id: string) => {
  const response = await api.get(`/reports/${id}/download`, {
    responseType: 'blob',
  })
  return response.data
}

// Settings API
export const fetchSettings = async () => {
  const response = await api.get('/settings')
  return response.data
}

export const updateSettings = async (settings: Record<string, unknown>) => {
  const response = await api.put('/settings', settings)
  return response.data
}

// Cloud Provider Connections
export const connectCloudProvider = async (provider: string, credentials: Record<string, string>) => {
  const response = await api.post('/providers/connect', { provider, credentials })
  return response.data
}

export const fetchConnectedProviders = async () => {
  const response = await api.get('/providers')
  return response.data
}

export default api
