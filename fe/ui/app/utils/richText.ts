import MarkdownIt from 'markdown-it'
import { parse as parseBBCode } from '@bbob/parser'
import type { RenderRule } from 'markdown-it/lib/renderer.mjs'

type BBCodeNode = string | {
  tag?: string
  attrs?: Record<string, string>
  content?: BBCodeNode[]
}

const markdown = new MarkdownIt({
  html: false,
  linkify: true,
  breaks: true
})

const supportedBBCodeTags = ['b', 'strong', 'i', 'em', 'u', 's', 'strike', 'del', 'color', 'size', 'url', 'img', 'quote', 'code', 'list', '*']

const defaultLinkOpenRule = markdown.renderer.rules.link_open
const defaultImageRule = markdown.renderer.rules.image

markdown.renderer.rules.link_open = (tokens, idx, options, env, self) => {
  tokens[idx].attrSet('target', '_blank')
  tokens[idx].attrSet('rel', 'noopener noreferrer')
  return renderToken(defaultLinkOpenRule, tokens, idx, options, env, self)
}

markdown.renderer.rules.image = (tokens, idx, options, env, self) => {
  tokens[idx].attrSet('loading', 'lazy')
  return renderToken(defaultImageRule, tokens, idx, options, env, self)
}

export function renderUserMarkdown(input: string) {
  const source = normalizeBBCodeToMarkdown(input)
  return markdown.render(source)
}

export function normalizeBBCodeToMarkdown(input: string) {
  const source = String(input || '')
  if (!hasSupportedBBCode(source)) return source

  try {
    return renderBBCodeNodesAsMarkdown(parseBBCode(source, { onlyAllowTags: supportedBBCodeTags }) as BBCodeNode[])
  } catch {
    return source
  }
}

function hasSupportedBBCode(input: string) {
  return /\[\/?(?:b|strong|i|em|u|s|strike|del|color|size|url|img|quote|code|list|\*)(?:[\]=\s]|\/?\])/i.test(input)
}

function renderBBCodeNodesAsMarkdown(nodes: BBCodeNode[]) {
  return nodes.map(renderBBCodeNodeAsMarkdown).join('')
}

function renderBBCodeNodeAsMarkdown(node: BBCodeNode): string {
  if (typeof node === 'string') return node

  const tag = String(node.tag || '').toLowerCase()
  const content = renderBBCodeNodesAsMarkdown(node.content || [])

  switch (tag) {
    case 'b':
    case 'strong':
      return `**${content}**`
    case 'i':
    case 'em':
      return `*${content}*`
    case 'u':
      return content
    case 's':
    case 'strike':
    case 'del':
      return `~~${content}~~`
    case 'color':
    case 'size':
      return content
    case 'url':
      return renderBBCodeUrl(node, content)
    case 'img':
      return renderBBCodeImage(content)
    case 'quote':
      return renderBBCodeQuote(node, content)
    case 'code':
      return `\n\`\`\`\n${content.trim()}\n\`\`\`\n`
    case 'list':
      return renderBBCodeList(node.content || [])
    case '*':
      return ''
    default:
      return content
  }
}

function renderBBCodeUrl(node: Exclude<BBCodeNode, string>, content: string) {
  const href = getBBCodePrimaryAttr(node) || content.trim()
  if (!isSafeMarkdownUrl(href)) return content

  const label = content.trim() || href
  return `[${escapeMarkdownLinkText(label)}](${href})`
}

function renderBBCodeImage(content: string) {
  const src = content.trim()
  if (!isSafeMarkdownUrl(src)) return content
  return `\n![](${src})\n`
}

function renderBBCodeQuote(node: Exclude<BBCodeNode, string>, content: string) {
  const author = getBBCodePrimaryAttr(node)
  const quote = content
    .trim()
    .split('\n')
    .map((line) => `> ${line}`)
    .join('\n')
  return author ? `\n> ${author}:\n${quote ? `${quote}\n` : ''}\n` : `\n${quote}\n`
}

function renderBBCodeList(nodes: BBCodeNode[]) {
  const items: string[] = []
  let current = ''

  for (const node of nodes) {
    if (typeof node !== 'string' && String(node.tag || '').toLowerCase() === '*') {
      if (current.trim()) items.push(current.trim())
      current = ''
      continue
    }
    current += renderBBCodeNodeAsMarkdown(node)
  }
  if (current.trim()) items.push(current.trim())

  return items.length > 0 ? `\n${items.map((item) => `- ${item}`).join('\n')}\n` : ''
}

function getBBCodePrimaryAttr(node: Exclude<BBCodeNode, string>) {
  const attrs = node.attrs || {}
  return attrs.url || attrs.href || attrs.default || Object.values(attrs)[0] || ''
}

function isSafeMarkdownUrl(value: string) {
  try {
    const url = new URL(value)
    return url.protocol === 'http:' || url.protocol === 'https:'
  } catch {
    return false
  }
}

function escapeMarkdownLinkText(value: string) {
  return value.replace(/([\\\]])/g, '\\$1')
}

function renderToken(rule: RenderRule | undefined, ...args: Parameters<RenderRule>) {
  return rule ? rule(...args) : args[4].renderToken(args[0], args[1], args[2])
}
