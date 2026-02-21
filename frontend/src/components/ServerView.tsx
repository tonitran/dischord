import { useState, useEffect, useRef } from 'react'
import { User, Server, Post, Message } from '../types'
import { api } from '../api/client'
import PostCard from './PostCard'
import CreatePostModal from './CreatePostModal'

interface Props {
  serverId: string
  currentUser: User
  view: 'posts' | 'messages'
  onSetView: (v: 'posts' | 'messages') => void
}

export default function ServerView({ serverId, currentUser, view, onSetView }: Props) {
  const [server, setServer] = useState<Server | null>(null)
  const [posts, setPosts] = useState<Post[]>([])
  const [messages, setMessages] = useState<Message[]>([])
  const [userCache, setUserCache] = useState<Record<string, User>>({})
  const [loading, setLoading] = useState(true)
  const [showCreatePost, setShowCreatePost] = useState(false)
  const [messageInput, setMessageInput] = useState('')
  const [sendingMsg, setSendingMsg] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const prevServerRef = useRef<string>('')

  useEffect(() => {
    if (prevServerRef.current === serverId) return
    prevServerRef.current = serverId

    let cancelled = false
    setLoading(true)
    setPosts([])
    setMessages([])

    async function load() {
      try {
        const s: Server = await api.getServer(serverId)
        if (cancelled) return
        setServer(s)

        // Load posts
        const postResults = await Promise.all(
          s.post_ids.map((id: string) => api.getPost(serverId, id).catch(() => null))
        )
        const msgs = await api.getMessages(serverId).catch(() => [] as Message[])
        if (cancelled) return

        const validPosts = postResults.filter(Boolean) as Post[]
        setPosts(validPosts)
        setMessages(msgs ?? [])

        // Prefetch all unique authors + members
        const authorIds = new Set<string>([
          ...s.member_ids,
          ...validPosts.map(p => p.author_id),
          ...(msgs ?? []).map((m: Message) => m.author_id),
        ])
        const userEntries = await Promise.all(
          [...authorIds].map(id =>
            api.getUser(id).then((u: User) => [id, u] as [string, User]).catch(() => null)
          )
        )
        if (cancelled) return
        const cache: Record<string, User> = {}
        userEntries.forEach(entry => { if (entry) cache[entry[0]] = entry[1] })
        setUserCache(cache)
      } finally {
        if (!cancelled) setLoading(false)
      }
    }

    load()
    return () => {
      cancelled = true
      prevServerRef.current = ''
    }
  }, [serverId])

  // Cache newly seen authors on the fly
  const ensureUser = async (id: string) => {
    if (userCache[id]) return
    try {
      const u: User = await api.getUser(id)
      setUserCache(prev => ({ ...prev, [id]: u }))
    } catch { /* ignore */ }
  }

  useEffect(() => {
    if (view === 'messages') {
      messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
    }
  }, [messages, view])

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!messageInput.trim() || sendingMsg) return
    setSendingMsg(true)
    try {
      const msg: Message = await api.createMessage(serverId, currentUser.user_id, messageInput.trim())
      setMessages(prev => [...prev, msg])
      setMessageInput('')
      await ensureUser(msg.author_id)
    } finally {
      setSendingMsg(false)
    }
  }

  const handlePostCreated = async (post: Post) => {
    setPosts(prev => [post, ...prev])
    setShowCreatePost(false)
    await ensureUser(post.author_id)
  }

  const handlePostUpdated = (post: Post) => {
    setPosts(prev => prev.map(p => p.post_id === post.post_id ? post : p))
  }

  const handlePostDeleted = (postId: string) => {
    setPosts(prev => prev.filter(p => p.post_id !== postId))
  }

  if (loading) {
    return (
      <div className="flex-1 flex items-center justify-center bg-[#313338]">
        <div className="text-[#949ba4] animate-pulse">Loading...</div>
      </div>
    )
  }

  if (!server) {
    return (
      <div className="flex-1 flex items-center justify-center bg-[#313338]">
        <div className="text-[#f23f43]">Failed to load server.</div>
      </div>
    )
  }

  return (
    <div className="flex-1 flex flex-col bg-[#313338] overflow-hidden">

      {/* ── Header ── */}
      <header className="flex-shrink-0 h-12 bg-[#313338] border-b border-[#1e1f22] flex items-center px-4 gap-3 shadow-sm z-10">
        <div className="flex items-center gap-2 mr-2">
          <span className="text-white font-bold text-lg leading-none">{server.name}</span>
        </div>

        {/* Tab switcher */}
        <div className="flex gap-1">
          <button
            onClick={() => onSetView('posts')}
            className={`px-3 py-1 rounded text-sm font-medium transition-colors ${
              view === 'posts'
                ? 'bg-[#5865f2] text-white'
                : 'text-[#949ba4] hover:text-white hover:bg-[#35373c]'
            }`}
          >
            Posts
          </button>
          <button
            onClick={() => onSetView('messages')}
            className={`px-3 py-1 rounded text-sm font-medium transition-colors ${
              view === 'messages'
                ? 'bg-[#5865f2] text-white'
                : 'text-[#949ba4] hover:text-white hover:bg-[#35373c]'
            }`}
          >
            Chat
          </button>
        </div>

        <div className="ml-auto flex items-center gap-3">
          <span className="text-xs text-[#6d6f78]">
            {server.member_ids.length} member{server.member_ids.length !== 1 ? 's' : ''}
          </span>

          {/* Search button (placeholder) */}
          <button
            className="text-[#949ba4] hover:text-white p-1 hover:bg-[#35373c] rounded transition-colors"
            title="Search posts"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="w-4 h-4">
              <circle cx="11" cy="11" r="8" />
              <line x1="21" y1="21" x2="16.65" y2="16.65" />
            </svg>
          </button>

          {/* Pin button (placeholder) */}
          <button
            className="text-[#949ba4] hover:text-white p-1 hover:bg-[#35373c] rounded transition-colors"
            title="Pinned posts"
          >
            <svg viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4">
              <path d="M16 12V4h1V2H7v2h1v8l-2 2v2h5.2v6h1.6v-6H18v-2l-2-2z" />
            </svg>
          </button>

          {/* Share server ID */}
          <button
            onClick={() => { navigator.clipboard.writeText(server.server_id) }}
            className="text-xs text-[#949ba4] hover:text-white flex items-center gap-1 px-2 py-1 hover:bg-[#35373c] rounded transition-colors"
            title="Copy server ID to share"
          >
            <svg viewBox="0 0 24 24" fill="currentColor" className="w-3.5 h-3.5">
              <path d="M18 16.08c-.76 0-1.44.3-1.96.77L8.91 12.7c.05-.23.09-.46.09-.7s-.04-.47-.09-.7l7.05-4.11c.54.5 1.25.81 2.04.81 1.66 0 3-1.34 3-3s-1.34-3-3-3-3 1.34-3 3c0 .24.04.47.09.7L8.04 9.81C7.5 9.31 6.79 9 6 9c-1.66 0-3 1.34-3 3s1.34 3 3 3c.79 0 1.5-.31 2.04-.81l7.12 4.16c-.05.21-.08.43-.08.65 0 1.61 1.31 2.92 2.92 2.92s2.92-1.31 2.92-2.92-1.31-2.92-2.92-2.92z" />
            </svg>
            Share
          </button>
          {view === 'posts' && (
            <button
              onClick={() => setShowCreatePost(true)}
              className="bg-[#5865f2] text-white text-sm px-3 py-1 rounded hover:bg-[#4752c4] transition-colors flex items-center gap-1"
            >
              <svg viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4">
                <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z" />
              </svg>
              New Post
            </button>
          )}
        </div>
      </header>

      {/* ── Body row: content + members panel ── */}
      <div className="flex-1 flex overflow-hidden">

        {/* ── Content + Compose ── */}
        <div className="flex-1 flex flex-col overflow-hidden">
          {view === 'posts' ? (
            <div className="flex-1 overflow-y-auto p-4 space-y-3">
              {posts.length === 0 ? (
                <div className="flex flex-col items-center justify-center h-full text-center py-16">
                  <div className="text-5xl mb-4">📋</div>
                  <div className="text-[#b5bac1] font-medium mb-1">No posts yet</div>
                  <div className="text-[#6d6f78] text-sm">Be the first to post something!</div>
                  <button
                    onClick={() => setShowCreatePost(true)}
                    className="mt-4 bg-[#5865f2] text-white px-4 py-2 rounded-md text-sm hover:bg-[#4752c4] transition-colors"
                  >
                    Create Post
                  </button>
                </div>
              ) : (
                posts.map(post => (
                  <PostCard
                    key={post.post_id}
                    post={post}
                    currentUser={currentUser}
                    author={userCache[post.author_id]}
                    onUpdated={handlePostUpdated}
                    onDeleted={handlePostDeleted}
                  />
                ))
              )}
            </div>
          ) : (
            <div className="flex-1 overflow-y-auto p-4 space-y-1">
              {messages.length === 0 ? (
                <div className="flex flex-col items-center justify-center h-full text-center py-16">
                  <div className="text-5xl mb-4">💬</div>
                  <div className="text-[#b5bac1] font-medium mb-1">No messages yet</div>
                  <div className="text-[#6d6f78] text-sm">Be the first to say something!</div>
                </div>
              ) : (
                messages.map(msg => {
                  const author = userCache[msg.author_id]
                  const initial = author ? author.username[0].toUpperCase() : '?'
                  const name = author ? author.username : msg.author_id.slice(0, 8)
                  const time = new Date(msg.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
                  return (
                    <div key={msg.message_id} className="flex items-start gap-3 px-2 py-1 rounded hover:bg-[#2e3035] group">
                      <div className="flex-shrink-0 w-8 h-8 rounded-full bg-[#5865f2] flex items-center justify-center text-white text-sm font-bold">
                        {initial}
                      </div>
                      <div className="min-w-0">
                        <div className="flex items-baseline gap-2">
                          <span className="text-[#f2f3f5] text-sm font-medium">{name}</span>
                          <span className="text-[#4e5058] text-xs">{time}</span>
                        </div>
                        <div className="text-[#dcddde] text-sm break-words">{msg.content}</div>
                      </div>
                    </div>
                  )
                })
              )}
              <div ref={messagesEndRef} />
            </div>
          )}

          {/* ── Compose / Message bar ── */}
          <div className="flex-shrink-0 px-4 pb-4 pt-2">
            {view === 'posts' ? (
              <button
                onClick={() => setShowCreatePost(true)}
                className="w-full bg-[#383a40] text-[#6d6f78] rounded-lg px-4 py-2.5 text-left hover:bg-[#404249] transition-colors text-sm"
              >
                Create a new post...
              </button>
            ) : (
              <form onSubmit={handleSendMessage} className="flex items-center gap-2 bg-[#383a40] rounded-full px-4 py-2">
                <input
                  value={messageInput}
                  onChange={e => setMessageInput(e.target.value)}
                  placeholder="Type a message..."
                  className="flex-1 bg-transparent text-[#dcddde] placeholder-[#6d6f78] text-sm outline-none"
                />
                <button
                  type="submit"
                  disabled={sendingMsg || !messageInput.trim()}
                  className="flex-shrink-0 w-7 h-7 bg-[#5865f2] hover:bg-[#4752c4] disabled:opacity-40 disabled:cursor-not-allowed rounded-full flex items-center justify-center transition-colors"
                >
                  <svg viewBox="0 0 24 24" fill="currentColor" className="w-3.5 h-3.5 text-white">
                    <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z" />
                  </svg>
                </button>
              </form>
            )}
          </div>
        </div>

        {/* ── Members panel ── */}
        <aside className="w-48 flex-shrink-0 border-l border-[#1e1f22] bg-[#2b2d31] overflow-y-auto">
          <h3 className="px-3 pt-4 pb-2 text-[#949ba4] text-xs font-semibold uppercase tracking-wide">
            Members — {server.member_ids.length}
          </h3>
          <div className="px-2 pb-4 space-y-0.5">
            {server.member_ids.map(id => {
              const member = userCache[id]
              const initial = member ? member.username[0].toUpperCase() : '?'
              const name = member ? member.username : id.slice(0, 8)
              const isOwner = id === server.owner_id
              return (
                <div key={id} className="flex items-center gap-2 px-1 py-1.5 rounded hover:bg-[#35373c] transition-colors">
                  <div className="flex-shrink-0 w-8 h-8 rounded-full bg-[#5865f2] flex items-center justify-center text-white text-sm font-bold">
                    {initial}
                  </div>
                  <div className="min-w-0 flex-1">
                    <div className="text-[#b5bac1] text-sm truncate">{name}</div>
                    {isOwner && (
                      <div className="text-[#f0b132] text-xs leading-none">Owner</div>
                    )}
                  </div>
                </div>
              )
            })}
          </div>
        </aside>

      </div>

      {showCreatePost && (
        <CreatePostModal
          serverId={serverId}
          currentUser={currentUser}
          onCreated={handlePostCreated}
          onClose={() => setShowCreatePost(false)}
        />
      )}
    </div>
  )
}
