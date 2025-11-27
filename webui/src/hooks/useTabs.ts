import { useState } from 'react'
import { ChatTab, Message } from '../types'
import { generateId } from '../utils'

export function useTabs() {
  const [tabs, setTabs] = useState<ChatTab[]>([
    {
      id: generateId(),
      title: 'New Chat',
      messages: [],
      input: ''
    }
  ])
  const [activeTabId, setActiveTabId] = useState(tabs[0].id)

  const addTab = () => {
    const newTab: ChatTab = {
      id: generateId(),
      title: 'New Chat',
      messages: [],
      input: ''
    }
    setTabs([...tabs, newTab])
    setActiveTabId(newTab.id)
  }

  const updateTabTitle = (tabId: string, title: string) => {
    setTabs(prevTabs => prevTabs.map(tab =>
      tab.id === tabId ? { ...tab, title } : tab
    ))
  }

  const closeTab = (tabId: string) => {
    if (tabs.length === 1) return // Keep at least one tab
    
    setTabs(prevTabs => {
      const newTabs = prevTabs.filter(tab => tab.id !== tabId)
      
      // If closing active tab, switch to another
      if (tabId === activeTabId) {
        setActiveTabId(newTabs[0].id)
      }
      
      return newTabs
    })
  }

  const updateTabInput = (tabId: string, input: string) => {
    setTabs(prevTabs => prevTabs.map(tab => 
      tab.id === tabId ? { ...tab, input } : tab
    ))
  }

  const addMessage = (tabId: string, message: Message) => {
    setTabs(prevTabs => prevTabs.map(tab =>
      tab.id === tabId
        ? { ...tab, messages: [...tab.messages, message] }
        : tab
    ))
  }

  const sendMessage = (tabId: string) => {
    const tab = tabs.find(t => t.id === tabId)
    if (!tab || !tab.input.trim()) return

    const userMessage: Message = {
      id: generateId(),
      role: 'user',
      content: tab.input
    }

    // Add user message and clear input
    setTabs(prevTabs => prevTabs.map(t => 
      t.id === tabId 
        ? { 
            ...t, 
            messages: [...t.messages, userMessage],
            input: '' 
          }
        : t
    ))

    // TODO: Call AI API and add assistant response
    // For now, just simulate a response
    setTimeout(() => {
      const assistantMessage: Message = {
        id: generateId(),
        role: 'assistant',
        content: 'This is a simulated response. AI integration coming soon!'
      }
      addMessage(tabId, assistantMessage)
    }, 500)
  }

  const activeTab = tabs.find(tab => tab.id === activeTabId)

  return {
    tabs,
    activeTab,
    activeTabId,
    setActiveTabId,
    addTab,
    closeTab,
    updateTabTitle,
    updateTabInput,
    sendMessage
  }
}

