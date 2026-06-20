import type { Photo } from '../../types/lobby'
import { t } from '../../lib/i18n'

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
    </section>
  )
}
