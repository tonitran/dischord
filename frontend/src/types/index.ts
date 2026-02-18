export interface User {
  id: string
  username: string
  email: string
  created_at: string
}

export interface Server {
  id: string
  name: string
  owner_id: string
  member_ids: string[]
  post_ids: string[]
  created_at: string
}

export interface Post {
  id: string
  server_id: string
  author_id: string
  title: string
  body: string
  created_at: string
  updated_at: string
}

export interface Message {
  id: string
  server_id: string
  author_id: string
  content: string
  created_at: string
}
