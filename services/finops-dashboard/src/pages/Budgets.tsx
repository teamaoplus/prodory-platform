import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Plus, Bell, Edit2, Trash2, AlertTriangle, CheckCircle } from 'lucide-react'
import { fetchBudgets } from '../services/api'

interface Budget {
  id: string
  name: string
  amount: number
  spent: number
  period: string
  alerts: number[]
  status: 'active' | 'exceeded' | 'warning'
}

export default function Budgets() {
  const [showCreateModal, setShowCreateModal] = useState(false)
  const { data: budgets, isLoading } = useQuery({
    queryKey: ['budgets'],
    queryFn: fetchBudgets,
  })

  const mockBudgets: Budget[] = [
    {
      id: '1',
      name: 'Production Infrastructure',
      amount: 15000,
      spent: 14200,
      period: 'monthly',
      alerts: [80, 100],
      status: 'warning',
    },
    {
      id: '2',
      name: 'Development Environment',
      amount: 5000,
      spent: 3200,
      period: 'monthly',
      alerts: [80, 100],
      status: 'active',
    },
    {
      id: '3',
      name: 'Data Science Workloads',
      amount: 8000,
      spent: 8500,
      period: 'monthly',
      alerts: [80, 100],
      status: 'exceeded',
    },
    {
      id: '4',
      name: 'Q1 2024 Total Budget',
      amount: 45000,
      spent: 28900,
      period: 'quarterly',
      alerts: [75, 90, 100],
      status: 'active',
    },
  ]

  const displayBudgets = budgets || mockBudgets

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="loading-spinner"></div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Budgets & Alerts</h1>
          <p className="text-gray-500 mt-1">Set spending limits and get notified</p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="btn-primary flex items-center gap-2"
        >
          <Plus className="w-4 h-4" />
          Create Budget
        </button>
      </div>

      {/* Budget Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {displayBudgets.map((budget) => {
          const percentage = (budget.spent / budget.amount) * 100
          return (
            <div key={budget.id} className="card">
              <div className="flex items-start justify-between mb-4">
                <div>
                  <h3 className="font-semibold text-gray-900">{budget.name}</h3>
                  <p className="text-sm text-gray-500 capitalize">{budget.period} budget</p>
                </div>
                <div className="flex items-center gap-2">
                  {budget.status === 'exceeded' && (
                    <span className="flex items-center gap-1 px-2 py-1 bg-red-100 text-danger rounded-full text-xs font-medium">
                      <AlertTriangle className="w-3 h-3" />
                      Exceeded
                    </span>
                  )}
                  {budget.status === 'warning' && (
                    <span className="flex items-center gap-1 px-2 py-1 bg-yellow-100 text-warning rounded-full text-xs font-medium">
                      <AlertTriangle className="w-3 h-3" />
                      Warning
                    </span>
                  )}
                  {budget.status === 'active' && (
                    <span className="flex items-center gap-1 px-2 py-1 bg-green-100 text-success rounded-full text-xs font-medium">
                      <CheckCircle className="w-3 h-3" />
                      On Track
                    </span>
                  )}
                  <button className="p-1 text-gray-400 hover:text-gray-600">
                    <Edit2 className="w-4 h-4" />
                  </button>
                  <button className="p-1 text-gray-400 hover:text-danger">
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>

              <div className="mb-4">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-2xl font-bold text-gray-900">
                    ${budget.spent.toLocaleString()}
                  </span>
                  <span className="text-sm text-gray-500">
                    of ${budget.amount.toLocaleString()}
                  </span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-3">
                  <div
                    className={`h-3 rounded-full transition-all ${
                      percentage > 100
                        ? 'bg-danger'
                        : percentage > 80
                        ? 'bg-warning'
                        : 'bg-success'
                    }`}
                    style={{ width: `${Math.min(percentage, 100)}%` }}
                  />
                </div>
                <p className="text-sm text-gray-500 mt-2">{percentage.toFixed(1)}% used</p>
              </div>

              <div className="flex items-center gap-2 text-sm text-gray-500">
                <Bell className="w-4 h-4" />
                Alerts at {budget.alerts.join('%, ')}%
              </div>
            </div>
          )
        })}
      </div>

      {/* Alert History */}
      <div className="card">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Recent Alerts</h3>
        <div className="space-y-3">
          {[
            {
              message: 'Production Infrastructure budget reached 80%',
              time: '2 hours ago',
              type: 'warning',
            },
            {
              message: 'Data Science Workloads budget exceeded',
              time: '1 day ago',
              type: 'danger',
            },
            {
              message: 'Daily spend spike detected (+45%)',
              time: '2 days ago',
              type: 'info',
            },
          ].map((alert, index) => (
            <div
              key={index}
              className={`flex items-center gap-3 p-3 rounded-lg ${
                alert.type === 'danger'
                  ? 'bg-red-50'
                  : alert.type === 'warning'
                  ? 'bg-yellow-50'
                  : 'bg-blue-50'
              }`}
            >
              <AlertTriangle
                className={`w-5 h-5 ${
                  alert.type === 'danger'
                    ? 'text-danger'
                    : alert.type === 'warning'
                    ? 'text-warning'
                    : 'text-primary-500'
                }`}
              />
              <div className="flex-1">
                <p className="text-sm font-medium text-gray-900">{alert.message}</p>
                <p className="text-xs text-gray-500">{alert.time}</p>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Create Budget Modal */}
      {showCreateModal && (
        <div
          className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
          onClick={() => setShowCreateModal(false)}
        >
          <div
            className="bg-white rounded-xl shadow-xl max-w-lg w-full mx-4"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="p-6 border-b border-gray-200">
              <h2 className="text-xl font-bold text-gray-900">Create New Budget</h2>
            </div>
            <div className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Budget Name
                </label>
                <input type="text" className="input-field" placeholder="e.g., Production Team" />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Budget Amount
                </label>
                <input type="number" className="input-field" placeholder="5000" />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Period</label>
                <select className="input-field">
                  <option value="daily">Daily</option>
                  <option value="weekly">Weekly</option>
                  <option value="monthly">Monthly</option>
                  <option value="quarterly">Quarterly</option>
                  <option value="yearly">Yearly</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Alert Thresholds (%)
                </label>
                <div className="flex gap-2">
                  {[80, 90, 100].map((threshold) => (
                    <label key={threshold} className="flex items-center gap-2">
                      <input type="checkbox" defaultChecked className="rounded" />
                      <span className="text-sm">{threshold}%</span>
                    </label>
                  ))}
                </div>
              </div>
            </div>
            <div className="p-6 border-t border-gray-200 flex justify-end gap-3">
              <button onClick={() => setShowCreateModal(false)} className="btn-secondary">
                Cancel
              </button>
              <button onClick={() => setShowCreateModal(false)} className="btn-primary">
                Create Budget
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
