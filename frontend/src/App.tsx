import { useState } from 'react'
import { Routes, Route, NavLink } from 'react-router-dom'
import Dashboard from './pages/Dashboard'
import Launchpad from './pages/Launchpad'
import Settings from './pages/Settings'

function App() {
  const [resultText, setResultText] = useState<string>("Ready to scan...");

  return (
    <div id="app" className="flex h-screen bg-gray-900 text-gray-100">
      <aside className="w-64 bg-gray-800 shadow-lg flex flex-col">
        <div className="p-6 text-2xl font-semibold text-purple-400">LMNinja</div>
        <nav className="flex-grow">
          <NavLink
            to="/dashboard"
            className={({ isActive }) =>
              `block py-3 px-6 text-lg hover:bg-gray-700 transition-colors duration-200 ${isActive ? 'bg-gray-700 border-l-4 border-purple-500' : ''}`
            }
          >
            Dashboard
          </NavLink>
          <NavLink
            to="/launchpad"
            className={({ isActive }) =>
              `block py-3 px-6 text-lg hover:bg-gray-700 transition-colors duration-200 ${isActive ? 'bg-gray-700 border-l-4 border-purple-500' : ''}`
            }
          >
            Launchpad
          </NavLink>
          <NavLink
            to="/settings"
            className={({ isActive }) =>
              `block py-3 px-6 text-lg hover:bg-gray-700 transition-colors duration-200 ${isActive ? 'bg-gray-700 border-l-4 border-purple-500' : ''}`
            }
          >
            Settings
          </NavLink>
        </nav>
        <div className="p-6 text-sm text-gray-500">
          LMNinja v0.1.0
        </div>
      </aside>

      <main className="flex-1 p-8 overflow-auto">
        <Routes>
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/launchpad" element={<Launchpad setResultText={setResultText} />} />
          <Route path="/settings" element={<Settings />} />
          <Route path="*" element={<Dashboard />} /> {/* Default route */}
        </Routes>
      </main>
    </div>
  )
}

export default App