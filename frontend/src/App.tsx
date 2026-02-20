import { useState, useEffect } from 'react'
import { User, Server } from './types'
import { api } from './api/client'
import Sidebar from './components/Sidebar'
import ServerView from './components/ServerView'
import LoginModal from './components/LoginModal'
import CreateServerModal from './components/CreateServerModal'

export default function App() {
  const [currentUser, setCurrentUser] = useState<User | null>(null)
  const [serverIds, setServerIds] = useState<string[]>([])
  const [currentServerId, setCurrentServerId] = useState<string | null>(null)
  const [view, setView] = useState<'posts' | 'messages'>('posts')
  const [showCreateServer, setShowCreateServer] = useState(false)
  const [loadingUser, setLoadingUser] = useState(true)

  useEffect(() => {
    const storedUserId = localStorage.getItem('dischord_user_id')
    const storedServerIds: string[] = JSON.parse(localStorage.getItem('dischord_server_ids') ?? '[]')
    setServerIds(storedServerIds)

    if (storedUserId) {
      api.getUser(storedUserId)
        .then(user => setCurrentUser(user))
        .catch(() => localStorage.removeItem('dischord_user_id'))
        .finally(() => setLoadingUser(false))
    } else {
      setLoadingUser(false)
    }
  }, [])

  const handleLogin = (user: User) => {
    setCurrentUser(user)
    localStorage.setItem('dischord_user_id', user.user_id)
  }

  const handleServerCreated = (server: Server) => {
    const newIds = [...serverIds, server.server_id]
    setServerIds(newIds)
    localStorage.setItem('dischord_server_ids', JSON.stringify(newIds))
    setCurrentServerId(server.server_id)
    setView('posts')
    setShowCreateServer(false)
  }

  const handleJoinServer = (serverId: string) => {
    if (!serverIds.includes(serverId)) {
      const newIds = [...serverIds, serverId]
      setServerIds(newIds)
      localStorage.setItem('dischord_server_ids', JSON.stringify(newIds))
    }
    setCurrentServerId(serverId)
    setView('posts')
  }

  const handleSelectServer = (id: string) => {
    setCurrentServerId(id)
    setView('posts')
  }

  if (loadingUser) {
    return (
      <div className="flex h-screen items-center justify-center bg-[#313338]">
        <div className="text-[#949ba4] animate-pulse text-lg">Loading DisChord...</div>
      </div>
    )
  }

  if (!currentUser) {
    return <LoginModal onLogin={handleLogin} />
  }

  return (
    <div className="flex h-screen bg-[#313338] text-white overflow-hidden">
      <Sidebar
        currentUser={currentUser}
        serverIds={serverIds}
        currentServerId={currentServerId}
        onSelectServer={handleSelectServer}
        onCreateServer={() => setShowCreateServer(true)}
        onJoinServer={handleJoinServer}
      />

      <main className="flex-1 flex flex-col overflow-hidden min-w-0">
        {currentServerId ? (
          <ServerView
            serverId={currentServerId}
            currentUser={currentUser}
            view={view}
            onSetView={setView}
          />
        ) : (
          <div className="flex-1 flex flex-col items-center justify-center text-center gap-4">
            <div className="text-6xl">🎵</div>
            <div>
              <h2 className="text-xl font-bold text-white mb-1">Welcome to DisChord</h2>
              <p className="text-[#949ba4] text-sm">Select a server from the sidebar, or create one to get started.</p>
            </div>
            <button
              onClick={() => setShowCreateServer(true)}
              className="bg-[#5865f2] text-white px-5 py-2.5 rounded-md font-medium hover:bg-[#4752c4] transition-colors"
            >
              Create a Server
            </button>
          </div>
        )}
      </main>

      {showCreateServer && (
        <CreateServerModal
          currentUser={currentUser}
          onCreated={handleServerCreated}
          onClose={() => setShowCreateServer(false)}
        />
      )}
    </div>
  )
}
