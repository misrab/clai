export interface Message {
  id: string
  role: 'user' | 'assistant'
  content: string
}

export interface ChatTab {
  id: string
  title: string
  messages: Message[]
  input: string
}

