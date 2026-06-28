document.addEventListener('DOMContentLoaded', () => {
  const tablist = document.querySelector('[role="tablist"]')
  const tabs = tablist
    ? Array.from(tablist.querySelectorAll('[role="tab"]'))
    : []
  const panels = Array.from(document.querySelectorAll('[role="tabpanel"]'))
  const lightbox = document.getElementById('lightbox')
  const lightboxContent = lightbox?.querySelector('.lightbox-content')
  const lightboxClose = lightbox?.querySelector('.lightbox-close')
  let lastFocusedElement = null
  let lightboxImage = null

  function activateTab(nextTab) {
    const targetId = nextTab.getAttribute('aria-controls')
    if (!targetId) return

    tabs.forEach((tab) => {
      const selected = tab === nextTab
      tab.classList.toggle('active', selected)
      tab.setAttribute('aria-selected', selected ? 'true' : 'false')
      tab.tabIndex = selected ? 0 : -1
    })

    panels.forEach((panel) => {
      const active = panel.id === targetId
      panel.classList.toggle('active', active)
      panel.hidden = !active
    })

    nextTab.focus()
  }

  tabs.forEach((tab, index) => {
    tab.tabIndex = tab.getAttribute('aria-selected') === 'true' ? 0 : -1

    tab.addEventListener('click', () => {
      activateTab(tab)
    })

    tab.addEventListener('keydown', (event) => {
      let nextIndex = index

      if (event.key === 'ArrowRight' || event.key === 'ArrowDown') {
        event.preventDefault()
        nextIndex = (index + 1) % tabs.length
      } else if (event.key === 'ArrowLeft' || event.key === 'ArrowUp') {
        event.preventDefault()
        nextIndex = (index - 1 + tabs.length) % tabs.length
      } else if (event.key === 'Home') {
        event.preventDefault()
        nextIndex = 0
      } else if (event.key === 'End') {
        event.preventDefault()
        nextIndex = tabs.length - 1
      } else {
        return
      }

      activateTab(tabs[nextIndex])
    })
  })

  document.querySelectorAll('.mockup-trigger').forEach((trigger) => {
    trigger.addEventListener('click', () => {
      if (!lightbox || !lightboxContent) return

      const img = trigger.querySelector('img')
      if (!img) return

      lastFocusedElement = trigger
      if (!lightboxImage) {
        lightboxImage = document.createElement('img')
        lightboxImage.className = 'lightbox-image'
        lightboxContent.appendChild(lightboxImage)
      }
      lightboxImage.src = img.currentSrc || img.src
      lightboxImage.alt = img.alt
      lightbox.showModal()
      lightboxClose?.focus()
    })
  })

  function closeLightbox() {
    if (!lightbox?.open) return
    lightbox.close()
    if (lightboxImage) {
      lightboxImage.remove()
      lightboxImage = null
    }
    lastFocusedElement?.focus()
  }

  lightboxClose?.addEventListener('click', closeLightbox)

  lightbox?.addEventListener('click', (event) => {
    const rect = lightbox.getBoundingClientRect()
    const clickedBackdrop =
      event.clientX < rect.left ||
      event.clientX > rect.right ||
      event.clientY < rect.top ||
      event.clientY > rect.bottom
    if (clickedBackdrop) {
      closeLightbox()
    }
  })

  lightbox?.addEventListener('cancel', (event) => {
    event.preventDefault()
    closeLightbox()
  })
})
