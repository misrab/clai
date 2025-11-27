import { ChatTab, Message } from './types'

const API_BASE = '/api'

export interface ChatResponse {
  id: string
  title: string
  created_at: string
  updated_at: string
  messages: Message[]
}

export interface ChatListItem {
  id: string
  title: string
  created_at: string
  updated_at: string
}

// Chats API
export const api = {
  // List all chats
  async listChats(): Promise<ChatListItem[]> {
    const response = await fetch(`${API_BASE}/chats`)
    if (!response.ok) throw new Error('Failed to list chats')
    return response.json()
  },

  // Get a single chat with messages
  async getChat(id: string): Promise<ChatResponse> {
    const response = await fetch(`${API_BASE}/chats/${id}`)
    if (!response.ok) throw new Error('Failed to get chat')
    return response.json()
  },

  // Create a new chat
  async createChat(id: string, title: string): Promise<ChatResponse> {
    const response = await fetch(`${API_BASE}/chats`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id, title })
    })
    if (!response.ok) throw new Error('Failed to create chat')
    return response.json()
  },

  // Update chat title
  async updateChatTitle(id: string, title: string): Promise<void> {
    const response = await fetch(`${API_BASE}/chats/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title })
    })
    if (!response.ok) throw new Error('Failed to update chat title')
  },

  // Delete a chat
  async deleteChat(id: string): Promise<void> {
    const response = await fetch(`${API_BASE}/chats/${id}`, {
      method: 'DELETE'
    })
    if (!response.ok) throw new Error('Failed to delete chat')
  },

  // Add a message to a chat
  async createMessage(chatId: string, message: Message): Promise<Message> {
    const response = await fetch(`${API_BASE}/chats/${chatId}/messages`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(message)
    })
    if (!response.ok) throw new Error('Failed to create message')
    return response.json()
  }
}

