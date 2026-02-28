import { useState, useEffect, useRef } from 'react'
import { User, Message } from '../types'
import { api } from '../api/client'

interface Props {
  serverId: string | null
  currentUser: User
}

export default function ChatPanel({ serverId, currentUser }: Props) {
  const [messages, setMessages] = useState<Message[]>([])
  const [userCache, setUserCache] = useState<Record<string, User>>({})
  const [input, setInput] = useState('')
  const [sending, setSending] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!serverId) {
      setMessages([])
      return
    }

    let cancelled = false
    setMessages([])

    async function load() {
      const msgs = await api.getMessages(serverId!).catch(() => [] as Message[])
      if (cancelled) return
      setMessages(msgs ?? [])

      const authorIds = new Set<string>((msgs ?? []).map((m: Message) => m.author_id))
      const entries = await Promise.all(
        [...authorIds].map(id =>
          api.getUser(id).then((u: User) => [id, u] as [string, User]).catch(() => null)
        )
      )
      if (cancelled) return
      const cache: Record<string, User> = {}
      entries.forEach(e => { if (e) cache[e[0]] = e[1] })
      setUserCache(cache)
    }

    load()
    return () => { cancelled = true }
  }, [serverId])

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const ensureUser = async (id: string) => {
    if (userCache[id]) return
    try {
      const u: User = await api.getUser(id)
      setUserCache(prev => ({ ...prev, [id]: u }))
    } catch { /* ignore */ }
  }

  const handleSend = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!input.trim() || sending || !serverId) return
    setSending(true)
    try {
      const msg: Message = await api.createMessage(serverId, currentUser.user_id, input.trim())
      setMessages(prev => [...prev, msg])
      setInput('')
      await ensureUser(msg.author_id)
    } finally {
      setSending(false)
    }
  }

  if (!serverId) return null

  return (
    <div className="flex-shrink-0 flex flex-col border-t border-[#1e1f22]">
      <div className="px-3 pt-3 pb-1">
        <span className="text-xs font-semibold text-[#949ba4] uppercase tracking-wide">Chat</span>
      </div>

      {/* Message list */}
      <div className="h-48 overflow-y-auto px-2 space-y-0.5">
        {messages.length === 0 ? (
          <div className="flex items-center justify-center h-full">
            <span className="text-[#6d6f78] text-xs">No messages yet</span>
          </div>
        ) : (
          messages.map(msg => {
            const author = userCache[msg.author_id]
            const initial = author ? author.username[0].toUpperCase() : '?'
            const name = author ? author.username : msg.author_id.slice(0, 8)
            const time = new Date(msg.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
            return (
              <div key={msg.message_id} className="flex items-start gap-2 px-1 py-0.5 rounded hover:bg-[#35373c] group">
                <div className="flex-shrink-0 w-6 h-6 rounded-full bg-[#5865f2] flex items-center justify-center text-white text-xs font-bold mt-0.5">
                  {initial}
                </div>
                <div className="min-w-0">
                  <div className="flex items-baseline gap-1.5">
                    <span className="text-[#f2f3f5] text-xs font-medium">{name}</span>
                    <span className="text-[#4e5058] text-[10px]">{time}</span>
                  </div>
                  <div className="text-[#dcddde] text-xs break-words">{msg.content}</div>
                </div>
              </div>
            )
          })
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Compose */}
      <div className="px-2 py-2">
        <form onSubmit={handleSend} className="flex items-center gap-1.5 bg-[#383a40] rounded-full px-3 py-1.5">
          <input
            value={input}
            onChange={e => setInput(e.target.value)}
            placeholder="Message..."
            className="flex-1 bg-transparent text-[#dcddde] placeholder-[#6d6f78] text-xs outline-none"
          />
          <button
            type="submit"
            disabled={sending || !input.trim()}
            className="flex-shrink-0 w-5 h-5 bg-[#5865f2] hover:bg-[#4752c4] disabled:opacity-40 disabled:cursor-not-allowed rounded-full flex items-center justify-center transition-colors"
          >
            <svg viewBox="0 0 24 24" fill="currentColor" className="w-3 h-3 text-white">
              <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z" />
            </svg>
          </button>
        </form>
      </div>
    </div>
  )
}
