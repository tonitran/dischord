import { useState, useEffect } from 'react'
import { User, Post } from '../types'
import { api } from '../api/client'
import Avatar from './Avatar'

interface Props {
  post: Post
  currentUser: User
  author?: User
  onUpdated: (post: Post) => void
  onDeleted: (postId: string) => void
}

export default function PostCard({ post, currentUser, author, onUpdated, onDeleted }: Props) {
  const [editing, setEditing] = useState(false)
  const [editTitle, setEditTitle] = useState(post.title)
  const [editBody, setEditBody] = useState(post.body)
  const [saving, setSaving] = useState(false)
  const [deleting, setDeleting] = useState(false)
  const [myVote, setMyVote] = useState<number>(0)
  const [voteLoading, setVoteLoading] = useState(false)

  useEffect(() => {
    let cancelled = false
    api.getVote(post.server_id, post.post_id, currentUser.user_id)
      .then(v => { if (!cancelled) setMyVote(v.vote) })
      .catch(() => { /* 404 = no vote yet, stays 0 */ })
    return () => { cancelled = true }
  }, [post.post_id, currentUser.user_id])

  const handleVote = async (value: -1 | 1) => {
    if (voteLoading) return
    const next = myVote === value ? 0 : value
    const voteDelta = next - myVote
    setVoteLoading(true)
    try {
      await api.putVote(post.server_id, post.post_id, currentUser.user_id, next)
      setMyVote(next)
      onUpdated({ ...post, votes: post.votes + voteDelta })
    } finally {
      setVoteLoading(false)
    }
  }

  const isOwn = post.author_id === currentUser.user_id
  const authorName = author?.username ?? 'Unknown'
  const wasEdited = post.updated_at !== post.created_at

  const handleSave = async () => {
    if (!editTitle.trim()) return
    setSaving(true)
    try {
      const updated = await api.updatePost(post.server_id, post.post_id, editTitle.trim(), editBody.trim())
      onUpdated(updated)
      setEditing(false)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async () => {
    if (!confirm('Delete this post?')) return
    setDeleting(true)
    try {
      await api.deletePost(post.server_id, post.post_id)
      onDeleted(post.post_id)
    } finally {
      setDeleting(false)
    }
  }

  const formatTime = (iso: string) => {
    const d = new Date(iso)
    const now = new Date()
    const diffMs = now.getTime() - d.getTime()
    const diffMins = Math.floor(diffMs / 60000)
    if (diffMins < 1) return 'just now'
    if (diffMins < 60) return `${diffMins}m ago`
    const diffHrs = Math.floor(diffMins / 60)
    if (diffHrs < 24) return `${diffHrs}h ago`
    return d.toLocaleDateString()
  }

  if (editing) {
    return (
      <div className="bg-[#2b2d31] rounded-lg p-4 border-2 border-[#5865f2]">
        <input
          type="text"
          value={editTitle}
          onChange={e => setEditTitle(e.target.value)}
          className="input-field mb-2 text-lg font-semibold"
          placeholder="Title"
        />
        <textarea
          value={editBody}
          onChange={e => setEditBody(e.target.value)}
          className="input-field mb-3 resize-none"
          rows={4}
          placeholder="Body (optional)"
        />
        <div className="flex gap-2 justify-end">
          <button
            onClick={() => { setEditing(false); setEditTitle(post.title); setEditBody(post.body) }}
            className="px-3 py-1.5 text-sm text-[#949ba4] hover:text-white transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleSave}
            disabled={saving || !editTitle.trim()}
            className="px-3 py-1.5 text-sm bg-[#5865f2] text-white rounded-md hover:bg-[#4752c4] transition-colors disabled:opacity-50"
          >
            {saving ? 'Saving...' : 'Save'}
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="bg-[#2b2d31] rounded-lg p-4 hover:bg-[#32353b] transition-colors group">
      <div className="flex gap-3">
        <Avatar name={authorName} size="md" />
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1 flex-wrap">
            <span className={`text-sm font-semibold ${isOwn ? 'text-[#5865f2]' : 'text-white'}`}>
              {authorName}
            </span>
            {isOwn && (
              <span className="text-[10px] bg-[#5865f2]/30 text-[#5865f2] px-1.5 py-0.5 rounded-full font-medium">
                you
              </span>
            )}
            <span className="text-xs text-[#6d6f78]">{formatTime(post.created_at)}</span>
            {wasEdited && <span className="text-xs text-[#6d6f78]">(edited)</span>}
          </div>
          <h3 className="text-white font-semibold text-base mb-1 leading-snug">{post.title}</h3>
          {post.body && (
            <p className="text-[#b5bac1] text-sm whitespace-pre-wrap leading-relaxed">{post.body}</p>
          )}
          <div className="flex items-center gap-1 mt-2">
            <button
              onClick={() => handleVote(1)}
              disabled={voteLoading}
              className={`p-1 rounded transition-colors disabled:opacity-40 ${
                myVote === 1
                  ? 'text-[#57f287]'
                  : 'text-[#6d6f78] hover:text-[#57f287] hover:bg-[#383a40]'
              }`}
              title="Upvote"
            >
              <svg viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4">
                <path d="M7 14l5-5 5 5H7z" />
              </svg>
            </button>
            <span className={`text-xs font-semibold min-w-[1.5rem] text-center tabular-nums ${
              post.votes > 0 ? 'text-[#57f287]' : post.votes < 0 ? 'text-[#f23f43]' : 'text-[#949ba4]'
            }`}>
              {post.votes}
            </span>
            <button
              onClick={() => handleVote(-1)}
              disabled={voteLoading}
              className={`p-1 rounded transition-colors disabled:opacity-40 ${
                myVote === -1
                  ? 'text-[#f23f43]'
                  : 'text-[#6d6f78] hover:text-[#f23f43] hover:bg-[#383a40]'
              }`}
              title="Downvote"
            >
              <svg viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4">
                <path d="M7 10l5 5 5-5H7z" />
              </svg>
            </button>
          </div>
        </div>
        {isOwn && (
          <div className="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0 self-start">
            <button
              onClick={() => setEditing(true)}
              className="p-1.5 text-[#949ba4] hover:text-[#5865f2] hover:bg-[#383a40] rounded-md transition-colors"
              title="Edit post"
            >
              <svg viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4">
                <path d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04a1 1 0 0 0 0-1.41l-2.34-2.34a1 1 0 0 0-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z" />
              </svg>
            </button>
            <button
              onClick={handleDelete}
              disabled={deleting}
              className="p-1.5 text-[#949ba4] hover:text-[#f23f43] hover:bg-[#383a40] rounded-md transition-colors disabled:opacity-50"
              title="Delete post"
            >
              <svg viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4">
                <path d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z" />
              </svg>
            </button>
          </div>
        )}
      </div>
    </div>
  )
}
