import assert from 'node:assert/strict'
import { test } from 'node:test'
import { composeChatMessage, createChatImage, createChatQuote, splitChatMessage } from '../app/utils/chat.ts'
import { renderUserMarkdown } from '../app/utils/richText.ts'

test('quote replies survive a plain-text send/read round trip', () => {
  const quote = createChatQuote('Alice', 'Still seeding.\nFeel free to download.')
  const content = composeChatMessage(' Thank you! ', quote)
  assert.deepEqual(splitChatMessage(content), {
    quote: 'Alice: Still seeding. Feel free to download.',
    body: 'Thank you!'
  })
})

test('replying to a reply quotes only its own body', () => {
  assert.equal(createChatQuote('Bob', '> Alice: First message\n\nMy reply'), 'Bob: My reply')
  assert.equal(createChatQuote('Bob', '> Alice: First message'), '')
})

test('only leading quote lines are separated and HTML remains plain text', () => {
  assert.deepEqual(splitChatMessage('Hello\n> ordinary text\n<script>alert(1)</script>'), {
    quote: '', body: 'Hello\n> ordinary text\n<script>alert(1)</script>'
  })
  assert.deepEqual(splitChatMessage('> first\r\n> second\r\n\r\nreply'), {
    quote: 'first\nsecond', body: 'reply'
  })
})

test('Unicode excerpts stay within 150 characters and authors cannot add quote lines', () => {
  const quote = createChatQuote('Alice\n> Bob', '\u{1F600}'.repeat(160))
  assert.equal(Array.from(quote.slice('Alice > Bob: '.length)).length, 150)
  assert.ok(quote.endsWith('...'))
  assert.equal(quote.includes('\n'), false)
  assert.equal(composeChatMessage(' reply ', ''), 'reply')
})

test('combined content includes quote overhead in the 1000-character limit', () => {
  const quote = createChatQuote('Alice', 'Hello')
  const overhead = Array.from(composeChatMessage('', quote)).length
  assert.equal(Array.from(composeChatMessage('x'.repeat(1000 - overhead), quote)).length, 1000)
  assert.equal(Array.from(composeChatMessage('x'.repeat(1001 - overhead), quote)).length, 1001)
})

test('image insertion validates URLs and quotes use an image summary', () => {
  assert.equal(createChatImage('javascript:alert(1)'), '')
  assert.equal(createChatImage('data:image/png;base64,abc'), '')
  assert.equal(createChatImage('https://user:password@example.com/image.png'), '')
  assert.equal(createChatImage('not a url'), '')
  assert.equal(createChatImage(' https://example.com/image.png '), '![](<https://example.com/image.png>)')
  assert.equal(createChatQuote('Alice', createChatImage('https://example.com/image.png'), '[picture]'), 'Alice: [picture]')
  assert.equal(createChatQuote('Bob', '[img]https://example.com/image.png[/img] Nice!', '[picture]'), 'Bob: [picture] Nice!')
  assert.equal(createChatQuote('Bob', '![](<https://example.com/photo(1).png>) Nice!', '[picture]'), 'Bob: [picture] Nice!')
})

test('chat image rendering reuses safe Markdown and BBCode rendering', () => {
  for (const content of ['![](<https://example.com/photo.png>)', '[img]https://example.com/photo.png[/img]']) {
    const html = renderUserMarkdown(content)
    assert.match(html, /<img src="https:\/\/example.com\/photo.png"/)
    assert.match(html, /loading="lazy"/)
  }
  assert.doesNotMatch(renderUserMarkdown('<img src=x onerror=alert(1)>'), /<img/)
  assert.doesNotMatch(renderUserMarkdown('![](javascript:alert(1))'), /<img/)
})
