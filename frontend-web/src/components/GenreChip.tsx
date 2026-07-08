import styles from './GenreChip.module.css'

export default function GenreChip({ genre }: { genre?: string }) {
  if (!genre) return null
  const vinheta = genre.toUpperCase() === 'VINHETA'
  return (
    <span className={`${styles.chip} ${vinheta ? styles.vinheta : ''}`}>{genre}</span>
  )
}
