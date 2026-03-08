import { useState } from 'react'
import { Download, FileText, Calendar, Filter, ChevronDown } from 'lucide-react'

interface Report {
  id: string
  name: string
  type: string
  createdAt: string
  size: string
  status: 'ready' | 'generating' | 'scheduled'
}

export default function Reports() {
  const [selectedType, setSelectedType] = useState('all')

  const reports: Report[] = [
    {
      id: '1',
      name: 'Monthly Cost Summary - January 2024',
      type: 'cost-summary',
      createdAt: '2024-02-01',
      size: '2.4 MB',
      status: 'ready',
    },
    {
      id: '2',
      name: 'AWS Resource Utilization Report',
      type: 'utilization',
      createdAt: '2024-01-28',
      size: '5.1 MB',
      status: 'ready',
    },
    {
      id: '3',
      name: 'Reserved Instance Analysis',
      type: 'ri-analysis',
      createdAt: '2024-01-25',
      size: '1.8 MB',
      status: 'ready',
    },
    {
      id: '4',
      name: 'Weekly Cost Trend Report',
      type: 'trend',
      createdAt: 'Generating...',
      size: '-',
      status: 'generating',
    },
  ]

  const reportTypes = [
    { id: 'cost-summary', name: 'Cost Summary', description: 'Monthly spend overview' },
    { id: 'utilization', name: 'Resource Utilization', description: 'Usage and efficiency metrics' },
    { id: 'ri-analysis', name: 'RI/SP Analysis', description: 'Reserved capacity recommendations' },
    { id: 'trend', name: 'Cost Trends', description: 'Historical spending patterns' },
    { id: 'forecast', name: 'Forecast Report', description: 'Predicted future costs' },
  ]

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Reports</h1>
        <p className="text-gray-500 mt-1">Generate and download cost reports</p>
      </div>

      {/* Generate Report Section */}
      <div className="card">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Generate New Report</h3>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Report Type</label>
            <select className="input-field">
              {reportTypes.map((type) => (
                <option key={type.id} value={type.id}>
                  {type.name}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Date Range</label>
            <select className="input-field">
              <option>Last 7 days</option>
              <option>Last 30 days</option>
              <option>Last 90 days</option>
              <option>Custom range</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Cloud Provider</label>
            <select className="input-field">
              <option>All Providers</option>
              <option>AWS</option>
              <option>Azure</option>
              <option>GCP</option>
            </select>
          </div>
          <div className="flex items-end">
            <button className="btn-primary w-full flex items-center justify-center gap-2">
              <FileText className="w-4 h-4" />
              Generate Report
            </button>
          </div>
        </div>

        {/* Report Type Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6">
          {reportTypes.slice(0, 3).map((type) => (
            <div
              key={type.id}
              className="p-4 border border-gray-200 rounded-lg hover:border-primary-500 hover:bg-primary-50 cursor-pointer transition-colors"
            >
              <FileText className="w-8 h-8 text-primary-500 mb-2" />
              <h4 className="font-medium text-gray-900">{type.name}</h4>
              <p className="text-sm text-gray-500">{type.description}</p>
            </div>
          ))}
        </div>
      </div>

      {/* Scheduled Reports */}
      <div className="card">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold text-gray-900">Scheduled Reports</h3>
          <button className="text-primary-600 hover:text-primary-700 text-sm font-medium">
            + Schedule New
          </button>
        </div>
        <div className="space-y-3">
          {[
            { name: 'Weekly Executive Summary', schedule: 'Every Monday', recipients: '5 people' },
            { name: 'Monthly Cost Report', schedule: '1st of month', recipients: '12 people' },
          ].map((schedule, index) => (
            <div
              key={index}
              className="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div>
                <p className="font-medium text-gray-900">{schedule.name}</p>
                <p className="text-sm text-gray-500">
                  {schedule.schedule} • {schedule.recipients}
                </p>
              </div>
              <div className="flex items-center gap-2">
                <button className="text-sm text-gray-500 hover:text-gray-700">Edit</button>
                <button className="text-sm text-danger hover:text-red-700">Delete</button>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Report History */}
      <div className="card">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold text-gray-900">Report History</h3>
          <div className="flex items-center gap-2">
            <Filter className="w-4 h-4 text-gray-400" />
            <select
              value={selectedType}
              onChange={(e) => setSelectedType(e.target.value)}
              className="text-sm border-none bg-transparent focus:ring-0"
            >
              <option value="all">All Types</option>
              {reportTypes.map((t) => (
                <option key={t.id} value={t.id}>
                  {t.name}
                </option>
              ))}
            </select>
          </div>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200">
                <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">Report Name</th>
                <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">Type</th>
                <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">Created</th>
                <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">Size</th>
                <th className="text-left py-3 px-4 text-sm font-medium text-gray-500">Status</th>
                <th className="text-right py-3 px-4 text-sm font-medium text-gray-500">Action</th>
              </tr>
            </thead>
            <tbody>
              {reports.map((report) => (
                <tr key={report.id} className="border-b border-gray-100 hover:bg-gray-50">
                  <td className="py-4 px-4">
                    <div className="flex items-center gap-3">
                      <FileText className="w-5 h-5 text-gray-400" />
                      <span className="font-medium text-gray-900">{report.name}</span>
                    </div>
                  </td>
                  <td className="py-4 px-4 text-sm text-gray-600 capitalize">
                    {report.type.replace('-', ' ')}
                  </td>
                  <td className="py-4 px-4 text-sm text-gray-600">{report.createdAt}</td>
                  <td className="py-4 px-4 text-sm text-gray-600">{report.size}</td>
                  <td className="py-4 px-4">
                    <span
                      className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                        report.status === 'ready'
                          ? 'bg-green-100 text-success'
                          : report.status === 'generating'
                          ? 'bg-blue-100 text-primary-600'
                          : 'bg-gray-100 text-gray-600'
                      }`}
                    >
                      {report.status === 'generating' && (
                        <div className="w-3 h-3 border-2 border-primary-500 border-t-transparent rounded-full animate-spin mr-1" />
                      )}
                      {report.status.charAt(0).toUpperCase() + report.status.slice(1)}
                    </span>
                  </td>
                  <td className="py-4 px-4 text-right">
                    {report.status === 'ready' && (
                      <button className="text-primary-600 hover:text-primary-700 text-sm font-medium flex items-center gap-1 ml-auto">
                        <Download className="w-4 h-4" />
                        Download
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
