import { Navigate, Route, Routes } from 'react-router-dom'
import { HomePage } from './pages/HomePage'
import { LobbyPage } from './pages/LobbyPage'
import './index.css'

function App() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/lobby/:id" element={<LobbyPage />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

export default App
