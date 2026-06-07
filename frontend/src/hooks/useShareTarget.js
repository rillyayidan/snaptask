import { useEffect, useState } from 'react'

export function useShareTarget() {
  const [share, setShare] = useState({ image: null, error: '' })

  useEffect(() => {
    async function readServiceWorkerShare() {
      const params = new URLSearchParams(window.location.search)
      const shareError = params.get('share_error')
      if (shareError) {
        setShare({ image: null, error: shareErrorMessage(shareError) })
        window.history.replaceState({}, '', '/')
        return
      }
      if (!params.has('shared')) return

      try {
        const response = await fetch('/shared-image', { cache: 'no-store' })
        if (!response.ok) {
          setShare({ image: null, error: 'No shared screenshot was found. Try sharing it again.' })
          return
        }

        const blob = await response.blob()
        if (blob.type && !blob.type.startsWith('image/')) {
          setShare({ image: null, error: 'The shared file was not an image.' })
          return
        }

        setShare({
          image: new File([blob], sharedFileName(response.headers.get('X-SnapTask-Filename')), {
            type: blob.type || 'image/png'
          }),
          error: ''
        })
      } catch (err) {
        setShare({ image: null, error: err.message || 'Could not read the shared screenshot.' })
      } finally {
        window.history.replaceState({}, '', '/')
      }
    }

    async function readLaunchQueue() {
      if (!('launchQueue' in window)) return
      window.launchQueue.setConsumer(async (launchParams) => {
        const fileHandle = launchParams.files?.[0]
        if (!fileHandle) return
        const file = await fileHandle.getFile()
        if (file.type && !file.type.startsWith('image/')) {
          setShare({ image: null, error: 'The shared file was not an image.' })
          return
        }
        setShare({ image: file, error: '' })
      })
    }

    readServiceWorkerShare()
    readLaunchQueue()
  }, [])

  return share
}

function sharedFileName(headerValue) {
  if (!headerValue) return 'shared-screenshot.png'

  try {
    return decodeURIComponent(headerValue).replace(/[/\\]/g, '-')
  } catch {
    return 'shared-screenshot.png'
  }
}

function shareErrorMessage(value) {
  if (value === 'missing-image') {
    return 'The shared content did not include a screenshot image.'
  }
  if (value === 'too-large') {
    return 'The shared screenshot is too large. Use an image under 10 MB.'
  }
  return 'Could not import the shared screenshot.'
}
