import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'

export default function Dashboard() {
  const [stats, setStats] = useState({ total: 0, favorites: 0, recent: 0, weak: 0, reused: 0 })

  useEffect(() => {
    setStats({ total: 42, favorites: 8, recent: 5, weak: 3, reused: 2 })
  }, [])

  const cards = [
    { title: 'Total Passwords', value: stats.total, to: '/vault', color: 'bg-blue-600' },
    { title: 'Favorites', value: stats.favorites, to: '/vault', color: 'bg-yellow-600' },
    { title: 'Recently Added', value: stats.recent, to: '/vault', color: 'bg-green-600' },
    { title: 'Weak Passwords', value: stats.weak, to: '/search?weak=1', color: 'bg-red-600' },
    { title: 'Reused Passwords', value: stats.reused, to: '/search?reused=1', color: 'bg-orange-600' },
  ]

  return (
    <div>
      <h1 className="text-3xl font-bold text-white mb-8">Dashboard</h1>
      <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-6">
        {cards.map((card) => (
          <Link
            key={card.title}
            to={card.to}
            className="bg-gray-800 rounded-lg p-6 hover:bg-gray-750 transition-colors"
          >
            <div className={`${card.color} w-12 h-12 rounded-lg flex items-center justify-center mb-4`}>
              <span className="text-2xl font-bold">{card.value}</span>
            </div>
            <h3 className="text-sm font-medium text-gray-400">{card.title}</h3>
          </Link>
        ))}
      </div>
      <div className="mt-8 grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="bg-gray-800 rounded-lg p-6">
          <h2 className="text-xl font-semibold text-white mb-4">Quick Actions</h2>
          <div className="space-y-3">
            <Link to="/generator" className="block w-full text-center bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded">
              Generate Password
            </Link>
            <Link to="/vault" className="block w-full text-center bg-gray-700 hover:bg-gray-600 text-white font-medium py-2 px-4 rounded">
              View Vault
            </Link>
          </div>
        </div>
        <div className="bg-gray-800 rounded-lg p-6">
          <h2 className="text-xl font-semibold text-white mb-4">Recent Activity</h2>
          <p className="text-gray-400">No recent activity to show.</p>
        </div>
      </div>
    </div>
  )
}