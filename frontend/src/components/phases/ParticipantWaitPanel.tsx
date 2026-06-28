import { useEffect, useState } from 'react'
import { IconBulb, IconSpinner } from '../ui/Icons'
import { getLocalizedTips, t } from '../../lib/i18n'

export function ParticipantWaitPanel() {
  const tips = getLocalizedTips()
  const [index, setIndex] = useState(0)
  const [paused, setPaused] = useState(false)
  const [transitionState, setTransitionState] = useState<
    'active' | 'exit' | 'enter'
  >('active')

  useEffect(() => {
    if (paused || tips.length === 0) {
      return
    }

    const timer = window.setInterval(() => {
      setTransitionState('exit')

      const exitTimer = window.setTimeout(() => {
        setIndex((prev) => (prev + 1) % tips.length)
        setTransitionState('enter')

        const enterTimer = window.setTimeout(() => {
          setTransitionState('active')
        }, 50)

        return () => window.clearTimeout(enterTimer)
      }, 300)

      return () => window.clearTimeout(exitTimer)
    }, 8000)

    return () => window.clearInterval(timer)
  }, [paused, tips.length])

  const currentTip = tips[index]

  let transitionClass = 'tip-panel__tip--active'
  if (transitionState === 'enter') {
    transitionClass = 'tip-panel__tip--enter'
  } else if (transitionState === 'exit') {
    transitionClass = 'tip-panel__tip--exit'
  }

  function goToTip(nextIndex: number) {
    setTransitionState('exit')
    window.setTimeout(() => {
      setIndex(nextIndex)
      setTransitionState('enter')
      window.setTimeout(() => {
        setTransitionState('active')
      }, 50)
    }, 300)
  }

  return (
    <div className="phase-panel tip-panel" aria-live="polite">
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

      <div className="tip-panel__controls">
        <button
          type="button"
          className="button button--secondary tip-panel__pause"
          onClick={() => setPaused((value) => !value)}
          aria-pressed={paused}
        >
          {paused ? t('wait.resumeCarousel') : t('wait.pauseCarousel')}
        </button>

        <div className="tip-panel__dots" aria-label={t('wait.carouselNav')}>
          {tips.map((_, idx) => (
            <button
              key={idx}
              type="button"
              className={`tip-panel__dot ${idx === index ? 'tip-panel__dot--active' : ''}`}
              onClick={() => goToTip(idx)}
              aria-label={t('wait.goTip', { idx: idx + 1 })}
              aria-current={idx === index ? 'true' : 'false'}
            />
          ))}
        </div>
      </div>

      <p className="tip-panel__status">
        <IconSpinner className="button__spinner" />
        <span>{t('wait.waitingForAdmin')}</span>
      </p>
    </div>
  )
}
