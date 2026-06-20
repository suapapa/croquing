import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { createLobby } from '../api/lobbyApi'
import {
  IconLink,
  IconLogo,
  IconSpinner,
  IconTimer,
  IconUsers,
} from '../components/ui/Icons'
import { saveAdminToken } from '../lib/adminStorage'

const STEPS = [
  {
    icon: IconLink,
    title: 'Create & share',
    description: 'Start a lobby and send the link to your drawing group.',
  },
  {
    icon: IconUsers,
    title: 'Pick references',
    description: 'The admin selects photos from Pixabay for everyone to draw.',
  },
  {
    icon: IconTimer,
    title: 'Draw in sync',
    description: 'Timed rounds with the same photo and countdown for all.',
  },
] as const

export function HomePage() {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function handleCreateLobby() {
    setLoading(true)
    setError(null)

    try {
      const created = await createLobby()
      saveAdminToken(created.id, created.admin_token)
      navigate(`/lobby/${created.id}`)
    } catch (err) {
      const message =
        err instanceof Error ? err.message : 'Failed to create lobby'
      setError(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="home-page">
      <div className="home-page__shell">
        <section className="home-page__hero" aria-labelledby="home-title">
          <p className="home-page__eyebrow">Real-time croquis meetups</p>
          <div className="home-page__brand-header">
            <IconLogo className="home-page__logo" aria-hidden="true" />
            <h1 id="home-title">Croquis King</h1>
          </div>
          <p className="home-page__lead">
            Create a lobby, share the link, and draw together with synchronized
            photos and timers — no screen sharing required.
          </p>

          <button
            type="button"
            className="button button--primary button--large"
            onClick={() => void handleCreateLobby()}
            disabled={loading}
            aria-busy={loading}
          >
            {loading ? (
              <>
                <IconSpinner className="button__spinner" />
                Creating lobby…
              </>
            ) : (
              'Create lobby'
            )}
          </button>

          {error ? (
            <p className="home-page__error" role="alert">
              {error}
            </p>
          ) : null}
        </section>

        <section className="home-page__steps" aria-labelledby="how-it-works">
          <h2 id="how-it-works" className="home-page__steps-title">
            How it works
          </h2>
          <ol className="home-page__steps-list">
            {STEPS.map((step, index) => {
              const Icon = step.icon
              return (
                <li
                  key={step.title}
                  className="home-page__step"
                  style={{ animationDelay: `${index * 60}ms` }}
                >
                  <span className="home-page__step-icon" aria-hidden="true">
                    <Icon />
                  </span>
                  <div>
                    <h3>{step.title}</h3>
                    <p>{step.description}</p>
                  </div>
                </li>
              )
            })}
          </ol>
        </section>
      </div>
    </main>
  )
}
