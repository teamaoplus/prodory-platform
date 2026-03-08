import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  AreaChart,
  Area,
} from 'recharts'
import { Filter, Download, Calendar } from 'lucide-react'
import { fetchCostAnalysis } from '../services/api'

const providers = ['All', 'AWS', 'Azure', 'GCP']
const services = ['All', 'Compute', 'Storage', 'Network', 'Database', 'Other']

export default function CostAnalysis() {
  const [selectedProvider, setSelectedProvider] = useState('All')
  const [selectedService, setSelectedService] = useState('All')
  const [dateRange, setDateRange] = useState('30d')

  const { data, isLoading } = useQuery({
    queryKey: ['costAnalysis', selectedProvider, selectedService, dateRange],
    queryFn: () =>
      fetchCostAnalysis({
        provider: selectedProvider === 'All' ? undefined : selectedProvider,
        service: selectedService === 'All' ? undefined : selectedService,
        period: dateRange,
      }),
  })

  const dailyCosts = data?.dailyCosts || [
    { date: '01', cost: 1200 },
    { date: '02', cost: 1350 },
    { date: '03', cost: 1100 },
    { date: '04', cost: 1450 },
    { date: '05', cost: 1300 },
    { date: '06', cost: 1600 },
    { date: '07', cost: 1500 },
  ]

  const serviceCosts = data?.serviceCosts || [
    { service: 'EC2', cost: 5200, change: 12 },
    { service: 'S3', cost: 1800, change: -5 },
    { service: 'RDS', cost: 2400, change: 8 },
    { service: 'Lambda', cost: 800, change: 25 },
    { service: 'CloudFront', cost: 600, change: -2 },
  ]

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Cost Analysis</h1>
          <p className="text-gray-500 mt-1">Deep dive into your cloud spending patterns</p>
        </div>
        <button className="btn-secondary flex items-center gap-2">
          <Download className="w-4 h-4" />
          Export Report
        </button>
      </div>

      {/* Filters */}
      <div className="card flex flex-wrap items-center gap-4">
        <div className="flex items-center gap-2">
          <Filter className="w-5 h-5 text-gray-400" />
          <span className="text-sm font-medium text-gray-700">Filters:</span>
        </div>

        <select
          value={selectedProvider}
          onChange={(e) => setSelectedProvider(e.target.value)}
          className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
        >
          {providers.map((p) => (
            <option key={p} value={p}>
              {p} Provider
            </option>
          ))}
        </select>

        <select
          value={selectedService}
          onChange={(e) => setSelectedService(e.target.value)}
          className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
        >
          {services.map((s) => (
            <option key={s} value={s}>
              {s} Service
            </option>
          ))}
        </select>

        <div className="flex items-center gap-2">
          <Calendar className="w-4 h-4 text-gray-400" />
          <select
            value={dateRange}
            onChange={(e) => setDateRange(e.target.value)}
            className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
          >
            <option value="7d">Last 7 days</option>
            <option value="30d">Last 30 days</option>
            <option value="90d">Last 90 days</option>
            <option value="1y">Last year</option>
          </select>
        </div>
      </div>

      {/* Daily Cost Chart */}
      <div className="card">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Daily Cost Breakdown</h3>
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={dailyCosts}>
              <defs>
                <linearGradient id="colorCost" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#0B7DF0" stopOpacity={0.3} />
                  <stop offset="95%" stopColor="#0B7DF0" stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="#E5E7EB" />
              <XAxis dataKey="date" stroke="#6B7280" />
              <YAxis stroke="#6B7280" tickFormatter={(value) => `$${value}`} />
              <Tooltip
                contentStyle={{ backgroundColor: '#1F2937', border: 'none', borderRadius: '8px' }}
                labelStyle={{ color: '#F3F4F6' }}
                itemStyle={{ color: '#F3F4F6' }}
                formatter={(value: number) => [`$${value}`, 'Cost']}
              />
              <Area
                type="monotone"
                dataKey="cost"
                stroke="#0B7DF0"
                fillOpacity={1}
                fill="url(#colorCost)"
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Service Cost Table */}
      <div className="card">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Cost by Service</h3>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200">
                <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">Service</th>
                <th className="text-right py-3 px-4 text-sm font-medium text-gray-500">Monthly Cost</th>
                <th className="text-right py-3 px-4 text-sm font-medium text-gray-500">% of Total</th>
                <th className="text-right py-3 px-4 text-sm font-medium text-gray-500">Change</th>
                <th className="text-right py-3 px-4 text-sm font-medium text-gray-500">Trend</th>
              </tr>
            </thead>
            <tbody>
              {serviceCosts.map((service, index) => (
                <tr key={index} className="border-b border-gray-100 hover:bg-gray-50">
                  <td className="py-4 px-4 font-medium text-gray-900">{service.service}</td>
                  <td className="py-4 px-4 text-right">${service.cost.toLocaleString()}</td>
                  <td className="py-4 px-4 text-right">
                    {((service.cost / serviceCosts.reduce((a, b) => a + b.cost, 0)) * 100).toFixed(1)}%
                  </td>
                  <td className="py-4 px-4 text-right">
                    <span
                      className={`inline-flex items-center ${
                        service.change > 0 ? 'text-danger' : 'text-success'
                      }`}
                    >
                      {service.change > 0 ? '↑' : '↓'} {Math.abs(service.change)}%
                    </span>
                  </td>
                  <td className="py-4 px-4 text-right">
                    <div className="w-24 h-8 ml-auto">
                      <ResponsiveContainer width="100%" height="100%">
                        <BarChart data={[{ v: service.cost * 0.8 }, { v: service.cost * 0.9 }, { v: service.cost }]}>
                          <Bar dataKey="v" fill={service.change > 0 ? '#EF4444' : '#10B981'} />
                        </BarChart>
                      </ResponsiveContainer>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Cost Allocation */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Cost by Team</h3>
          <div className="space-y-4">
            {[
              { team: 'Engineering', cost: 8500, budget: 10000 },
              { team: 'Data Science', cost: 4200, budget: 5000 },
              { team: 'DevOps', cost: 3100, budget: 4000 },
              { team: 'QA', cost: 1800, budget: 2500 },
            ].map((item, index) => (
              <div key={index}>
                <div className="flex items-center justify-between mb-1">
                  <span className="text-sm font-medium text-gray-700">{item.team}</span>
                  <span className="text-sm text-gray-500">
                    ${item.cost.toLocaleString()} / ${item.budget.toLocaleString()}
                  </span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className={`h-2 rounded-full ${
                      item.cost / item.budget > 0.9 ? 'bg-danger' : 'bg-primary-500'
                    }`}
                    style={{ width: `${(item.cost / item.budget) * 100}%` }}
                  />
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Cost by Environment</h3>
          <div className="space-y-4">
            {[
              { env: 'Production', cost: 15200, percentage: 65 },
              { env: 'Staging', cost: 4800, percentage: 20 },
              { env: 'Development', cost: 2400, percentage: 10 },
              { env: 'Testing', cost: 1200, percentage: 5 },
            ].map((item, index) => (
              <div key={index} className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div
                    className={`w-3 h-3 rounded-full ${
                      item.env === 'Production'
                        ? 'bg-danger'
                        : item.env === 'Staging'
                        ? 'bg-warning'
                        : item.env === 'Development'
                        ? 'bg-success'
                        : 'bg-gray-400'
                    }`}
                  />
                  <span className="text-sm font-medium text-gray-700">{item.env}</span>
                </div>
                <div className="text-right">
                  <span className="text-sm font-medium text-gray-900">
                    ${item.cost.toLocaleString()}
                  </span>
                  <span className="text-sm text-gray-500 ml-2">({item.percentage}%)</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
