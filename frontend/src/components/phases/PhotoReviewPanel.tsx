import { useState, useRef, useEffect } from 'react'
import type { Photo } from '../../types/lobby'
import { t } from '../../lib/i18n'
import { IconClose, IconLink } from '../ui/Icons'

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
  const [selectedPhoto, setSelectedPhoto] = useState<Photo | null>(null)
  const dialogRef = useRef<HTMLDialogElement>(null)

  useEffect(() => {
    const dialog = dialogRef.current
    if (!dialog) return

    const handleClose = () => {
      setSelectedPhoto(null)
    }

    dialog.addEventListener('close', handleClose)
    return () => {
      dialog.removeEventListener('close', handleClose)
    }
  }, [])

  useEffect(() => {
    const dialog = dialogRef.current
    if (!dialog) return

    if (selectedPhoto) {
      if (!dialog.open) {
        dialog.showModal()
      }
    } else {
      if (dialog.open) {
        dialog.close()
      }
    }
  }, [selectedPhoto])

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
                onClick={() => setSelectedPhoto(photo)}
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
        aria-label="Photo Preview"
      >
        {selectedPhoto && (
          <div className="photo-modal__content">
            <button
              type="button"
              className="photo-modal__close-btn"
              onClick={() => setSelectedPhoto(null)}
              aria-label="Close preview"
            >
              <IconClose />
            </button>
            <img
              src={selectedPhoto.large_image_url}
              alt=""
              className="photo-modal__image"
            />
            <div className="photo-modal__footer">
              <a
                href={selectedPhoto.page_url}
                target="_blank"
                rel="noopener noreferrer"
                className="photo-modal__link"
              >
                <IconLink className="icon" />
                <span>{t('draw.attribution')} Pixabay</span>
              </a>
            </div>
          </div>
        )}
      </dialog>
    </section>
  )
}
