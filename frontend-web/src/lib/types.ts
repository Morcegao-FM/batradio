// Tipos espelhando o JSON do gateway Go.

export interface Song {
  file: string
  artist: string
  title: string
  album?: string
  genre?: string
  time: number
  pos: number
  id: number
  /** Vindos do catálogo do site pelo gateway; ausentes quando ele não conhece a faixa. */
  imageUrl?: string
  displayName?: string
  year?: number
  kind?: string
}

export interface QueueItem extends Song {
  nextPresentation: string
}

export interface Status {
  state: 'play' | 'pause' | 'stop'
  song: number
  elapsed: number
  repeat: boolean
  random: boolean
  crossfade: boolean
  playlistLength: number
  current?: Song
  next?: Song
}

export interface Page<T> {
  items: T[]
  total: number
  offset: number
  limit: number
}

export interface QueuePage extends Page<QueueItem> {
  currentPos: number
  version: number
}

export interface PlaylistInfo {
  name: string
  lastModified?: string
}

export interface ServerInfo {
  nodeOk: boolean
  nodeHost: string
  libraryCount: number
  libraryRefreshedAt?: string
  streamUrl: string
}
