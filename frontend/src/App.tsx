import { getApiBase } from './config/env'
import './index.css'

function App() {
  return (
    <main className="app">
      <h1>Croquis King</h1>
      <p className="tagline">
        Real-time croquis meetups — synchronized photos and timers.
      </p>
      <p className="meta">
        API base: <code>{getApiBase()}</code>
      </p>
    </main>
  )
}

export default App
