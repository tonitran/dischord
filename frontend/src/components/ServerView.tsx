import { useState, useEffect, useRef } from 'react'
import { User, Server, Post, Message } from '../types'
import { api } from '../api/client'
import Avatar from './Avatar'
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

        // Load posts and messages in parallel
        const [postResults, msgs] = await Promise.all([
          Promise.all(s.post_ids.map((id: string) => api.getPost(serverId, id).catch(() => null))),
          api.getMessages(serverId).catch(() => [] as Message[]),
        ])
        if (cancelled) return

        const validPosts = postResults.filter(Boolean) as Post[]
        setPosts(validPosts)
        setMessages(msgs ?? [])

        // Prefetch all unique authors
        const authorIds = new Set<string>([
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
    return () => { cancelled = true }
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
      const msg: Message = await api.createMessage(serverId, currentUser.id, messageInput.trim())
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
    setPosts(prev => prev.map(p => p.id === post.id ? post : p))
  }

  const handlePostDeleted = (postId: string) => {
    setPosts(prev => prev.filter(p => p.id !== postId))
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

  const serverSlug = server.name.toLowerCase().replace(/\s+/g, '-')

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
          {/* Share server ID */}
          <button
            onClick={() => { navigator.clipboard.writeText(server.id) }}
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

      {/* ── Content ── */}
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
                key={post.id}
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
        <div className="flex-1 overflow-y-auto px-4 py-3 space-y-0">
          {messages.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-full text-center py-16">
              <div className="text-5xl mb-4">💬</div>
              <div className="text-[#b5bac1] font-medium mb-1">No messages yet</div>
              <div className="text-[#6d6f78] text-sm">Say something to kick things off!</div>
            </div>
          ) : (
            messages.map((msg, i) => {
              const author = userCache[msg.author_id]
              const isOwn = msg.author_id === currentUser.id
              const prevMsg = messages[i - 1]
              const grouped = prevMsg && prevMsg.author_id === msg.author_id

              return (
                <div
                  key={msg.id}
                  className={`flex items-start gap-3 ${grouped ? 'mt-0.5' : 'mt-4'}`}
                >
                  {grouped ? (
                    <div className="w-8 flex-shrink-0" />
                  ) : (
                    <Avatar name={author?.username ?? '?'} size="sm" />
                  )}
                  <div className="min-w-0 flex-1">
                    {!grouped && (
                      <div className="flex items-baseline gap-2 mb-0.5">
                        <span className={`text-sm font-semibold ${isOwn ? 'text-[#5865f2]' : 'text-white'}`}>
                          {author?.username ?? 'Unknown'}
                        </span>
                        <span className="text-xs text-[#6d6f78]">
                          {new Date(msg.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                        </span>
                      </div>
                    )}
                    <p className="text-[#dbdee1] text-sm leading-relaxed break-words">{msg.content}</p>
                  </div>
                </div>
              )
            })
          )}
          <div ref={messagesEndRef} />
        </div>
      )}

      {/* ── Compose bar ── */}
      <div className="flex-shrink-0 px-4 pb-4 pt-2">
        {view === 'messages' ? (
          <form onSubmit={handleSendMessage} className="flex gap-2">
            <input
              type="text"
              value={messageInput}
              onChange={e => setMessageInput(e.target.value)}
              placeholder={`Message #${serverSlug}`}
              className="flex-1 bg-[#383a40] text-[#dbdee1] rounded-lg px-4 py-2.5 focus:outline-none focus:ring-1 focus:ring-[#5865f2] placeholder-[#6d6f78] text-sm"
            />
            <button
              type="submit"
              disabled={!messageInput.trim() || sendingMsg}
              className="bg-[#5865f2] text-white px-4 py-2.5 rounded-lg hover:bg-[#4752c4] disabled:opacity-40 disabled:cursor-not-allowed transition-colors flex-shrink-0"
            >
              <svg viewBox="0 0 24 24" fill="currentColor" className="w-5 h-5">
                <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z" />
              </svg>
            </button>
          </form>
        ) : (
          <button
            onClick={() => setShowCreatePost(true)}
            className="w-full bg-[#383a40] text-[#6d6f78] rounded-lg px-4 py-2.5 text-left hover:bg-[#404249] transition-colors text-sm"
          >
            Create a new post...
          </button>
        )}
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
