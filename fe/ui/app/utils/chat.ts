import MarkdownIt from 'markdown-it'

const quoteMarkdown = new MarkdownIt({ html: false })

export function splitChatMessage(content: string) {
  const lines = content.replace(/\r\n?/g, '\n').split('\n')
  let end = 0
  while (end < lines.length && lines[end]?.startsWith('>')) end += 1
  if (end === 0) return { quote: '', body: content }

  return {
    quote: lines.slice(0, end).map(line => line.replace(/^> ?/, '')).join('\n'),
    body: lines.slice(end).join('\n').replace(/^\n+/, '')
  }
}

export function createChatQuote(author: string, content: string, imageLabel = '[image]') {
  const source = splitChatMessage(content).body.replace(/\[img\][\s\S]*?\[\/img\]/gi, () => imageLabel)
  const body = quoteMarkdown.parseInline(source, {})
    .flatMap(token => token.children || [])
    .map(token => token.type === 'image' ? imageLabel : ['softbreak', 'hardbreak'].includes(token.type) ? ' ' : token.content)
    .join('').replace(/\s+/g, ' ').trim()
  if (!body) return ''
  const characters = Array.from(body)
  const excerpt = characters.length > 150 ? characters.slice(0, 147).join('') + '...' : body
  return `${author.replace(/\s+/g, ' ').trim()}: ${excerpt}`
}

export function composeChatMessage(body: string, quote: string) {
  return quote ? `> ${quote}\n\n${body.trim()}` : body.trim()
}

export function createChatImage(url: string) {
  const value = url.trim()
  try {
    const parsed = new URL(value)
    if (!['http:', 'https:'].includes(parsed.protocol) || parsed.username || parsed.password) return ''
    return `![](<${parsed.href.replace(/</g, '%3C').replace(/>/g, '%3E')}>)`
  } catch {
    return ''
  }
}
