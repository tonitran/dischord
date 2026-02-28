import { useState } from 'react'
import { User, Post } from '../types'
import { api } from '../api/client'

interface Props {
  serverId: string
  currentUser: User
  onCreated: (post: Post) => void
  onClose: () => void
}

export default function CreatePostModal({ serverId, currentUser, onCreated, onClose }: Props) {
  const [title, setTitle] = useState('')
  const [body, setBody] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!title.trim()) return
    setError('')
    setLoading(true)
    try {
      const post = await api.createPost(serverId, currentUser.user_id, title.trim(), body.trim())
      onCreated(post)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to create post')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div
      className="modal-backdrop"
      onClick={e => { if (e.target === e.currentTarget) onClose() }}
    >
      <div className="modal-card max-w-lg">
        <h2 className="text-xl font-bold text-white mb-4">Create Post</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="field-label">
              Title <span className="text-[#f23f43]">*</span>
            </label>
            <input
              type="text"
              value={title}
              onChange={e => setTitle(e.target.value)}
              className="input-field"
              placeholder="What's on your mind?"
              required
              autoFocus
            />
          </div>
          <div>
            <label className="field-label">
              Body
            </label>
            <textarea
              value={body}
              onChange={e => setBody(e.target.value)}
              className="input-field resize-none"
              rows={5}
              placeholder="Add more details (optional)"
            />
          </div>
          {error && <p className="field-error">{error}</p>}
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
              disabled={loading || !title.trim()}
              className="px-4 py-2 text-sm bg-[#5865f2] text-white rounded-md hover:bg-[#4752c4] transition-colors disabled:opacity-50"
            >
              {loading ? 'Posting...' : 'Post'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
