import { useState } from 'react'
import { User } from '../types'
import { api } from '../api/client'

interface Props {
  onLogin: (user: User) => void
}

export default function LoginModal({ onLogin }: Props) {
  const [tab, setTab] = useState<'create' | 'existing'>('create')
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [userId, setUserId] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const user = await api.createUser(username.trim(), email.trim())
      onLogin(user)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to create user')
    } finally {
      setLoading(false)
    }
  }

  const handleExisting = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const user = await api.getUser(userId.trim())
      onLogin(user)
    } catch {
      setError('User not found. Check the ID and try again.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex h-screen bg-[#313338] items-center justify-center">
      <div className="bg-[#2b2d31] rounded-xl p-8 w-full max-w-md shadow-2xl">
        <div className="text-center mb-6">
          <div className="text-4xl mb-2">🎵</div>
          <h1 className="text-2xl font-bold text-white">Welcome to DisChord</h1>
          <p className="text-[#949ba4] text-sm mt-1">Chat. Post. Connect.</p>
        </div>

        <div className="flex rounded-lg overflow-hidden mb-6 bg-[#1e1f22]">
          <button
            className={`flex-1 py-2 text-sm font-medium transition-colors ${tab === 'create' ? 'bg-[#5865f2] text-white' : 'text-[#949ba4] hover:text-white'}`}
            onClick={() => { setTab('create'); setError('') }}
          >
            Create Account
          </button>
          <button
            className={`flex-1 py-2 text-sm font-medium transition-colors ${tab === 'existing' ? 'bg-[#5865f2] text-white' : 'text-[#949ba4] hover:text-white'}`}
            onClick={() => { setTab('existing'); setError('') }}
          >
            Log In
          </button>
        </div>

        {tab === 'create' ? (
          <form onSubmit={handleCreate} className="space-y-4">
            <div>
              <label className="block text-xs font-semibold text-[#b5bac1] uppercase tracking-wide mb-1.5">
                Username
              </label>
              <input
                type="text"
                value={username}
                onChange={e => setUsername(e.target.value)}
                className="w-full bg-[#1e1f22] text-white rounded-md px-3 py-2.5 focus:outline-none focus:ring-2 focus:ring-[#5865f2] placeholder-[#6d6f78]"
                placeholder="cooluser123"
                required
                autoFocus
              />
            </div>
            <div>
              <label className="block text-xs font-semibold text-[#b5bac1] uppercase tracking-wide mb-1.5">
                Email
              </label>
              <input
                type="email"
                value={email}
                onChange={e => setEmail(e.target.value)}
                className="w-full bg-[#1e1f22] text-white rounded-md px-3 py-2.5 focus:outline-none focus:ring-2 focus:ring-[#5865f2] placeholder-[#6d6f78]"
                placeholder="user@example.com"
                required
              />
            </div>
            {error && <p className="text-[#f23f43] text-sm">{error}</p>}
            <button
              type="submit"
              disabled={loading}
              className="w-full bg-[#5865f2] text-white rounded-md py-2.5 font-medium hover:bg-[#4752c4] transition-colors disabled:opacity-50 mt-2"
            >
              {loading ? 'Creating...' : 'Create Account'}
            </button>
          </form>
        ) : (
          <form onSubmit={handleExisting} className="space-y-4">
            <div>
              <label className="block text-xs font-semibold text-[#b5bac1] uppercase tracking-wide mb-1.5">
                Your User ID
              </label>
              <input
                type="text"
                value={userId}
                onChange={e => setUserId(e.target.value)}
                className="w-full bg-[#1e1f22] text-white rounded-md px-3 py-2.5 focus:outline-none focus:ring-2 focus:ring-[#5865f2] placeholder-[#6d6f78] font-mono text-sm"
                placeholder="Paste your user ID here"
                required
                autoFocus
              />
              <p className="text-xs text-[#6d6f78] mt-1">
                Find your ID in the bottom-left of the app after logging in.
              </p>
            </div>
            {error && <p className="text-[#f23f43] text-sm">{error}</p>}
            <button
              type="submit"
              disabled={loading}
              className="w-full bg-[#5865f2] text-white rounded-md py-2.5 font-medium hover:bg-[#4752c4] transition-colors disabled:opacity-50"
            >
              {loading ? 'Loading...' : 'Log In'}
            </button>
          </form>
        )}
      </div>
    </div>
  )
}
