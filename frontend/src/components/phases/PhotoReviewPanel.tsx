import type { Photo } from '../../types/lobby'

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
        <h2 id="photo-review-title">{photos.length} photos saved</h2>
        <p>
          Hover a thumbnail to preview it at full size. When you are happy with
          the set, shuffle and lock the order to start.
        </p>
      </header>

      {photos.length > 0 ? (
        <ul className="photo-review__list">
          {photos.map((photo, index) => (
            <li key={photo.pixabay_id} className="photo-review__item">
              <button
                type="button"
                className="photo-review__thumb"
                aria-label={`Preview saved photo ${index + 1} of ${photos.length}`}
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
          Edit selection
        </button>
        <button
          type="button"
          className="button button--primary"
          disabled={loading || photos.length === 0}
          onClick={onConfirm}
        >
          {loading ? 'Shuffling…' : 'Selection complete'}
        </button>
      </div>
    </section>
  )
}
