import { Navigate, Route, Routes } from 'react-router-dom'
import { AppFooter } from './components/layout/AppFooter'
import { HomePage } from './pages/HomePage'
import { LobbyPage } from './pages/LobbyPage'
import { useLanguage } from './lib/i18n'
import './index.css'

function App() {
  useLanguage() // React to global language changes to trigger re-render

  return (
    <div className="app-shell">
      <div className="app-shell__main">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/lobby/:id" element={<LobbyPage />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </div>
      <AppFooter />
    </div>
  )
}

export default App
