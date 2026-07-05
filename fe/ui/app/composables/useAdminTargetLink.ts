export interface AdminTargetLinkTarget {
  type: string
  id: number
  parent_type?: string
  parent_id?: number
  status?: string
}

export function useAdminTargetLink() {
  function getTargetPath(target: AdminTargetLinkTarget) {
    if (target.status !== 'normal') return ''
    if (target.type === 'catalog_torrent') return `/catalog/torrents/${target.id}`
    if (target.type === 'catalog_comment' && target.parent_type === 'catalog_torrent' && target.parent_id) {
      return `/catalog/torrents/${target.parent_id}#comment-${target.id}`
    }
    if (target.type === 'catalog_subtitle' && target.parent_type === 'catalog_torrent' && target.parent_id) {
      return `/catalog/torrents/${target.parent_id}`
    }
    if (target.type === 'forum_topic') return `/forum/topics/${target.id}`
    if (target.type === 'forum_reply' && target.parent_type === 'forum_topic' && target.parent_id) {
      return `/forum/topics/${target.parent_id}#reply-${target.id}`
    }
    return ''
  }

  return {
    getTargetPath
  }
}
