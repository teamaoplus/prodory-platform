import { useState } from 'react'
import { Cloud, Key, Bell, Users, Shield, Database } from 'lucide-react'

export default function Settings() {
  const [activeTab, setActiveTab] = useState('providers')

  const tabs = [
    { id: 'providers', name: 'Cloud Providers', icon: Cloud },
    { id: 'integrations', name: 'Integrations', icon: Database },
    { id: 'notifications', name: 'Notifications', icon: Bell },
    { id: 'team', name: 'Team Members', icon: Users },
    { id: 'security', name: 'Security', icon: Shield },
    { id: 'api', name: 'API Keys', icon: Key },
  ]

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900">Settings</h1>
        <p className="text-gray-500 mt-1">Manage your account and preferences</p>
      </div>

      <div className="flex gap-6">
        {/* Sidebar */}
        <div className="w-64 flex-shrink-0">
          <nav className="space-y-1">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg text-sm font-medium transition-colors ${
                  activeTab === tab.id
                    ? 'bg-primary-50 text-primary-600'
                    : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                }`}
              >
                <tab.icon className="w-5 h-5" />
                {tab.name}
              </button>
            ))}
          </nav>
        </div>

        {/* Content */}
        <div className="flex-1">
          {activeTab === 'providers' && (
            <div className="space-y-6">
              <div className="card">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Connected Cloud Providers</h3>
                <div className="space-y-4">
                  {/* AWS */}
                  <div className="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
                    <div className="flex items-center gap-4">
                      <div className="w-12 h-12 bg-orange-100 rounded-lg flex items-center justify-center">
                        <span className="text-orange-600 font-bold">AWS</span>
                      </div>
                      <div>
                        <p className="font-medium text-gray-900">Amazon Web Services</p>
                        <p className="text-sm text-gray-500">Connected • Last synced 5 min ago</p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <button className="text-sm text-gray-500 hover:text-gray-700">Sync Now</button>
                      <button className="text-sm text-danger hover:text-red-700">Disconnect</button>
                    </div>
                  </div>

                  {/* Azure */}
                  <div className="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
                    <div className="flex items-center gap-4">
                      <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
                        <span className="text-blue-600 font-bold">AZ</span>
                      </div>
                      <div>
                        <p className="font-medium text-gray-900">Microsoft Azure</p>
                        <p className="text-sm text-gray-500">Not connected</p>
                      </div>
                    </div>
                    <button className="btn-primary text-sm">Connect</button>
                  </div>

                  {/* GCP */}
                  <div className="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
                    <div className="flex items-center gap-4">
                      <div className="w-12 h-12 bg-red-100 rounded-lg flex items-center justify-center">
                        <span className="text-red-600 font-bold">GCP</span>
                      </div>
                      <div>
                        <p className="font-medium text-gray-900">Google Cloud Platform</p>
                        <p className="text-sm text-gray-500">Not connected</p>
                      </div>
                    </div>
                    <button className="btn-primary text-sm">Connect</button>
                  </div>
                </div>
              </div>

              <div className="card">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">AWS Configuration</h3>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      AWS Access Key ID
                    </label>
                    <input
                      type="text"
                      value="AKIA..."
                      readOnly
                      className="input-field bg-gray-50"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Default Region
                    </label>
                    <select className="input-field">
                      <option>us-east-1</option>
                      <option>us-west-2</option>
                      <option>eu-west-1</option>
                      <option>ap-south-1</option>
                    </select>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Connected Accounts
                    </label>
                    <div className="flex flex-wrap gap-2">
                      <span className="px-3 py-1 bg-gray-100 text-gray-700 rounded-full text-sm">
                        123456789012 (Production)
                      </span>
                      <span className="px-3 py-1 bg-gray-100 text-gray-700 rounded-full text-sm">
                        098765432109 (Development)
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'notifications' && (
            <div className="card">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">Notification Preferences</h3>
              <div className="space-y-4">
                {[
                  { name: 'Budget Alerts', description: 'Get notified when budgets reach thresholds', enabled: true },
                  { name: 'Daily Summary', description: 'Receive daily cost summary email', enabled: true },
                  { name: 'Anomaly Detection', description: 'Alert on unusual spending patterns', enabled: true },
                  { name: 'Weekly Report', description: 'Weekly cost analysis report', enabled: false },
                  { name: 'New Recommendations', description: 'Notify when AI finds savings opportunities', enabled: true },
                ].map((item, index) => (
                  <div key={index} className="flex items-center justify-between py-3 border-b border-gray-100 last:border-0">
                    <div>
                      <p className="font-medium text-gray-900">{item.name}</p>
                      <p className="text-sm text-gray-500">{item.description}</p>
                    </div>
                    <label className="relative inline-flex items-center cursor-pointer">
                      <input type="checkbox" defaultChecked={item.enabled} className="sr-only peer" />
                      <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary-500"></div>
                    </label>
                  </div>
                ))}
              </div>

              <div className="mt-6 pt-6 border-t border-gray-200">
                <h4 className="font-medium text-gray-900 mb-3">Notification Channels</h4>
                <div className="space-y-3">
                  <div className="flex items-center gap-3">
                    <input type="checkbox" defaultChecked className="rounded" />
                    <span className="text-sm">Email (admin@prodory.com)</span>
                  </div>
                  <div className="flex items-center gap-3">
                    <input type="checkbox" className="rounded" />
                    <span className="text-sm">Slack (#cloud-costs)</span>
                  </div>
                  <div className="flex items-center gap-3">
                    <input type="checkbox" className="rounded" />
                    <span className="text-sm">Webhook</span>
                  </div>
                </div>
              </div>
            </div>
          )}

          {activeTab === 'team' && (
            <div className="card">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-semibold text-gray-900">Team Members</h3>
                <button className="btn-primary text-sm">Invite Member</button>
              </div>
              <div className="space-y-3">
                {[
                  { name: 'Admin User', email: 'admin@prodory.com', role: 'Owner', status: 'active' },
                  { name: 'John Doe', email: 'john@prodory.com', role: 'Admin', status: 'active' },
                  { name: 'Jane Smith', email: 'jane@prodory.com', role: 'Viewer', status: 'pending' },
                ].map((member, index) => (
                  <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-primary-100 rounded-full flex items-center justify-center">
                        <span className="text-primary-600 font-medium">
                          {member.name.split(' ').map((n) => n[0]).join('')}
                        </span>
                      </div>
                      <div>
                        <p className="font-medium text-gray-900">{member.name}</p>
                        <p className="text-sm text-gray-500">{member.email}</p>
                      </div>
                    </div>
                    <div className="flex items-center gap-3">
                      <span
                        className={`px-2 py-1 rounded-full text-xs font-medium ${
                          member.role === 'Owner'
                            ? 'bg-purple-100 text-purple-700'
                            : member.role === 'Admin'
                            ? 'bg-blue-100 text-blue-700'
                            : 'bg-gray-100 text-gray-700'
                        }`}
                      >
                        {member.role}
                      </span>
                      {member.status === 'pending' && (
                        <span className="text-xs text-warning">Pending</span>
                      )}
                      <button className="text-sm text-gray-500 hover:text-gray-700">Edit</button>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {activeTab === 'api' && (
            <div className="card">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-semibold text-gray-900">API Keys</h3>
                <button className="btn-primary text-sm">Generate New Key</button>
              </div>
              <div className="space-y-4">
                {[
                  { name: 'Production API Key', key: 'pk_prod_...x7y9z', created: 'Jan 15, 2024', lastUsed: '2 hours ago' },
                  { name: 'Development API Key', key: 'pk_dev_...a1b2c', created: 'Jan 10, 2024', lastUsed: '1 day ago' },
                ].map((apiKey, index) => (
                  <div key={index} className="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
                    <div>
                      <p className="font-medium text-gray-900">{apiKey.name}</p>
                      <p className="text-sm text-gray-500 font-mono mt-1">{apiKey.key}</p>
                      <p className="text-xs text-gray-400 mt-1">
                        Created: {apiKey.created} • Last used: {apiKey.lastUsed}
                      </p>
                    </div>
                    <div className="flex items-center gap-2">
                      <button className="text-sm text-gray-500 hover:text-gray-700">Copy</button>
                      <button className="text-sm text-danger hover:text-red-700">Revoke</button>
                    </div>
                  </div>
                ))}
              </div>

              <div className="mt-6 p-4 bg-blue-50 rounded-lg">
                <h4 className="font-medium text-blue-900 mb-2">API Documentation</h4>
                <p className="text-sm text-blue-700 mb-3">
                  Use our REST API to programmatically access cost data and manage your account.
                </p>
                <a href="#" className="text-sm text-primary-600 hover:text-primary-700 font-medium">
                  View API Docs →
                </a>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
