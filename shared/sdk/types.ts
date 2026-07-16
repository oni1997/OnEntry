export interface PasswordOptions {
  length: number
  uppercase: boolean
  lowercase: boolean
  numbers: boolean
  symbols: boolean
  excludeSimilar: boolean
  pronounceable: boolean
}

export interface PasswordEntry {
  id: string
  title: string
  username: string
  password: string
  website: string
  notes: string
  folder: string
  favorite: boolean
  tags: string[]
  created_at: string
  updated_at: string
}

export interface AuthResponse {
  access_token: string
  refresh_token: string
  user: {
    id: string
    email: string
  }
  expires_in: number
}