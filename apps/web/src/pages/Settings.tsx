import { useNavigate } from 'react-router-dom'

export default function Settings() {
  const navigate = useNavigate()

  const handleExport = () => {
    const entries = localStorage.getItem('vault_entries')
    const blob = new Blob([entries || '[]'], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `onentry-backup-${new Date().toISOString().split('T')[0]}.json`
    a.click()
    URL.revokeObjectURL(url)
  }

  const handleImport = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    const reader = new FileReader()
    reader.onload = (event) => {
      try {
        const data = JSON.parse(event.target?.result as string)
        localStorage.setItem('vault_entries', JSON.stringify(data))
        alert('Import successful')
      } catch {
        alert('Invalid file format')
      }
    }
    reader.readAsText(file)
  }

  const handleLogout = () => {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user_id')
    navigate('/login')
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-white mb-8">Settings</h1>
      <div className="max-w-2xl space-y-6">
        <div className="bg-gray-800 rounded-lg p-6">
          <h2 className="text-xl font-semibold text-white mb-4">Data Management</h2>
          <div className="space-y-3">
            <button onClick={handleExport} className="block w-full text-left px-4 py-3 bg-gray-700 hover:bg-gray-600 rounded text-white">
              Export Vault
            </button>
            <label className="block w-full text-left px-4 py-3 bg-gray-700 hover:bg-gray-600 rounded text-white cursor-pointer">
              Import Vault
              <input type="file" accept=".json" onChange={handleImport} className="hidden" />
            </label>
          </div>
        </div>

        <div className="bg-gray-800 rounded-lg p-6">
          <h2 className="text-xl font-semibold text-white mb-4">Account</h2>
          <div className="space-y-3">
            <button onClick={handleLogout} className="block w-full text-left px-4 py-3 bg-red-900/50 hover:bg-red-900 text-red-200 rounded">
              Logout
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}