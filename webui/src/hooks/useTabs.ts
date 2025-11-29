import { useState, useEffect } from 'react'
import { ChatTab, Message } from '../types'
import { generateId } from '../utils'
import { api } from '../api'

export function useTabs() {
  const [tabs, setTabs] = useState<ChatTab[]>([])
  const [activeTabId, setActiveTabId] = useState<string>('')
  const [loading, setLoading] = useState(true)

  // Load chats from API on mount
  useEffect(() => {
    loadChats()
  }, [])

  const loadChats = async (retryCount = 0) => {
    try {
      const chatList = await api.listChats()
      
      if (chatList.length === 0) {
        // Create initial chat if none exist
        const newId = generateId()
        await api.createChat(newId, 'New Chat')
        setTabs([{
          id: newId,
          title: 'New Chat',
          messages: [],
          input: ''
        }])
        setActiveTabId(newId)
      } else {
        // Load existing chats (without messages initially)
        const loadedTabs: ChatTab[] = chatList.map(chat => ({
          id: chat.id,
          title: chat.title,
          messages: [], // Will be loaded when tab is activated
          input: ''
        }))
        setTabs(loadedTabs)
        setActiveTabId(loadedTabs[0].id)
        
        // Load messages for the first tab
        loadChatMessages(loadedTabs[0].id)
      }
      setLoading(false)
    } catch (error) {
      console.error('Failed to load chats:', error)
      
      // Retry up to 3 times if it's a connection error (backend not ready)
      if (retryCount < 3 && error instanceof TypeError) {
        console.log(`Backend not ready, retrying in 2 seconds... (attempt ${retryCount + 1}/3)`)
        setTimeout(() => loadChats(retryCount + 1), 2000)
        return
      }
      
      // Fallback to a default chat after retries exhausted
      console.log('Creating fallback chat...')
      const newId = generateId()
      setTabs([{
        id: newId,
        title: 'New Chat',
        messages: [],
        input: ''
      }])
      setActiveTabId(newId)
      setLoading(false)
    }
  }

  const loadChatMessages = async (chatId: string) => {
    try {
      const chat = await api.getChat(chatId)
      setTabs(prevTabs => prevTabs.map(tab =>
        tab.id === chatId
          ? { ...tab, messages: chat.messages }
          : tab
      ))
    } catch (error) {
      console.error('Failed to load chat messages:', error)
    }
  }

  // Load messages when switching tabs
  useEffect(() => {
    if (activeTabId) {
      const activeTab = tabs.find(t => t.id === activeTabId)
      if (activeTab && activeTab.messages.length === 0) {
        loadChatMessages(activeTabId)
      }
    }
  }, [activeTabId])

  const addTab = async () => {
    const newId = generateId()
    const title = 'New Chat'
    
    try {
      await api.createChat(newId, title)
      const newTab: ChatTab = {
        id: newId,
        title,
        messages: [],
        input: ''
      }
      setTabs([...tabs, newTab])
      setActiveTabId(newId)
    } catch (error) {
      console.error('Failed to create chat:', error)
    }
  }

  const updateTabTitle = async (tabId: string, title: string) => {
    try {
      await api.updateChatTitle(tabId, title)
      setTabs(prevTabs => prevTabs.map(tab =>
        tab.id === tabId ? { ...tab, title } : tab
      ))
    } catch (error) {
      console.error('Failed to update chat title:', error)
    }
  }

  const closeTab = async (tabId: string) => {
    if (tabs.length === 1) return // Keep at least one tab
    
    try {
      await api.deleteChat(tabId)
      
      setTabs(prevTabs => {
        const newTabs = prevTabs.filter(tab => tab.id !== tabId)
        
        // If closing active tab, switch to another
        if (tabId === activeTabId) {
          setActiveTabId(newTabs[0].id)
        }
        
        return newTabs
      })
    } catch (error) {
      console.error('Failed to delete chat:', error)
    }
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

  const sendMessage = async (tabId: string) => {
    const tab = tabs.find(t => t.id === tabId)
    if (!tab || !tab.input.trim()) return

    const userMessageId = generateId()
    const userMessage: Message = {
      id: userMessageId,
      role: 'user',
      content: tab.input
    }

    const messageContent = tab.input
    const assistantMessageId = generateId()

    // Add user message to local state and clear input
    setTabs(prevTabs => prevTabs.map(t => 
      t.id === tabId 
        ? { 
            ...t, 
            messages: [...t.messages, userMessage],
            input: '' 
          }
        : t
    ))

    // Add empty assistant message placeholder for streaming
    const assistantMessage: Message = {
      id: assistantMessageId,
      role: 'assistant',
      content: ''
    }
    addMessage(tabId, assistantMessage)

    try {
      // Send message and handle streaming chunks
      await api.sendMessage(
        tabId, 
        userMessageId, 
        messageContent,
        (chunk: string) => {
          // Update assistant message with new chunk
          setTabs(prevTabs => prevTabs.map(t => {
            if (t.id !== tabId) return t
            
            return {
              ...t,
              messages: t.messages.map(m => 
                m.id === assistantMessageId
                  ? { ...m, content: m.content + chunk }
                  : m
              )
            }
          }))
        }
      )
    } catch (error) {
      console.error('Failed to send message:', error)
      
      // Replace placeholder with error message
      setTabs(prevTabs => prevTabs.map(t => {
        if (t.id !== tabId) return t
        
        return {
          ...t,
          messages: t.messages.map(m => 
            m.id === assistantMessageId
              ? { 
                  ...m, 
                  content: `Error: ${error instanceof Error ? error.message : 'Failed to get AI response. Make sure Ollama is running.'}`
                }
              : m
          )
        }
      }))
    }
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
    sendMessage,
    loading
  }
}

