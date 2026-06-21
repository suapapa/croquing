import { useEffect, useState } from 'react'
import { IconBulb, IconSpinner } from '../ui/Icons'
import { getLocalizedTips, t } from '../../lib/i18n'

export function ParticipantWaitPanel() {
  const tips = getLocalizedTips()
  const [index, setIndex] = useState(0)
  const [transitionState, setTransitionState] = useState<
    'active' | 'exit' | 'enter'
  >('active')

  useEffect(() => {
    const timer = setInterval(() => {
      setTransitionState('exit')

      const exitTimer = setTimeout(() => {
        setIndex((prev) => (prev + 1) % tips.length)
        setTransitionState('enter')

        const enterTimer = setTimeout(() => {
          setTransitionState('active')
        }, 50)

        return () => clearTimeout(enterTimer)
      }, 300)

      return () => clearTimeout(exitTimer)
    }, 8000)

    return () => clearInterval(timer)
  }, [tips.length])

  const currentTip = tips[index]

  let transitionClass = 'tip-panel__tip--active'
  if (transitionState === 'enter') {
    transitionClass = 'tip-panel__tip--enter'
  } else if (transitionState === 'exit') {
    transitionClass = 'tip-panel__tip--exit'
  }

  return (
    <div className="phase-panel glass-card tip-panel" aria-live="polite">
      <div className="tip-panel__icon-wrap">
        <IconBulb className="button__icon" />
      </div>
      <h2 className="tip-panel__title">{t('wait.settingPhotos')}</h2>
      <p className="tip-panel__subtitle">{t('wait.adminChoosing')}</p>

      <div className="tip-panel__carousel">
        {currentTip ? (
          <div className={`tip-panel__tip ${transitionClass}`}>
            <span className="tip-panel__tip-title">
              {t('wait.tipLabel', { title: currentTip.title })}
            </span>
            <p className="tip-panel__tip-text">{currentTip.text}</p>
          </div>
        ) : null}
      </div>

      <div className="tip-panel__dots" aria-label="Carousel navigation">
        {tips.map((_, idx) => (
          <button
            key={idx}
            type="button"
            className={`tip-panel__dot ${idx === index ? 'tip-panel__dot--active' : ''}`}
            onClick={() => {
              setTransitionState('exit')
              setTimeout(() => {
                setIndex(idx)
                setTransitionState('enter')
                setTimeout(() => {
                  setTransitionState('active')
                }, 50)
              }, 300)
            }}
            aria-label={t('wait.goTip', { idx: idx + 1 })}
            aria-current={idx === index ? 'true' : 'false'}
          />
        ))}
      </div>

      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          gap: 'var(--space-2)',
          marginTop: 'var(--space-8)',
          color: 'var(--color-accent)',
          fontWeight: 600,
          fontSize: '0.875rem',
        }}
      >
        <IconSpinner className="button__spinner" />
        <span>{t('wait.waitingForAdmin')}</span>
      </div>
    </div>
  )
}
