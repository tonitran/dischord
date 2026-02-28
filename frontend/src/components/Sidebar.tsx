import { useState, useEffect } from 'react'
import { User, Server } from '../types'
import { api } from '../api/client'
import Avatar from './Avatar'

interface Props {
  currentUser: User
  serverIds: string[]
  currentServerId: string | null
  onSelectServer: (id: string) => void
  onCreateServer: () => void
  onJoinServer: (id: string) => void
}

export default function Sidebar({
  currentUser,
  serverIds,
  currentServerId,
  onSelectServer,
  onCreateServer,
  onJoinServer,
}: Props) {
  const [friends, setFriends] = useState<User[]>([])
  const [servers, setServers] = useState<Server[]>([])
  const [showAddFriend, setShowAddFriend] = useState(false)
  const [friendInput, setFriendInput] = useState('')
  const [friendError, setFriendError] = useState('')
  const [showJoinServer, setShowJoinServer] = useState(false)
  const [serverInput, setServerInput] = useState('')
  const [serverError, setServerError] = useState('')

  useEffect(() => {
    api.getFriends(currentUser.user_id).then(data => setFriends(data ?? [])).catch(() => {})
  }, [currentUser.user_id])

  useEffect(() => {
    if (serverIds.length === 0) { setServers([]); return }
    Promise.all(serverIds.map(id => api.getServer(id).catch(() => null)))
      .then(results => setServers(results.filter(Boolean) as Server[]))
  }, [serverIds])

  const handleAddFriend = async (e: React.FormEvent) => {
    e.preventDefault()
    setFriendError('')
    try {
      await api.addFriend(currentUser.user_id, friendInput.trim())
      const updated = await api.getFriends(currentUser.user_id)
      setFriends(updated)
      setFriendInput('')
      setShowAddFriend(false)
    } catch (err: unknown) {
      setFriendError(err instanceof Error ? err.message : 'Failed to add friend')
    }
  }

  const handleJoinServer = async (e: React.FormEvent) => {
    e.preventDefault()
    setServerError('')
    try {
      await api.getServer(serverInput.trim())
      onJoinServer(serverInput.trim())
      setServerInput('')
      setShowJoinServer(false)
    } catch {
      setServerError('Server not found')
    }
  }

  const handleCopyId = () => {
    navigator.clipboard.writeText(currentUser.user_id)
  }

  return (
    <aside className="w-60 bg-[#2b2d31] flex flex-col h-screen flex-shrink-0">

      {/* ── Servers ── */}
      <div className="flex-1 overflow-y-auto min-h-0">
        <div className="px-3 pb-1 flex items-center justify-between">
          <span className="text-xs font-semibold text-[#949ba4] uppercase tracking-wide">
            Servers
          </span>
          <div className="flex gap-1">
            <button
              onClick={() => { setShowJoinServer(v => !v); setServerError('') }}
              className="w-5 h-5 flex items-center justify-center text-[#949ba4] hover:text-white rounded transition-colors"
              title="Join server by ID"
            >
              <svg viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4">
                <path d="M11 17h2v-4h4v-2h-4V7h-2v4H7v2h4zm1 5C6.48 22 2 17.52 2 12S6.48 2 12 2s10 4.48 10 10-4.48 10-10 10z" />
              </svg>
            </button>
            <button
              onClick={onCreateServer}
              className="w-5 h-5 flex items-center justify-center text-[#949ba4] hover:text-white rounded transition-colors"
              title="Create server"
            >
              <svg viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4">
                <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z" />
              </svg>
            </button>
          </div>
        </div>

        {showJoinServer && (
          <form onSubmit={handleJoinServer} className="px-3 pb-2">
            <input
              type="text"
              value={serverInput}
              onChange={e => setServerInput(e.target.value)}
              placeholder="Paste server ID"
              className="w-full bg-[#1e1f22] text-white text-xs rounded px-2 py-1.5 focus:outline-none focus:ring-1 focus:ring-[#5865f2] placeholder-[#6d6f78]"
              autoFocus
            />
            {serverError && <p className="text-[#f23f43] text-xs mt-1">{serverError}</p>}
            <button
              type="submit"
              className="mt-1.5 w-full bg-[#5865f2] text-white text-xs rounded py-1 hover:bg-[#4752c4] transition-colors"
            >
              Join Server
            </button>
          </form>
        )}

        <div className="px-2 space-y-0.5 pb-2">
          {servers.map(server => (
            <button
              key={server.server_id}
              onClick={() => onSelectServer(server.server_id)}
              className={`w-full flex items-center gap-2.5 px-2 py-2 rounded-md text-left transition-colors ${
                currentServerId === server.server_id
                  ? 'bg-[#404249] text-white'
                  : 'text-[#949ba4] hover:bg-[#35373c] hover:text-white'
              }`}
            >
              <div className="w-8 h-8 bg-[#1e1f22] rounded-xl flex items-center justify-center text-sm font-bold flex-shrink-0 text-white">
                {server.name[0].toUpperCase()}
              </div>
              <span className="text-sm truncate font-medium">{server.name}</span>
            </button>
          ))}
          {servers.length === 0 && (
            <p className="text-xs text-[#6d6f78] px-2 py-1">No servers yet</p>
          )}
        </div>
      </div>

      {/* ── Friends ── */}
      <div className="flex-shrink-0 border-t border-[#1e1f22]">
        <div className="px-3 pt-2 pb-1 flex items-center justify-between">
          <span className="text-xs font-semibold text-[#949ba4] uppercase tracking-wide">
            Friends
          </span>
          <button
            onClick={() => { setShowAddFriend(v => !v); setFriendError('') }}
            className="w-5 h-5 flex items-center justify-center text-[#949ba4] hover:text-white rounded transition-colors"
            title="Add friend by ID"
          >
            <svg viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4">
              <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z" />
            </svg>
          </button>
        </div>

        {showAddFriend && (
          <form onSubmit={handleAddFriend} className="px-3 pb-2">
            <input
              type="text"
              value={friendInput}
              onChange={e => setFriendInput(e.target.value)}
              placeholder="Paste friend ID"
              className="w-full bg-[#1e1f22] text-white text-xs rounded px-2 py-1.5 focus:outline-none focus:ring-1 focus:ring-[#5865f2] placeholder-[#6d6f78]"
              autoFocus
            />
            {friendError && <p className="text-[#f23f43] text-xs mt-1">{friendError}</p>}
            <button
              type="submit"
              className="mt-1.5 w-full bg-[#5865f2] text-white text-xs rounded py-1 hover:bg-[#4752c4] transition-colors"
            >
              Add Friend
            </button>
          </form>
        )}

        <div className="px-2 pb-2 space-y-0.5 max-h-36 overflow-y-auto">
          {friends.length === 0 ? (
            <p className="text-xs text-[#6d6f78] px-2 py-1">No friends yet</p>
          ) : (
            friends.map(friend => (
              <div
                key={friend.user_id}
                className="flex items-center gap-2.5 px-2 py-1.5 rounded-md text-[#949ba4]"
              >
                <Avatar name={friend.username} size="sm" />
                <span className="text-sm truncate font-medium text-white">{friend.username}</span>
              </div>
            ))
          )}
        </div>
      </div>

      {/* ── User panel ── */}
      <div className="flex-shrink-0 bg-[#232428] px-2 py-2">
        <div className="flex items-center gap-2 mb-2">
          <div className="relative">
            <Avatar name={currentUser.username} size="sm" />
            <div className="absolute bottom-0 right-0 w-3 h-3 bg-[#23a55a] rounded-full border-2 border-[#232428]" />
          </div>
          <div className="flex-1 min-w-0">
            <div className="text-sm font-semibold text-white truncate">{currentUser.username}</div>
            <div className="text-xs text-[#949ba4] truncate">Online</div>
          </div>
          <button
            onClick={() => { localStorage.removeItem('dischord_user_id'); window.location.reload() }}
            className="p-1.5 text-[#949ba4] hover:text-white hover:bg-[#35373c] rounded-md transition-colors flex-shrink-0"
            title="Log out"
          >
            <svg viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4">
              <path d="M17 7l-1.41 1.41L18.17 11H8v2h10.17l-2.58 2.58L17 17l5-5zM4 5h8V3H4c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h8v-2H4V5z" />
            </svg>
          </button>
        </div>

        {/* Controls row */}
        <div className="flex gap-1">
          <button
            className="flex-1 flex items-center justify-center gap-1 py-1 px-1 text-xs text-[#949ba4] hover:text-white hover:bg-[#35373c] rounded-md transition-colors"
            title="Mic (UI only)"
          >
            <svg viewBox="0 0 24 24" fill="currentColor" className="w-3.5 h-3.5">
              <path d="M12 15c1.66 0 2.99-1.34 2.99-3L15 6c0-1.66-1.34-3-3-3S9 4.34 9 6v6c0 1.66 1.34 3 3 3zm5.3-3c0 3-2.54 5.1-5.3 5.1S6.7 15 6.7 12H5c0 3.41 2.72 6.23 6 6.72V22h2v-3.28c3.28-.48 6-3.3 6-6.72h-1.7z" />
            </svg>
            Mic
          </button>
          <button
            className="flex-1 flex items-center justify-center gap-1 py-1 px-1 text-xs text-[#949ba4] hover:text-white hover:bg-[#35373c] rounded-md transition-colors"
            title="Audio (UI only)"
          >
            <svg viewBox="0 0 24 24" fill="currentColor" className="w-3.5 h-3.5">
              <path d="M3 9v6h4l5 5V4L7 9H3zm13.5 3c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02z" />
            </svg>
            Audio
          </button>
          <button
            onClick={handleCopyId}
            className="flex-1 flex items-center justify-center gap-1 py-1 px-1 text-xs text-[#949ba4] hover:text-white hover:bg-[#35373c] rounded-md transition-colors"
            title="Copy your user ID"
          >
            <svg viewBox="0 0 24 24" fill="currentColor" className="w-3.5 h-3.5">
              <path d="M16 1H4c-1.1 0-2 .9-2 2v14h2V3h12V1zm3 4H8c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h11c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm0 16H8V7h11v14z" />
            </svg>
            Copy ID
          </button>
        </div>
      </div>
    </aside>
  )
}
