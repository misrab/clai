import { useState } from 'react'
import { Plus, X } from 'lucide-react'
import { ChatTab } from '../types'

interface TabBarProps {
  tabs: ChatTab[]
  activeTabId: string
  onTabSelect: (tabId: string) => void
  onTabClose: (tabId: string) => void
  onTabRename: (tabId: string, title: string) => void
  onNewTab: () => void
}

export function TabBar({ tabs, activeTabId, onTabSelect, onTabClose, onTabRename, onNewTab }: TabBarProps) {
  const [editingTabId, setEditingTabId] = useState<string | null>(null)
  const [editValue, setEditValue] = useState('')

  const handleDoubleClick = (tab: ChatTab) => {
    setEditingTabId(tab.id)
    setEditValue(tab.title)
  }

  const handleKeyDown = (e: React.KeyboardEvent, tabId: string) => {
    if (e.key === 'Enter') {
      e.preventDefault()
      finishEditing(tabId)
    } else if (e.key === 'Escape') {
      cancelEditing()
    }
  }

  const finishEditing = (tabId: string) => {
    if (editValue.trim()) {
      onTabRename(tabId, editValue.trim())
    }
    setEditingTabId(null)
  }

  const cancelEditing = () => {
    setEditingTabId(null)
    setEditValue('')
  }

  return (
    <div className="tab-bar">
      {tabs.map(tab => (
        <div
          key={tab.id}
          className={`tab ${tab.id === activeTabId ? 'active' : ''}`}
          onClick={() => onTabSelect(tab.id)}
        >
          {editingTabId === tab.id ? (
            <input
              type="text"
              className="tab-title-input"
              value={editValue}
              onChange={(e) => setEditValue(e.target.value)}
              onBlur={() => finishEditing(tab.id)}
              onKeyDown={(e) => handleKeyDown(e, tab.id)}
              autoFocus
              onClick={(e) => e.stopPropagation()}
            />
          ) : (
            <span
              className="tab-title"
              onDoubleClick={(e) => {
                e.stopPropagation()
                handleDoubleClick(tab)
              }}
            >
              {tab.title}
            </span>
          )}
          {tabs.length > 1 && (
            <button
              className="tab-close"
              onClick={(e) => {
                e.stopPropagation()
                onTabClose(tab.id)
              }}
              aria-label={`Close ${tab.title}`}
            >
              <X size={14} />
            </button>
          )}
        </div>
      ))}
      <button
        className="tab new-tab"
        onClick={onNewTab}
        aria-label="New chat"
      >
        <Plus size={16} />
      </button>
    </div>
  )
}

