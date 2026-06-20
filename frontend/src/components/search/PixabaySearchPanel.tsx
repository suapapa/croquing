import { useCallback, useEffect, useRef, useState, type ReactNode } from 'react'
import { ApiError } from '../../api/client'
import {
  hitToPhoto,
  searchPixabay,
  type PixabaySearchHit,
} from '../../api/pixabayApi'
import type { Photo } from '../../types/lobby'
import {
  IconChevronLeft,
  IconChevronRight,
  IconSearch,
  IconSpinner,
} from '../ui/Icons'

const RECOMMENDED_COUNT = 5
const PIXABAY_PER_PAGE = 24

interface PixabaySearchProps {
  lobbyId: string
  selectedPhotos: Photo[]
  onSelectionChange: (photos: Photo[]) => void
  readOnly?: boolean
  footerStart?: ReactNode
}

export function PixabaySearchPanel({
  lobbyId,
  selectedPhotos,
  onSelectionChange,
  readOnly = false,
  footerStart,
}: PixabaySearchProps) {
  const [query, setQuery] = useState('figure drawing reference')
  const [order, setOrder] = useState<'popular' | 'latest'>('popular')
  const [page, setPage] = useState(1)
  const [hits, setHits] = useState<PixabaySearchHit[]>([])
  const [totalPages, setTotalPages] = useState(1)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const initialSearchDone = useRef(false)

  const selectedIds = new Set(selectedPhotos.map((photo) => photo.pixabay_id))

  const runSearch = useCallback(
    async (nextPage: number) => {
      const trimmed = query.trim()
      if (!trimmed) {
        setError('Enter a search term')
        return
      }

      setLoading(true)
      setError(null)

      try {
        const result = await searchPixabay({
          lobbyId,
          query: trimmed,
          order,
          page: nextPage,
          perPage: PIXABAY_PER_PAGE,
        })
        setHits(result.hits)
        setPage(nextPage)
        setTotalPages(
          Math.max(1, Math.ceil(result.total_hits / PIXABAY_PER_PAGE)),
        )
      } catch (err) {
        const message =
          err instanceof ApiError ? err.message : 'Search failed'
        setError(message)
      } finally {
        setLoading(false)
      }
    },
    [lobbyId, order, query],
  )

  useEffect(() => {
    if (!readOnly && !initialSearchDone.current) {
      initialSearchDone.current = true
      void runSearch(1)
    }
  }, [readOnly, runSearch])

  function togglePhoto(hit: PixabaySearchHit) {
    if (readOnly) {
      return
    }

    if (selectedIds.has(hit.pixabay_id)) {
      onSelectionChange(
        selectedPhotos.filter((photo) => photo.pixabay_id !== hit.pixabay_id),
      )
      return
    }

    onSelectionChange([...selectedPhotos, hitToPhoto(hit)])
  }

  return (
    <section className="pixabay-search" aria-label="Pixabay search">
      {!readOnly ? (
        <form
          className="pixabay-search__form"
          onSubmit={(event) => {
            event.preventDefault()
            void runSearch(1)
          }}
        >
          <label className="pixabay-search__field">
            <span className="pixabay-search__label">Search Pixabay</span>
            <input
              type="search"
              value={query}
              onChange={(event) => setQuery(event.target.value)}
              placeholder="e.g. portrait, anatomy, gesture"
              autoComplete="off"
            />
          </label>

          <label className="pixabay-search__field pixabay-search__field--compact">
            <span className="pixabay-search__label">Sort</span>
            <select
              value={order}
              onChange={(event) =>
                setOrder(event.target.value as 'popular' | 'latest')
              }
            >
              <option value="popular">Popular</option>
              <option value="latest">Latest</option>
            </select>
          </label>

          <button
            type="submit"
            className="button button--primary button--icon-only"
            disabled={loading}
            aria-label={loading ? 'Searching' : 'Search'}
            title={loading ? 'Searching' : 'Search'}
          >
            {loading ? (
              <IconSpinner className="button__spinner" />
            ) : (
              <IconSearch className="button__icon" />
            )}
          </button>
        </form>
      ) : null}

      <p className="pixabay-search__hint">
        {selectedPhotos.length} selected · {RECOMMENDED_COUNT} recommended
      </p>

      {error ? (
        <p className="pixabay-search__error" role="alert">
          {error}
        </p>
      ) : null}

      <div className="pixabay-search__grid" role="list">
        {hits.map((hit) => {
          const selected = selectedIds.has(hit.pixabay_id)
          return (
            <button
              key={hit.pixabay_id}
              type="button"
              role="listitem"
              className={`pixabay-search__card${selected ? ' pixabay-search__card--selected' : ''
                }`}
              onClick={() => togglePhoto(hit)}
              disabled={readOnly}
              aria-pressed={selected}
              aria-label={`${selected ? 'Deselect' : 'Select'} photo ${hit.pixabay_id}`}
            >
              <img
                src={hit.preview_url}
                alt=""
                loading="lazy"
              />
              {selected ? (
                <span className="pixabay-search__check" aria-hidden="true">
                  ✓
                </span>
              ) : null}
            </button>
          )
        })}
      </div>

      {!readOnly && hits.length > 0 ? (
        <div className="pixabay-search__pagination">
          <button
            type="button"
            className="button button--secondary button--icon-only"
            disabled={loading || page <= 1}
            onClick={() => void runSearch(page - 1)}
            aria-label="Previous page"
            title="Previous page"
          >
            <IconChevronLeft className="button__icon" />
          </button>
          <span>
            Page {page} of {totalPages}
          </span>
          <button
            type="button"
            className="button button--secondary button--icon-only"
            disabled={loading || page >= totalPages}
            onClick={() => void runSearch(page + 1)}
            aria-label="Next page"
            title="Next page"
          >
            <IconChevronRight className="button__icon" />
          </button>
        </div>
      ) : null}

      <div
        className={`pixabay-search__footer${footerStart ? ' pixabay-search__footer--split' : ''
          }`}
      >
        {footerStart}
        <p className="pixabay-search__attribution">
          Images from{' '}
          <a href="https://pixabay.com" target="_blank" rel="noreferrer">
            Pixabay
          </a>
        </p>
      </div>
    </section>
  )
}
