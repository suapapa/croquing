import { useEffect, useState } from 'react'
import { IconBulb, IconSpinner } from '../ui/Icons'

const TIPS = [
  {
    title: 'Gesture First',
    text: 'Capture the main line of action and body flow first, ignoring all small details.',
  },
  {
    title: 'Squint to Simplify',
    text: 'Squinting helps block out high-frequency details so you can see the main shapes of light and shadow.',
  },
  {
    title: 'Draw with your Arm',
    text: 'Use your elbow and shoulder to make long, sweeping strokes instead of scratchy wrist movements.',
  },
  {
    title: 'Exaggerate the Pose',
    text: 'In short gesture sketches, it is better to draw a pose that looks more dynamic and expressive than the photo.',
  },
  {
    title: 'Negative Spaces',
    text: 'Look at the empty shapes formed between limbs to gauge proportions and angles accurately.',
  },
  {
    title: 'Proportions over Details',
    text: 'Leave hands, feet, and faces blank until the main torso and overall skeleton structure are locked in.',
  },
  {
    title: 'Keep it Loose',
    text: 'Remember, this is a croquis! Speed and feeling are more important than a polished, perfect drawing.',
  },
]

export function ParticipantWaitPanel() {
  const [index, setIndex] = useState(0)
  const [transitionState, setTransitionState] = useState<'active' | 'exit' | 'enter'>('active')

  useEffect(() => {
    const timer = setInterval(() => {
      setTransitionState('exit')
      
      const exitTimer = setTimeout(() => {
        setIndex((prev) => (prev + 1) % TIPS.length)
        setTransitionState('enter')
        
        const enterTimer = setTimeout(() => {
          setTransitionState('active')
        }, 50)
        
        return () => clearTimeout(enterTimer)
      }, 300)
      
      return () => clearTimeout(exitTimer)
    }, 8000)

    return () => clearInterval(timer)
  }, [])

  const currentTip = TIPS[index]

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
      <h2 className="tip-panel__title">Setting up reference photos...</h2>
      <p style={{ color: 'var(--color-muted)', fontSize: '0.875rem', marginBottom: 'var(--space-6)' }}>
        The admin is choosing photos. They stay hidden until the session starts.
      </p>

      <div className="tip-panel__carousel">
        <div className={`tip-panel__tip ${transitionClass}`}>
          <span className="tip-panel__tip-title">Tip: {currentTip.title}</span>
          <p className="tip-panel__tip-text">{currentTip.text}</p>
        </div>
      </div>

      <div className="tip-panel__dots" aria-label="Carousel navigation">
        {TIPS.map((_, idx) => (
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
            aria-label={`Go to tip ${idx + 1}`}
            aria-current={idx === index ? 'true' : 'false'}
          />
        ))}
      </div>

      <div style={{ display: 'flex', alignItems: 'center', gap: 'var(--space-2)', marginTop: 'var(--space-8)', color: 'var(--color-accent)', fontWeight: 600, fontSize: '0.875rem' }}>
        <IconSpinner className="button__spinner" />
        <span>Waiting for admin...</span>
      </div>
    </div>
  )
}
