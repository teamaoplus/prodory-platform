import { useState } from 'react'
import { useQuery, useMutation } from '@tanstack/react-query'
import { Lightbulb, CheckCircle, AlertCircle, TrendingDown, Play, Eye } from 'lucide-react'
import { fetchRecommendations, applyRecommendation } from '../services/api'

interface Recommendation {
  id: string
  title: string
  description: string
  category: string
  impact: 'high' | 'medium' | 'low'
  effort: 'high' | 'medium' | 'low'
  savings: number
  status: 'pending' | 'applied' | 'dismissed'
  resources: string[]
}

export default function Recommendations() {
  const [filter, setFilter] = useState<'all' | 'pending' | 'applied'>('all')
  const [selectedRec, setSelectedRec] = useState<Recommendation | null>(null)

  const { data: recommendations, isLoading } = useQuery({
    queryKey: ['recommendations'],
    queryFn: fetchRecommendations,
  })

  const applyMutation = useMutation({
    mutationFn: applyRecommendation,
  })

  const mockRecommendations: Recommendation[] = [
    {
      id: '1',
      title: 'Rightsize EC2 instances (t3.large → t3.medium)',
      description:
        '5 instances are running at average 15% CPU utilization. Downsizing to t3.medium will reduce costs by 40% with minimal performance impact.',
      category: 'Compute',
      impact: 'high',
      effort: 'low',
      savings: 1250,
      status: 'pending',
      resources: ['i-0a1b2c3d4e5f', 'i-1b2c3d4e5f6g', 'i-2c3d4e5f6g7h'],
    },
    {
      id: '2',
      title: 'Purchase Reserved Instances for steady workloads',
      description:
        '3 instances have consistent 24/7 usage patterns. Purchasing 1-year Reserved Instances can save up to 40% compared to On-Demand pricing.',
      category: 'Compute',
      impact: 'high',
      effort: 'medium',
      savings: 2400,
      status: 'pending',
      resources: ['prod-web-01', 'prod-db-01', 'prod-api-01'],
    },
    {
      id: '3',
      title: 'Delete unattached EBS volumes',
      description:
        '12 EBS volumes (450 GB total) are not attached to any EC2 instance. These have been unattached for over 30 days.',
      category: 'Storage',
      impact: 'medium',
      effort: 'low',
      savings: 180,
      status: 'pending',
      resources: ['vol-123', 'vol-456', 'vol-789'],
    },
    {
      id: '4',
      title: 'Enable S3 Intelligent-Tiering',
      description:
        'Buckets prod-backups and prod-logs contain infrequently accessed data. Enabling Intelligent-Tiering can reduce storage costs by 40%.',
      category: 'Storage',
      impact: 'medium',
      effort: 'low',
      savings: 320,
      status: 'pending',
      resources: ['prod-backups', 'prod-logs'],
    },
    {
      id: '5',
      title: 'Consolidate idle RDS instances',
      description:
        '2 RDS instances have zero connections in the last 30 days. Consider deleting or archiving these databases.',
      category: 'Database',
      impact: 'low',
      effort: 'high',
      savings: 450,
      status: 'pending',
      resources: ['old-analytics-db', 'test-replica'],
    },
  ]

  const recs = recommendations || mockRecommendations

  const filteredRecs = recs.filter((rec) => {
    if (filter === 'all') return true
    return rec.status === filter
  })

  const totalSavings = recs
    .filter((r) => r.status === 'pending')
    .reduce((acc, r) => acc + r.savings, 0)

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
          <h1 className="text-2xl font-bold text-gray-900">AI Recommendations</h1>
          <p className="text-gray-500 mt-1">Smart suggestions to optimize your cloud costs</p>
        </div>
        <div className="card bg-green-50 border-green-200">
          <div className="flex items-center gap-3">
            <TrendingDown className="w-6 h-6 text-success" />
            <div>
              <p className="text-sm text-gray-600">Potential Monthly Savings</p>
              <p className="text-2xl font-bold text-success">${totalSavings.toLocaleString()}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Filters */}
      <div className="flex items-center gap-4">
        {(['all', 'pending', 'applied'] as const).map((f) => (
          <button
            key={f}
            onClick={() => setFilter(f)}
            className={`px-4 py-2 rounded-lg text-sm font-medium capitalize transition-colors ${
              filter === f
                ? 'bg-primary-500 text-white'
                : 'bg-white text-gray-600 hover:bg-gray-100 border border-gray-200'
            }`}
          >
            {f} ({f === 'all' ? recs.length : recs.filter((r) => r.status === f).length})
          </button>
        ))}
      </div>

      {/* Recommendations Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {filteredRecs.map((rec) => (
          <div
            key={rec.id}
            className={`card hover:shadow-md transition-shadow cursor-pointer ${
              selectedRec?.id === rec.id ? 'ring-2 ring-primary-500' : ''
            }`}
            onClick={() => setSelectedRec(rec)}
          >
            <div className="flex items-start justify-between">
              <div className="flex items-start gap-3">
                <div
                  className={`w-10 h-10 rounded-lg flex items-center justify-center ${
                    rec.impact === 'high'
                      ? 'bg-red-100'
                      : rec.impact === 'medium'
                      ? 'bg-yellow-100'
                      : 'bg-green-100'
                  }`}
                >
                  <Lightbulb
                    className={`w-5 h-5 ${
                      rec.impact === 'high'
                        ? 'text-danger'
                        : rec.impact === 'medium'
                        ? 'text-warning'
                        : 'text-success'
                    }`}
                  />
                </div>
                <div>
                  <h3 className="font-semibold text-gray-900">{rec.title}</h3>
                  <p className="text-sm text-gray-500 mt-1 line-clamp-2">{rec.description}</p>
                </div>
              </div>
              <div className="text-right">
                <p className="text-lg font-bold text-success">${rec.savings}/mo</p>
                <span
                  className={`inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium ${
                    rec.impact === 'high'
                      ? 'bg-red-100 text-danger'
                      : rec.impact === 'medium'
                      ? 'bg-yellow-100 text-warning'
                      : 'bg-green-100 text-success'
                  }`}
                >
                  <AlertCircle className="w-3 h-3" />
                  {rec.impact} impact
                </span>
              </div>
            </div>

            <div className="flex items-center justify-between mt-4 pt-4 border-t border-gray-100">
              <div className="flex items-center gap-4 text-sm text-gray-500">
                <span>Category: {rec.category}</span>
                <span>Effort: {rec.effort}</span>
                <span>{rec.resources.length} resources</span>
              </div>
              <div className="flex items-center gap-2">
                <button
                  onClick={(e) => {
                    e.stopPropagation()
                    setSelectedRec(rec)
                  }}
                  className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
                >
                  <Eye className="w-4 h-4" />
                </button>
                {rec.status === 'pending' && (
                  <button
                    onClick={(e) => {
                      e.stopPropagation()
                      applyMutation.mutate(rec.id)
                    }}
                    className="btn-primary flex items-center gap-2 text-sm"
                  >
                    <Play className="w-4 h-4" />
                    Apply
                  </button>
                )}
                {rec.status === 'applied' && (
                  <span className="flex items-center gap-1 text-success text-sm">
                    <CheckCircle className="w-4 h-4" />
                    Applied
                  </span>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Detail Modal */}
      {selectedRec && (
        <div
          className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
          onClick={() => setSelectedRec(null)}
        >
          <div
            className="bg-white rounded-xl shadow-xl max-w-2xl w-full mx-4 max-h-[90vh] overflow-auto"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="p-6 border-b border-gray-200">
              <div className="flex items-start justify-between">
                <div>
                  <h2 className="text-xl font-bold text-gray-900">{selectedRec.title}</h2>
                  <p className="text-gray-500 mt-1">{selectedRec.category}</p>
                </div>
                <button
                  onClick={() => setSelectedRec(null)}
                  className="text-gray-400 hover:text-gray-600"
                >
                  ✕
                </button>
              </div>
            </div>

            <div className="p-6 space-y-6">
              <div>
                <h3 className="text-sm font-medium text-gray-700 mb-2">Description</h3>
                <p className="text-gray-600">{selectedRec.description}</p>
              </div>

              <div className="grid grid-cols-3 gap-4">
                <div className="bg-gray-50 rounded-lg p-4">
                  <p className="text-sm text-gray-500">Monthly Savings</p>
                  <p className="text-2xl font-bold text-success">${selectedRec.savings}</p>
                </div>
                <div className="bg-gray-50 rounded-lg p-4">
                  <p className="text-sm text-gray-500">Impact</p>
                  <p className="text-lg font-semibold text-gray-900 capitalize">{selectedRec.impact}</p>
                </div>
                <div className="bg-gray-50 rounded-lg p-4">
                  <p className="text-sm text-gray-500">Effort</p>
                  <p className="text-lg font-semibold text-gray-900 capitalize">{selectedRec.effort}</p>
                </div>
              </div>

              <div>
                <h3 className="text-sm font-medium text-gray-700 mb-2">Affected Resources</h3>
                <div className="flex flex-wrap gap-2">
                  {selectedRec.resources.map((resource, idx) => (
                    <span
                      key={idx}
                      className="px-3 py-1 bg-gray-100 text-gray-700 rounded-full text-sm"
                    >
                      {resource}
                    </span>
                  ))}
                </div>
              </div>

              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <h3 className="text-sm font-medium text-blue-900 mb-2">AI Analysis</h3>
                <p className="text-sm text-blue-700">
                  Based on 30 days of usage data, this recommendation could reduce your monthly
                  spend by ${selectedRec.savings} ({((selectedRec.savings / 15000) * 100).toFixed(1)}%
                  of total compute costs). The change has a low risk profile and can be rolled back
                  within 24 hours if needed.
                </p>
              </div>
            </div>

            <div className="p-6 border-t border-gray-200 flex justify-end gap-3">
              <button
                onClick={() => setSelectedRec(null)}
                className="btn-secondary"
              >
                Dismiss
              </button>
              {selectedRec.status === 'pending' && (
                <button
                  onClick={() => {
                    applyMutation.mutate(selectedRec.id)
                    setSelectedRec(null)
                  }}
                  className="btn-primary flex items-center gap-2"
                >
                  <Play className="w-4 h-4" />
                  Apply Recommendation
                </button>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
