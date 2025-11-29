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

  // Send a user message and get AI response (saves user message, calls Ollama, saves AI response)
  async sendMessage(chatId: string, userMessageId: string, content: string, model?: string): Promise<Message> {
    const response = await fetch(`${API_BASE}/chats/${chatId}/send`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ 
        userMessageId, 
        content,
        ...(model && { model })
      })
    })
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Failed to send message')
    }
    return response.json()
  }
}

