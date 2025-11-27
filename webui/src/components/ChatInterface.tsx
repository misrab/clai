import { useEffect, useRef } from 'react'
import { MessageSquare } from 'lucide-react'
import { ChatTab } from '../types'

interface ChatInterfaceProps {
  tab: ChatTab
  onInputChange: (input: string) => void
  onSend: () => void
}

export function ChatInterface({ tab, onInputChange, onSend }: ChatInterfaceProps) {
  const messagesEndRef = useRef<HTMLDivElement>(null)

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [tab.messages])

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      onSend()
    }
  }

  return (
    <div className="chat-interface">
      {/* Messages area */}
      <div className="messages">
        {tab.messages.length === 0 ? (
          <div className="empty-state">
            <MessageSquare size={48} strokeWidth={1.5} />
            <p>Start a conversation</p>
          </div>
        ) : (
          <>
            {tab.messages.map(message => (
              <div key={message.id} className={`message ${message.role}`}>
                <div className="message-content">{message.content}</div>
              </div>
            ))}
            <div ref={messagesEndRef} />
          </>
        )}
      </div>

      {/* Input area */}
      <div className="input-area">
        <textarea
          value={tab.input}
          onChange={(e) => onInputChange(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Type your message..."
          rows={1}
        />
        <button
          onClick={onSend}
          disabled={!tab.input.trim()}
          aria-label="Send message"
        >
          Send
        </button>
      </div>
    </div>
  )
}

