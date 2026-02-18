const BASE = '/api'

async function apiFetch(path: string, options?: RequestInit) {
  const res = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json', ...options?.headers },
    ...options,
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || `HTTP ${res.status}`)
  }
  if (res.status === 204) return null
  return res.json()
}

export const api = {
  // Users
  createUser: (username: string, email: string) =>
    apiFetch('/users', { method: 'POST', body: JSON.stringify({ username, email }) }),

  getUser: (id: string) =>
    apiFetch(`/users/${id}`),

  // Friends
  addFriend: (userId: string, friendId: string) =>
    apiFetch(`/users/${userId}/friends`, {
      method: 'POST',
      body: JSON.stringify({ friend_id: friendId }),
    }),

  getFriends: (userId: string) =>
    apiFetch(`/users/${userId}/friends`),

  // Servers
  createServer: (name: string, ownerId: string) =>
    apiFetch('/servers', {
      method: 'POST',
      body: JSON.stringify({ name, owner_id: ownerId }),
    }),

  getServer: (id: string) =>
    apiFetch(`/servers/${id}`),

  // Posts
  createPost: (serverId: string, authorId: string, title: string, body: string) =>
    apiFetch(`/servers/${serverId}/posts`, {
      method: 'POST',
      body: JSON.stringify({ author_id: authorId, title, body }),
    }),

  getPost: (serverId: string, postId: string) =>
    apiFetch(`/servers/${serverId}/posts/${postId}`),

  updatePost: (serverId: string, postId: string, title: string, body: string) =>
    apiFetch(`/servers/${serverId}/posts/${postId}`, {
      method: 'PUT',
      body: JSON.stringify({ title, body }),
    }),

  deletePost: (serverId: string, postId: string) =>
    apiFetch(`/servers/${serverId}/posts/${postId}`, { method: 'DELETE' }),

  // Messages
  createMessage: (serverId: string, authorId: string, content: string) =>
    apiFetch(`/servers/${serverId}/messages`, {
      method: 'POST',
      body: JSON.stringify({ author_id: authorId, content }),
    }),

  getMessages: (serverId: string) =>
    apiFetch(`/servers/${serverId}/messages`),
}
