import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { createLobby } from '../api/lobbyApi'
import {
  IconLink,
  IconSpinner,
  IconTimer,
  IconUsers,
} from '../components/ui/Icons'
import { saveAdminToken } from '../lib/adminStorage'
import { t } from '../lib/i18n'
import { useAppName } from '../hooks/useAppName'
import { useAppLogo } from '../hooks/useAppLogo'
import { useAppLogoLink } from '../hooks/useAppLogoLink'

const STEPS = [
  {
    icon: IconLink,
    titleKey: 'home.step1.title',
    descKey: 'home.step1.desc',
  },
  {
    icon: IconUsers,
    titleKey: 'home.step2.title',
    descKey: 'home.step2.desc',
  },
  {
    icon: IconTimer,
    titleKey: 'home.step3.title',
    descKey: 'home.step3.desc',
  },
] as const

export function HomePage() {
  const navigate = useNavigate()
  const appName = useAppName()
  const appLogo = useAppLogo()
  const appLogoLink = useAppLogoLink()
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
        err instanceof Error ? err.message : t('home.createLobbyFailed')
      setError(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="home-page">
      <div className="home-page__shell">
        <section className="home-page__hero" aria-labelledby="home-title">
          <p className="home-page__eyebrow">{t('home.eyebrow')}</p>
          <div className="home-page__brand-header">
            <a href={appLogoLink} target="_blank" rel="noopener noreferrer" className="home-page__logo-link">
              <img src={appLogo} alt="" className="home-page__logo-img" />
            </a>
            {appName.trim() && <h1 id="home-title">{appName}</h1>}
          </div>
          <p className="home-page__lead">{t('home.lead')}</p>

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
                {t('home.creatingLobby')}
              </>
            ) : (
              t('home.createLobby')
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
            {t('home.howItWorks')}
          </h2>
          <ol className="home-page__steps-list">
            {STEPS.map((step, index) => {
              const Icon = step.icon
              return (
                <li
                  key={step.titleKey}
                  className="home-page__step"
                  style={{ animationDelay: `${index * 60}ms` }}
                >
                  <span className="home-page__step-icon" aria-hidden="true">
                    <Icon />
                  </span>
                  <div>
                    <h3>{t(step.titleKey)}</h3>
                    <p>{t(step.descKey)}</p>
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
