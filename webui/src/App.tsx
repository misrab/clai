import { useState } from 'react'
import { MessageSquare } from 'lucide-react'
import { useTabs } from './hooks/useTabs'
import { TabBar } from './components/TabBar'
import { ChatInterface } from './components/ChatInterface'
import './App.css'

function App() {
  const [isCollapsed, setIsCollapsed] = useState(true)
  const [isMobileOpen, setIsMobileOpen] = useState(false)

  const {
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
  } = useTabs()

  if (loading) {
    return (
      <div className="app" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh' }}>
        <div>Loading...</div>
      </div>
    )
  }

  return (
    <div className="app">
      {/* Sidebar */}
      <div className={`sidebar ${isCollapsed ? 'collapsed' : ''} ${isMobileOpen ? 'mobile-open' : ''}`}>
        <div className="sidebar-content">
          <button 
            onClick={() => setIsCollapsed(!isCollapsed)} 
            className="collapse-btn"
            title={isCollapsed ? 'Expand' : 'Collapse'}
          >
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              {isCollapsed ? (
                <path d="M9 18l6-6-6-6" />
              ) : (
                <path d="M15 18l-6-6 6-6" />
              )}
            </svg>
          </button>
          <button 
            onClick={() => setIsMobileOpen(false)} 
            className="mobile-close-btn"
          >
            ×
          </button>
          <nav>
            <a href="#" className="nav-item active">
              <span className="nav-icon" aria-hidden="true">
                <MessageSquare size={18} />
              </span>
              <span className="nav-label">Chats</span>
            </a>
          </nav>
        </div>
      </div>

      {/* Backdrop for mobile */}
      {isMobileOpen && <div className="backdrop" onClick={() => setIsMobileOpen(false)} />}

      {/* Main content */}
      <div className={`main ${isCollapsed ? 'sidebar-collapsed' : 'sidebar-expanded'}`}>
        <header className="top-bar">
          <button
            onClick={() => setIsMobileOpen(true)}
            className="mobile-menu-btn"
            aria-label="Open navigation"
          >
            ☰
          </button>
          <div className="app-title">clai</div>
        </header>

        {/* Tab bar */}
        <TabBar
          tabs={tabs}
          activeTabId={activeTabId}
          onTabSelect={setActiveTabId}
          onTabClose={closeTab}
          onTabRename={updateTabTitle}
          onNewTab={addTab}
        />

        {/* Chat interface */}
        <div className="content">
          {activeTab && (
            <ChatInterface
              tab={activeTab}
              onInputChange={(input) => updateTabInput(activeTab.id, input)}
              onSend={() => sendMessage(activeTab.id)}
            />
          )}
        </div>
      </div>
    </div>
  )
}

export default App


