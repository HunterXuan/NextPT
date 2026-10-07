import { readFile } from 'node:fs/promises'
import { resolve } from 'node:path'

const imageTypes: Record<string, string> = {
  svg: 'image/svg+xml',
  png: 'image/png',
  webp: 'image/webp',
  ico: 'image/x-icon'
}

export default defineEventHandler(async (event) => {
  const file = getRouterParam(event, 'file') || ''
  const match = /^[a-zA-Z0-9_-]+\.(svg|png|webp|ico)$/.exec(file)
  if (!match) throw createError({ statusCode: 404 })

  // Mounted files are not part of Nitro's build-time static asset manifest.
  const directory = import.meta.dev ? 'public/branding' : '.output/public/branding'
  let image: Buffer
  try {
    image = await readFile(resolve(directory, file))
  } catch {
    throw createError({ statusCode: 404 })
  }

  setResponseHeader(event, 'Content-Type', imageTypes[match[1]!]!)
  setResponseHeader(event, 'Cache-Control', 'no-cache')
  setResponseHeader(event, 'X-Content-Type-Options', 'nosniff')
  return image
})
