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

  // Send a user message and get AI response via SSE streaming
  async sendMessage(
    chatId: string, 
    userMessageId: string, 
    content: string,
    onChunk: (chunk: string) => void,
    model?: string
  ): Promise<Message> {
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
      throw new Error('Failed to send message')
    }

    // Handle SSE stream
    const reader = response.body?.getReader()
    const decoder = new TextDecoder()
    
    if (!reader) {
      throw new Error('No response body')
    }

    let assistantMessage: Message | null = null

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      const chunk = decoder.decode(value)
      const lines = chunk.split('\n')

      for (const line of lines) {
        if (line.startsWith('data: ')) {
          const data = JSON.parse(line.slice(6))
          
          if (data.error) {
            throw new Error(data.error)
          }
          
          if (data.chunk) {
            onChunk(data.chunk)
          }
          
          if (data.done && data.content) {
            assistantMessage = {
              id: data.id,
              role: 'assistant',
              content: data.content
            }
          }
        }
      }
    }

    if (!assistantMessage) {
      throw new Error('No assistant message received')
    }

    return assistantMessage
  }
}

