self.addEventListener('fetch', (event) => {
  const request = event.request
  const url = new URL(request.url)

  if (request.method === 'POST' && url.pathname === '/') {
    event.respondWith(handleShareTarget(request))
    return
  }

  if (request.method === 'GET' && url.pathname === '/shared-image') {
    event.respondWith(readSharedImage())
  }
})

async function handleShareTarget(request) {
  const formData = await request.formData()
  const image = formData.get('image')
  if (image) {
    const cache = await caches.open('snaptask-share')
    await cache.put('/shared-image', new Response(image, {
      headers: { 'Content-Type': image.type || 'image/png' }
    }))
  }
  return Response.redirect('/?shared=1', 303)
}

async function readSharedImage() {
  const cache = await caches.open('snaptask-share')
  const response = await cache.match('/shared-image')
  if (!response) {
    return new Response('', { status: 404 })
  }
  await cache.delete('/shared-image')
  return response
}
