export interface User {
  id: string
  email: string
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
  user: User
  expires_in: number
}

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8082'

export async function api<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const token = localStorage.getItem('access_token')
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...options.headers,
  }

  const res = await fetch(`${API_BASE}${path}`, { ...options, headers })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error || res.statusText)
  }
  const json = await res.json()
  return json.data !== undefined ? json.data : json
}

export async function register(email: string, password: string, masterKey: string) {
  return api<{ user_id: string }>('/register', {
    method: 'POST',
    body: JSON.stringify({ email, password, master_key: masterKey }),
  })
}

export async function login(email: string, password: string) {
  return api<AuthResponse>('/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  })
}

export async function getVault(masterKey: string) {
  const params = new URLSearchParams({ master_key: masterKey })
  const data = await api<{ vault: string }>(`/vault?${params}`)
  return JSON.parse(data.vault)
}

export async function saveVault(vault: any, masterKey: string) {
  return api('/vault', {
    method: 'POST',
    body: JSON.stringify({ vault: JSON.stringify(vault), master_key: masterKey }),
  })
}

export async function createPassword(entry: Omit<PasswordEntry, 'id'>, masterKey: string) {
  const data = await api<{ id: string }>('/vault/entries', {
    method: 'POST',
    body: JSON.stringify({ ...entry, master_key: masterKey }),
  })
  return data.id
}

export async function generatePassword(length = 16, options?: any) {
  return api<{ password: string }>('/generate', {
    method: 'POST',
    body: JSON.stringify({ length, ...options }),
  })
}

export async function searchPasswords(query: string) {
  return api<PasswordEntry[]>(`/vault/search?q=${encodeURIComponent(query)}`)
}

export function isAuthenticated(): boolean {
  return !!localStorage.getItem('access_token')
}