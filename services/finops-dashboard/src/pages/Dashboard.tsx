import { useQuery } from '@tanstack/react-query'
import {
  TrendingUp,
  TrendingDown,
  DollarSign,
  Server,
  AlertTriangle,
  CheckCircle,
} from 'lucide-react'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from 'recharts'
import { fetchDashboardData } from '../services/api'

const COLORS = ['#0B7DF0', '#06B6D4', '#10B981', '#F59E0B', '#EF4444']

export default function Dashboard() {
  const { data, isLoading } = useQuery({
    queryKey: ['dashboard'],
    queryFn: fetchDashboardData,
  })

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="loading-spinner"></div>
      </div>
    )
  }

  const metrics = data?.metrics || {
    totalSpend: 45230,
    spendChange: 12.5,
    forecastedSpend: 52100,
    savings: 8500,
    resources: 142,
    alerts: 3,
  }

  const costTrend = data?.costTrend || [
    { date: 'Jan', cost: 38000 },
    { date: 'Feb', cost: 39500 },
    { date: 'Mar', cost: 41200 },
    { date: 'Apr', cost: 39800 },
    { date: 'May', cost: 42500 },
    { date: 'Jun', cost: 45230 },
  ]

  const serviceBreakdown = data?.serviceBreakdown || [
    { name: 'Compute', value: 45 },
    { name: 'Storage', value: 25 },
    { name: 'Network', value: 15 },
    { name: 'Database', value: 10 },
    { name: 'Other', value: 5 },
  ]

  return (
    <div className="space-y-6">
      {/* Page Title */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-gray-500 mt-1">Overview of your cloud spend and optimization opportunities</p>
      </div>

      {/* Metrics Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {/* Total Spend */}
        <div className="metric-card">
          <div className="flex items-center justify-between">
            <div>
              <p className="metric-label">Total Monthly Spend</p>
              <p className="metric-value">${metrics.totalSpend.toLocaleString()}</p>
              <p className={`metric-change ${metrics.spendChange > 0 ? 'negative' : 'positive'}`}>
                {metrics.spendChange > 0 ? (
                  <TrendingUp className="w-4 h-4 inline mr-1" />
                ) : (
                  <TrendingDown className="w-4 h-4 inline mr-1" />
                )}
                {Math.abs(metrics.spendChange)}% vs last month
              </p>
            </div>
            <div className="w-12 h-12 bg-primary-100 rounded-xl flex items-center justify-center">
              <DollarSign className="w-6 h-6 text-primary-600" />
            </div>
          </div>
        </div>

        {/* Forecasted Spend */}
        <div className="metric-card">
          <div className="flex items-center justify-between">
            <div>
              <p className="metric-label">Forecasted (EOY)</p>
              <p className="metric-value">${metrics.forecastedSpend.toLocaleString()}</p>
              <p className="metric-change positive">
                <CheckCircle className="w-4 h-4 inline mr-1" />
                On track with budget
              </p>
            </div>
            <div className="w-12 h-12 bg-accent-100 rounded-xl flex items-center justify-center">
              <TrendingUp className="w-6 h-6 text-accent-600" />
            </div>
          </div>
        </div>

        {/* Potential Savings */}
        <div className="metric-card">
          <div className="flex items-center justify-between">
            <div>
              <p className="metric-label">Potential Savings</p>
              <p className="metric-value text-success">${metrics.savings.toLocaleString()}</p>
              <p className="metric-change positive">
                <CheckCircle className="w-4 h-4 inline mr-1" />
                8 recommendations
              </p>
            </div>
            <div className="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center">
              <DollarSign className="w-6 h-6 text-success" />
            </div>
          </div>
        </div>

        {/* Active Resources */}
        <div className="metric-card">
          <div className="flex items-center justify-between">
            <div>
              <p className="metric-label">Active Resources</p>
              <p className="metric-value">{metrics.resources}</p>
              <p className="metric-change negative">
                <AlertTriangle className="w-4 h-4 inline mr-1" />
                {metrics.alerts} need attention
              </p>
            </div>
            <div className="w-12 h-12 bg-orange-100 rounded-xl flex items-center justify-center">
              <Server className="w-6 h-6 text-warning" />
            </div>
          </div>
        </div>
      </div>

      {/* Charts Row */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Cost Trend Chart */}
        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Cost Trend (6 Months)</h3>
          <div className="h-80">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={costTrend}>
                <CartesianGrid strokeDasharray="3 3" stroke="#E5E7EB" />
                <XAxis dataKey="date" stroke="#6B7280" />
                <YAxis stroke="#6B7280" tickFormatter={(value) => `$${value / 1000}k`} />
                <Tooltip
                  contentStyle={{ backgroundColor: '#1F2937', border: 'none', borderRadius: '8px' }}
                  labelStyle={{ color: '#F3F4F6' }}
                  itemStyle={{ color: '#F3F4F6' }}
                  formatter={(value: number) => [`$${value.toLocaleString()}`, 'Cost']}
                />
                <Line
                  type="monotone"
                  dataKey="cost"
                  stroke="#0B7DF0"
                  strokeWidth={2}
                  dot={{ fill: '#0B7DF0', strokeWidth: 2 }}
                  activeDot={{ r: 6 }}
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Service Breakdown */}
        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Cost by Service</h3>
          <div className="h-80">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={serviceBreakdown}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  outerRadius={100}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {serviceBreakdown.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip
                  contentStyle={{ backgroundColor: '#1F2937', border: 'none', borderRadius: '8px' }}
                  itemStyle={{ color: '#F3F4F6' }}
                  formatter={(value: number) => [`${value}%`, 'Percentage']}
                />
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>

      {/* AI Recommendations Preview */}
      <div className="card">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold text-gray-900">AI-Powered Recommendations</h3>
          <button className="text-primary-600 hover:text-primary-700 text-sm font-medium">
            View all →
          </button>
        </div>
        <div className="space-y-4">
          {[
            {
              title: 'Rightsize underutilized EC2 instances',
              description: '5 instances running at < 20% CPU average. Potential savings: $1,200/mo',
              impact: 'high',
              savings: 1200,
            },
            {
              title: 'Purchase Reserved Instances',
              description: '3 instances have consistent usage. RIs could save 40% on compute costs.',
              impact: 'medium',
              savings: 800,
            },
            {
              title: 'Delete unattached EBS volumes',
              description: '12 volumes (450 GB total) are not attached to any instance.',
              impact: 'low',
              savings: 150,
            },
          ].map((rec, index) => (
            <div
              key={index}
              className="flex items-start gap-4 p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
            >
              <div
                className={`w-2 h-2 rounded-full mt-2 ${
                  rec.impact === 'high'
                    ? 'bg-danger'
                    : rec.impact === 'medium'
                    ? 'bg-warning'
                    : 'bg-success'
                }`}
              />
              <div className="flex-1">
                <h4 className="font-medium text-gray-900">{rec.title}</h4>
                <p className="text-sm text-gray-500 mt-1">{rec.description}</p>
              </div>
              <div className="text-right">
                <p className="font-semibold text-success">${rec.savings}/mo</p>
                <button className="text-sm text-primary-600 hover:text-primary-700 mt-1">
                  Apply
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
