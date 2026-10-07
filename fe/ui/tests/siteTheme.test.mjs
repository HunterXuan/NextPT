import assert from 'node:assert/strict'
import { test } from 'node:test'
import { runInNewContext } from 'node:vm'
import { resolveSiteTheme, siteThemeBootstrap, siteThemeStyles, siteThemes } from '../app/utils/siteTheme.ts'

function bootstrap(defaultTheme, saved, storageUnavailable = false) {
  const document = { documentElement: { dataset: {} } }
  runInNewContext(siteThemeBootstrap(defaultTheme), {
    document,
    localStorage: {
      getItem() {
        if (storageUnavailable) throw new Error('Storage unavailable')
        return saved
      }
    }
  })
  return document.documentElement.dataset.siteTheme
}

test('all eight themes have light and dark styles and restore before hydration', () => {
  assert.equal(siteThemes.length, 8)
  assert.equal(new Set(siteThemes.map(theme => theme.id)).size, 8)
  for (const { id } of siteThemes) {
    assert.equal(resolveSiteTheme(id), id)
    assert.equal(bootstrap('sky', id), id)
    assert.ok(siteThemeStyles.includes(`html[data-site-theme="${id}"]`))
    assert.ok(siteThemeStyles.includes(`html.dark[data-site-theme="${id}"]`))
  }
})

test('invalid defaults and saved values fall back without injecting markup', () => {
  for (const value of [null, undefined, '', 'unknown', '</script>']) {
    assert.equal(resolveSiteTheme(value), 'sky')
    assert.equal(bootstrap('teal', value), 'teal')
  }
})

test('no override follows the runtime default; storage denial is harmless', () => {
  assert.equal(bootstrap('violet', null), 'violet')
  assert.equal(bootstrap('rose', 'green', true), 'rose')
})
