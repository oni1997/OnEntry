import { useState, useEffect } from 'react'
import { PasswordEntry } from '../api/client'

export default function Vault() {
  const [entries, setEntries] = useState<PasswordEntry[]>([])
  const [showForm, setShowForm] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [masterKey, setMasterKey] = useState('')

  useEffect(() => {
    const saved = localStorage.getItem('vault_entries')
    if (saved) {
      setEntries(JSON.parse(saved))
    }
  }, [])

  const saveEntries = (newEntries: PasswordEntry[]) => {
    setEntries(newEntries)
    localStorage.setItem('vault_entries', JSON.stringify(newEntries))
  }

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const form = e.target as HTMLFormElement
    const formData = new FormData(form)

    const entry: PasswordEntry = {
      id: editingId || crypto.randomUUID(),
      title: formData.get('title') as string,
      username: formData.get('username') as string,
      password: formData.get('password') as string,
      website: formData.get('website') as string,
      notes: formData.get('notes') as string,
      folder: formData.get('folder') as string || 'General',
      favorite: formData.get('favorite') === 'on',
      tags: [],
      created_at: editingId ? entries.find(e => e.id === editingId)?.created_at || new Date().toISOString() : new Date().toISOString(),
      updated_at: new Date().toISOString(),
    }

    if (editingId) {
      saveEntries(entries.map(e => e.id === editingId ? entry : e))
    } else {
      saveEntries([entry, ...entries])
    }

    setShowForm(false)
    setEditingId(null)
    form.reset()
  }

  const handleEdit = (entry: PasswordEntry) => {
    setEditingId(entry.id)
    setShowForm(true)
  }

  const handleDelete = (id: string) => {
    if (confirm('Are you sure you want to delete this entry?')) {
      saveEntries(entries.filter(e => e.id !== id))
    }
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold text-white">Passwords</h1>
        <button
          onClick={() => { setShowForm(true); setEditingId(null) }}
          className="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded"
        >
          Add Password
        </button>
      </div>

      {showForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <form onSubmit={handleSubmit} className="bg-gray-800 p-6 rounded-lg w-full max-w-md">
            <h2 className="text-xl font-semibold text-white mb-4">
              {editingId ? 'Edit Password' : 'New Password'}
            </h2>
            <div className="space-y-4">
              <input name="title" placeholder="Title" required defaultValue={editingId ? entries.find(e => e.id === editingId)?.title : ''} className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-white" />
              <input name="username" placeholder="Username" defaultValue={editingId ? entries.find(e => e.id === editingId)?.username : ''} className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-white" />
              <input name="password" type="password" placeholder="Password" required className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-white" />
              <input name="website" placeholder="Website" defaultValue={editingId ? entries.find(e => e.id === editingId)?.website : ''} className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-white" />
              <textarea name="notes" placeholder="Notes" defaultValue={editingId ? entries.find(e => e.id === editingId)?.notes : ''} className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-white" />
              <input name="folder" placeholder="Folder" defaultValue={editingId ? entries.find(e => e.id === editingId)?.folder : 'General'} className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-white" />
              <label className="flex items-center text-gray-300">
                <input type="checkbox" name="favorite" className="mr-2" defaultChecked={editingId ? entries.find(e => e.id === editingId)?.favorite : false} />
                Favorite
              </label>
            </div>
            <div className="mt-6 flex justify-end space-x-3">
              <button type="button" onClick={() => { setShowForm(false); setEditingId(null) }} className="px-4 py-2 text-gray-300 hover:text-white">Cancel</button>
              <button type="submit" className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded">{editingId ? 'Update' : 'Create'}</button>
            </div>
          </form>
        </div>
      )}

      <div className="bg-gray-800 rounded-lg overflow-hidden">
        <table className="min-w-full">
          <thead className="bg-gray-700">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase">Title</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase">Username</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase">Website</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase">Folder</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-700">
            {entries.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-6 py-8 text-center text-gray-400">
                  No passwords yet. Click "Add Password" to create one.
                </td>
              </tr>
            ) : (
              entries.map((entry) => (
                <tr key={entry.id} className="hover:bg-gray-750">
                  <td className="px-6 py-4 whitespace-nowrap text-white font-medium">{entry.title}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-gray-300">{entry.username}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-gray-300">{entry.website}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-gray-300">{entry.folder}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    <button onClick={() => handleEdit(entry)} className="text-blue-400 hover:text-blue-300 mr-3">Edit</button>
                    <button onClick={() => handleDelete(entry.id)} className="text-red-400 hover:text-red-300">Delete</button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}