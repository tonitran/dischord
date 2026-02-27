export interface User {
  user_id: string
  username: string
  email: string
  server_ids: string[]
  created_at: string
}

export interface Server {
  server_id: string
  name: string
  owner_id: string
  member_ids: string[]
  post_ids: string[]
  created_at: string
}

export interface Post {
  post_id: string
  server_id: string
  author_id: string
  title: string
  body: string
  votes: number
  created_at: string
  updated_at: string
}

export interface Message {
  message_id: string
  server_id: string
  author_id: string
  content: string
  created_at: string
}
