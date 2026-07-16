import { useState } from 'react'
import { generatePassword } from '../api/client'

export default function Generator() {
  const [length, setLength] = useState(16)
  const [uppercase, setUppercase] = useState(true)
  const [lowercase, setLowercase] = useState(true)
  const [numbers, setNumbers] = useState(true)
  const [symbols, setSymbols] = useState(true)
  const [excludeSimilar, setExcludeSimilar] = useState(false)
  const [pronounceable, setPronounceable] = useState(false)
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [copied, setCopied] = useState(false)

  const handleGenerate = async () => {
    setLoading(true)
    try {
      const data = await generatePassword(length, {
        uppercase, lowercase, numbers, symbols,
        excludeSimilar, pronounceable
      })
      setPassword(data.password)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const copyToClipboard = () => {
    navigator.clipboard.writeText(password)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  const strength = calculateStrength(password)

  return (
    <div>
      <h1 className="text-3xl font-bold text-white mb-8">Password Generator</h1>
      <div className="max-w-2xl">
        <div className="bg-gray-800 rounded-lg p-6 mb-6">
          <div className="flex items-center justify-between mb-4">
            <div className="flex-1">
              <input
                type="text"
                readOnly
                value={password}
                placeholder="Generated password"
                className="w-full bg-gray-700 border border-gray-600 rounded px-4 py-3 text-white font-mono text-lg"
              />
            </div>
            <button
              onClick={copyToClipboard}
              className="ml-4 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded"
            >
              {copied ? 'Copied!' : 'Copy'}
            </button>
          </div>
          {password && (
            <div className="flex items-center space-x-2">
              <div className="flex-1 bg-gray-700 rounded-full h-2">
                <div
                  className={`h-2 rounded-full ${strength.color}`}
                  style={{ width: `${strength.percentage}%` }}
                />
              </div>
              <span className={`text-sm font-medium ${strength.textColor}`}>{strength.label}</span>
            </div>
          )}
        </div>

        <div className="bg-gray-800 rounded-lg p-6 space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Length: {length}
            </label>
            <input
              type="range"
              min="4"
              max="128"
              value={length}
              onChange={(e) => setLength(Number(e.target.value))}
              className="w-full"
            />
          </div>

          <div className="space-y-3">
            {[
              { label: 'Uppercase', value: uppercase, set: setUppercase },
              { label: 'Lowercase', value: lowercase, set: setLowercase },
              { label: 'Numbers', value: numbers, set: setNumbers },
              { label: 'Symbols', value: symbols, set: setSymbols },
              { label: 'Exclude Similar', value: excludeSimilar, set: setExcludeSimilar },
              { label: 'Pronounceable', value: pronounceable, set: setPronounceable },
            ].map(({ label, value, set }) => (
              <label key={label} className="flex items-center text-gray-300">
                <input
                  type="checkbox"
                  checked={value}
                  onChange={(e) => set(e.target.checked)}
                  className="mr-3 h-4 w-4 text-blue-600 rounded border-gray-600 bg-gray-700"
                />
                {label}
              </label>
            ))}
          </div>

          <button
            onClick={handleGenerate}
            disabled={loading}
            className="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 px-4 rounded disabled:opacity-50"
          >
            {loading ? 'Generating...' : 'Generate Password'}
          </button>
        </div>
      </div>
    </div>
  )
}

function calculateStrength(password: string) {
  if (!password) return { percentage: 0, label: 'None', color: 'bg-gray-600', textColor: 'text-gray-400' }

  let score = 0
  if (password.length >= 8) score++
  if (password.length >= 12) score++
  if (/[A-Z]/.test(password)) score++
  if (/[a-z]/.test(password)) score++
  if (/[0-9]/.test(password)) score++
  if (/[^A-Za-z0-9]/.test(password)) score++

  if (score <= 2) return { percentage: 25, label: 'Weak', color: 'bg-red-500', textColor: 'text-red-400' }
  if (score <= 4) return { percentage: 50, label: 'Fair', color: 'bg-yellow-500', textColor: 'text-yellow-400' }
  if (score <= 5) return { percentage: 75, label: 'Strong', color: 'bg-blue-500', textColor: 'text-blue-400' }
  return { percentage: 100, label: 'Very Strong', color: 'bg-green-500', textColor: 'text-green-400' }
}