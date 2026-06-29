import { useState, useEffect, useRef } from 'react'
import {
  SUPPORTED_LANGUAGES,
  setLanguage,
  useLanguage,
  type Language,
} from '../../lib/i18n'

export function LanguageSelector() {
  const currentLang = useLanguage()
  const [isOpen, setIsOpen] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false)
      }
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') {
        setIsOpen(false)
      }
    }

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside)
      document.addEventListener('keydown', handleKeyDown)
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
      document.removeEventListener('keydown', handleKeyDown)
    }
  }, [isOpen])

  const toggleDropdown = () => setIsOpen((prev) => !prev)

  const handleSelect = (lang: Language) => {
    setLanguage(lang)
    setIsOpen(false)
  }

  return (
    <div className="language-selector" ref={containerRef}>
      <button
        type="button"
        className="language-selector__trigger"
        onClick={toggleDropdown}
        aria-haspopup="listbox"
        aria-expanded={isOpen}
        aria-label="Change language"
      >
        <svg
          className="language-selector__icon"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          aria-hidden="true"
        >
          <circle cx="12" cy="12" r="10" />
          <line x1="2" y1="12" x2="22" y2="12" />
          <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z" />
        </svg>
        <span className="language-selector__label">
          {SUPPORTED_LANGUAGES[currentLang]}
        </span>
        <svg
          className={`language-selector__chevron ${isOpen ? 'language-selector__chevron--open' : ''}`}
          width="12"
          height="12"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2.5"
          strokeLinecap="round"
          strokeLinejoin="round"
          aria-hidden="true"
        >
          <polyline points="6 9 12 15 18 9" />
        </svg>
      </button>

      {isOpen && (
        <ul
          className="language-selector__dropdown"
          role="listbox"
          aria-label="Languages"
        >
          {(Object.keys(SUPPORTED_LANGUAGES) as Language[]).map((lang) => (
            <li key={lang} role="option" aria-selected={lang === currentLang}>
              <button
                type="button"
                className={`language-selector__option ${
                  lang === currentLang
                    ? 'language-selector__option--active'
                    : ''
                }`}
                onClick={() => handleSelect(lang)}
              >
                {SUPPORTED_LANGUAGES[lang]}
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
