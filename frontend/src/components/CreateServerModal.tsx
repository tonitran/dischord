import { useState } from 'react'
import { User, Server } from '../types'
import { api } from '../api/client'

interface Props {
  currentUser: User
  onCreated: (server: Server) => void
  onClose: () => void
}

export default function CreateServerModal({ currentUser, onCreated, onClose }: Props) {
  const [name, setName] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!name.trim()) return
    setError('')
    setLoading(true)
    try {
      const server = await api.createServer(name.trim(), currentUser.user_id)
      onCreated(server)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to create server')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div
      className="fixed inset-0 bg-black/70 flex items-center justify-center z-50"
      onClick={e => { if (e.target === e.currentTarget) onClose() }}
    >
      <div className="bg-[#2b2d31] rounded-xl p-6 w-full max-w-md shadow-2xl">
        <h2 className="text-xl font-bold text-white mb-1">Create a Server</h2>
        <p className="text-[#949ba4] text-sm mb-5">
          Your server is where you and your friends post and hang out. Give it a name!
        </p>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-xs font-semibold text-[#b5bac1] uppercase tracking-wide mb-1.5">
              Server Name
            </label>
            <input
              type="text"
              value={name}
              onChange={e => setName(e.target.value)}
              className="w-full bg-[#1e1f22] text-white rounded-md px-3 py-2.5 focus:outline-none focus:ring-2 focus:ring-[#5865f2] placeholder-[#6d6f78]"
              placeholder="My Awesome Server"
              required
              autoFocus
            />
          </div>
          {error && <p className="text-[#f23f43] text-sm">{error}</p>}
          <div className="flex gap-2 justify-end pt-1">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-sm text-[#949ba4] hover:text-white transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading || !name.trim()}
              className="px-4 py-2 text-sm bg-[#5865f2] text-white rounded-md hover:bg-[#4752c4] transition-colors disabled:opacity-50"
            >
              {loading ? 'Creating...' : 'Create Server'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
