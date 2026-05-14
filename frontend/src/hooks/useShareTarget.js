import { useEffect, useState } from 'react'

export function useShareTarget() {
  const [image, setImage] = useState(null)

  useEffect(() => {
    async function readServiceWorkerShare() {
      const params = new URLSearchParams(window.location.search)
      if (!params.has('shared')) return
      const response = await fetch('/shared-image', { cache: 'no-store' })
      if (!response.ok) return
      const blob = await response.blob()
      setImage(new File([blob], 'shared-screenshot.png', { type: blob.type || 'image/png' }))
      window.history.replaceState({}, '', '/')
    }

    async function readLaunchQueue() {
      if (!('launchQueue' in window)) return
      window.launchQueue.setConsumer(async (launchParams) => {
        const fileHandle = launchParams.files?.[0]
        if (!fileHandle) return
        const file = await fileHandle.getFile()
        setImage(file)
      })
    }

    readServiceWorkerShare()
    readLaunchQueue()
  }, [])

  return image
}
