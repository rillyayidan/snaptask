const SHARE_CACHE = 'snaptask-share'
const SHARED_IMAGE_PATH = '/shared-image'

self.addEventListener('fetch', (event) => {
  const request = event.request
  const url = new URL(request.url)

  if (request.method === 'POST' && url.pathname === '/') {
    event.respondWith(handleShareTarget(request))
    return
  }

  if (request.method === 'GET' && url.pathname === SHARED_IMAGE_PATH) {
    event.respondWith(readSharedImage())
  }
})

async function handleShareTarget(request) {
  const formData = await request.formData()
  const image = firstSharedImage(formData.getAll('image'))
  if (!image) {
    return Response.redirect('/?share_error=missing-image', 303)
  }

  const cache = await caches.open(SHARE_CACHE)
  await cache.put(SHARED_IMAGE_PATH, new Response(image, {
    headers: {
      'Content-Type': image.type || 'application/octet-stream',
      'X-SnapTask-Filename': encodeURIComponent(image.name || 'shared-screenshot.png')
    }
  }))

  return Response.redirect('/?shared=1', 303)
}

async function readSharedImage() {
  const cache = await caches.open(SHARE_CACHE)
  const response = await cache.match(SHARED_IMAGE_PATH)
  if (!response) {
    return new Response('', { status: 404 })
  }
  await cache.delete(SHARED_IMAGE_PATH)
  return response
}

function firstSharedImage(values) {
  return values.find((value) => (
    value &&
    typeof value === 'object' &&
    typeof value.arrayBuffer === 'function' &&
    (!value.type || value.type.startsWith('image/'))
  )) ?? null
}
