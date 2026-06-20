import { useCallback, useEffect, useRef, useState } from 'react'
import { ApiError } from '../../api/client'
import {
  hitToPhoto,
  searchPixabay,
  type PixabaySearchHit,
} from '../../api/pixabayApi'
import type { Photo } from '../../types/lobby'

const RECOMMENDED_COUNT = 5

interface PixabaySearchProps {
  lobbyId: string
  selectedPhotos: Photo[]
  onSelectionChange: (photos: Photo[]) => void
  readOnly?: boolean
}

export function PixabaySearchPanel({
  lobbyId,
  selectedPhotos,
  onSelectionChange,
  readOnly = false,
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
        })
        setHits(result.hits)
        setPage(nextPage)
        setTotalPages(Math.max(1, Math.ceil(result.total_hits / 20)))
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
    <section className="pixabay-search" aria-label="PixaBay search">
      {!readOnly ? (
        <form
          className="pixabay-search__form"
          onSubmit={(event) => {
            event.preventDefault()
            void runSearch(1)
          }}
        >
          <label className="pixabay-search__field">
            <span className="pixabay-search__label">Search PixaBay</span>
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
            className="button button--primary"
            disabled={loading}
          >
            {loading ? 'Searching…' : 'Search'}
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
              className={`pixabay-search__card${
                selected ? ' pixabay-search__card--selected' : ''
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
                width={hit.width}
                height={hit.height}
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
            className="button button--secondary"
            disabled={loading || page <= 1}
            onClick={() => void runSearch(page - 1)}
          >
            Previous
          </button>
          <span>
            Page {page} of {totalPages}
          </span>
          <button
            type="button"
            className="button button--secondary"
            disabled={loading || page >= totalPages}
            onClick={() => void runSearch(page + 1)}
          >
            Next
          </button>
        </div>
      ) : null}

      <p className="pixabay-search__attribution">
        Images from{' '}
        <a href="https://pixabay.com" target="_blank" rel="noreferrer">
          PixaBay
        </a>
      </p>
    </section>
  )
}
