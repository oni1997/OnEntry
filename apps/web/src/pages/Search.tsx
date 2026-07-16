import { useState, useEffect } from 'react'
import { PasswordEntry } from '../api/client'

export default function Search() {
  const [query, setQuery] = useState('')
  const [results, setResults] = useState<PasswordEntry[]>([])

  useEffect(() => {
    if (!query.trim()) {
      setResults([])
      return
    }
    const saved = localStorage.getItem('vault_entries')
    const entries: PasswordEntry[] = saved ? JSON.parse(saved) : []
    const filtered = entries.filter(e =>
      e.title.toLowerCase().includes(query.toLowerCase()) ||
      e.username.toLowerCase().includes(query.toLowerCase()) ||
      e.website.toLowerCase().includes(query.toLowerCase()) ||
      e.tags.some(t => t.toLowerCase().includes(query.toLowerCase()))
    )
    setResults(filtered)
  }, [query])

  return (
    <div>
      <h1 className="text-3xl font-bold text-white mb-8">Search</h1>
      <div className="mb-6">
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search by title, username, website, or tags..."
          className="w-full px-4 py-3 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>
      <div className="bg-gray-800 rounded-lg overflow-hidden">
        {results.length === 0 ? (
          <div className="p-6 text-center text-gray-400">
            {query ? 'No results found.' : 'Enter a search query to find passwords.'}
          </div>
        ) : (
          <table className="min-w-full">
            <thead className="bg-gray-700">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase">Title</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase">Username</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase">Website</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase">Tags</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-700">
              {results.map((entry) => (
                <tr key={entry.id} className="hover:bg-gray-750">
                  <td className="px-6 py-4 whitespace-nowrap text-white font-medium">{entry.title}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-gray-300">{entry.username}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-gray-300">{entry.website}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-gray-300">
                    {entry.tags.map(t => (
                      <span key={t} className="inline-block bg-gray-700 text-gray-300 text-xs px-2 py-1 rounded mr-1 mb-1">{t}</span>
                    ))}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}