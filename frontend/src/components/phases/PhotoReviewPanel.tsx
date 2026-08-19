import { useState, useRef, useEffect } from 'react'
import type { Photo } from '../../types/lobby'
import { t } from '../../lib/i18n'
import {
  IconClose,
  IconLink,
  IconUser,
  IconShieldCheck,
  IconImage,
  IconChevronLeft,
  IconChevronRight,
} from '../ui/Icons'

interface PhotoReviewPanelProps {
  photos: Photo[]
  loading: boolean
  onEdit: () => void
  onConfirm: () => void
}

export function PhotoReviewPanel({
  photos,
  loading,
  onEdit,
  onConfirm,
}: PhotoReviewPanelProps) {
  const [selectedIndex, setSelectedIndex] = useState<number | null>(null)
  const dialogRef = useRef<HTMLDialogElement>(null)

  const selectedPhoto =
    selectedIndex !== null ? photos[selectedIndex] ?? null : null
  const canGoPrev = selectedIndex !== null && selectedIndex > 0
  const canGoNext =
    selectedIndex !== null && selectedIndex < photos.length - 1

  const goToPrev = () => {
    if (canGoPrev) {
      setSelectedIndex((index) => (index !== null ? index - 1 : null))
    }
  }

  const goToNext = () => {
    if (canGoNext) {
      setSelectedIndex((index) => (index !== null ? index + 1 : null))
    }
  }

  useEffect(() => {
    const dialog = dialogRef.current
    if (!dialog) return

    const handleClose = () => {
      setSelectedIndex(null)
    }

    dialog.addEventListener('close', handleClose)
    return () => {
      dialog.removeEventListener('close', handleClose)
    }
  }, [])

  useEffect(() => {
    const dialog = dialogRef.current
    if (!dialog) return

    if (selectedIndex !== null) {
      if (!dialog.open) {
        dialog.showModal()
      }
    } else if (dialog.open) {
      dialog.close()
    }
  }, [selectedIndex])

  useEffect(() => {
    if (selectedIndex === null) return

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'ArrowLeft') {
        event.preventDefault()
        setSelectedIndex((index) => {
          if (index === null || index <= 0) return index
          return index - 1
        })
      } else if (event.key === 'ArrowRight') {
        event.preventDefault()
        setSelectedIndex((index) => {
          if (index === null || index >= photos.length - 1) return index
          return index + 1
        })
      }
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => {
      window.removeEventListener('keydown', handleKeyDown)
    }
  }, [selectedIndex, photos.length])

  const handleBackdropClick = (e: React.MouseEvent<HTMLDialogElement>) => {
    const dialog = dialogRef.current
    if (!dialog) return
    if (e.target === dialog) {
      dialog.close()
    }
  }

  return (
    <section className="photo-review" aria-labelledby="photo-review-title">
      <header className="photo-review__header">
        <h2 id="photo-review-title">
          {t('review.photosSaved', { count: photos.length })}
        </h2>
        <p>{t('review.instruction')}</p>
      </header>

      {photos.length > 0 ? (
        <ul className="photo-review__list">
          {photos.map((photo, index) => (
            <li key={photo.pixabay_id} className="photo-review__item">
              <button
                type="button"
                className="photo-review__thumb"
                onClick={() => setSelectedIndex(index)}
                aria-label={t('review.previewAria', {
                  index: index + 1,
                  total: photos.length,
                })}
              >
                <img
                  className="photo-review__thumb-image"
                  src={photo.preview_url}
                  alt=""
                  loading="lazy"
                />
                <span className="photo-review__preview" aria-hidden="true">
                  <img
                    className="photo-review__preview-image"
                    src={photo.large_image_url}
                    alt=""
                    loading="lazy"
                  />
                </span>
              </button>
            </li>
          ))}
        </ul>
      ) : null}

      <div className="photo-review__actions">
        <button
          type="button"
          className="button button--secondary"
          disabled={loading}
          onClick={onEdit}
        >
          {t('review.editSelection')}
        </button>
        <button
          type="button"
          className="button button--primary"
          disabled={loading || photos.length === 0}
          onClick={onConfirm}
        >
          {loading ? t('review.shuffling') : t('review.selectionComplete')}
        </button>
      </div>

      <dialog
        ref={dialogRef}
        className="photo-modal"
        onClick={handleBackdropClick}
        aria-label={t('review.modalAria')}
      >
        {selectedPhoto && (
          <div className="photo-modal__content">
            <button
              type="button"
              className="photo-modal__close-btn"
              onClick={() => setSelectedIndex(null)}
              aria-label={t('review.closePreview')}
            >
              <IconClose />
            </button>
            <div className="photo-modal__body">
              <div className="photo-modal__image-wrapper">
                {photos.length > 1 ? (
                  <button
                    type="button"
                    className="photo-modal__nav-btn photo-modal__nav-btn--prev"
                    onClick={goToPrev}
                    disabled={!canGoPrev}
                    aria-label={t('review.prevPhoto')}
                  >
                    <IconChevronLeft />
                  </button>
                ) : null}
                <img
                  src={selectedPhoto.large_image_url}
                  alt={t('review.modalAlt', {
                    index: (selectedIndex ?? 0) + 1,
                    total: photos.length,
                  })}
                  className="photo-modal__image"
                />
                {photos.length > 1 ? (
                  <button
                    type="button"
                    className="photo-modal__nav-btn photo-modal__nav-btn--next"
                    onClick={goToNext}
                    disabled={!canGoNext}
                    aria-label={t('review.nextPhoto')}
                  >
                    <IconChevronRight />
                  </button>
                ) : null}
              </div>
              <aside className="photo-modal__info">
                <h3 className="photo-modal__info-title">
                  {t('review.modalAlt', {
                    index: (selectedIndex ?? 0) + 1,
                    total: photos.length,
                  })}
                </h3>
                <dl className="photo-modal__meta-list">
                  {selectedPhoto.user ? (
                    <div className="photo-modal__meta-item">
                      <dt className="photo-modal__meta-term">
                        <IconUser className="photo-modal__meta-icon" />
                        <span>{t('review.metaAuthor')}</span>
                      </dt>
                      <dd className="photo-modal__meta-value">
                        {selectedPhoto.user_id ? (
                          <a
                            href={`https://pixabay.com/users/${selectedPhoto.user}-${selectedPhoto.user_id}/`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="photo-modal__meta-author-link"
                          >
                            {selectedPhoto.user}
                          </a>
                        ) : (
                          <span>{selectedPhoto.user}</span>
                        )}
                      </dd>
                    </div>
                  ) : null}

                  <div className="photo-modal__meta-item">
                    <dt className="photo-modal__meta-term">
                      <IconShieldCheck className="photo-modal__meta-icon" />
                      <span>{t('review.metaLicense')}</span>
                    </dt>
                    <dd className="photo-modal__meta-value">
                      <a
                        href="https://pixabay.com/service/license-summary/"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="photo-modal__meta-license-link"
                      >
                        {t('review.metaLicenseFree')}
                      </a>
                    </dd>
                  </div>

                  {selectedPhoto.width && selectedPhoto.height ? (
                    <div className="photo-modal__meta-item">
                      <dt className="photo-modal__meta-term">
                        <IconImage className="photo-modal__meta-icon" />
                        <span>{t('review.metaDimensions')}</span>
                      </dt>
                      <dd className="photo-modal__meta-value">
                        {selectedPhoto.width} × {selectedPhoto.height} px
                      </dd>
                    </div>
                  ) : null}
                </dl>

                <div className="photo-modal__info-footer">
                  <a
                    href={selectedPhoto.page_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="button button--secondary photo-modal__pixabay-btn"
                  >
                    <IconLink className="icon" />
                    <span>{t('review.viewOnPixabay')}</span>
                  </a>
                </div>
              </aside>
            </div>
          </div>
        )}
      </dialog>
    </section>
  )
}
