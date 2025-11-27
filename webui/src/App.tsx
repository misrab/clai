import { useState } from 'react'
import { MessageSquare } from 'lucide-react'
import './App.css'

function App() {
  const [isCollapsed, setIsCollapsed] = useState(false)
  const [isMobileOpen, setIsMobileOpen] = useState(false)

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
        <div className="content">
          {/* Chats page content will go here */}
        </div>
      </div>
    </div>
  )
}

export default App


