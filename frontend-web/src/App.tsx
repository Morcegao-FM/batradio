import type { ReactNode } from 'react'
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import Login from './pages/Login'
import NowPlaying from './pages/NowPlaying'
import Playlists from './pages/Playlists'
import Settings from './pages/Settings'
import { useAuth } from './hooks/useAuth'

function RequireAuth({ children }: { children: ReactNode }) {
  const { email, loading } = useAuth()
  if (loading) {
    return null
  }
  if (!email) {
    return <Navigate to="/login" replace />
  }
  return children
}

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route
          path="/"
          element={
            <RequireAuth>
              <NowPlaying />
            </RequireAuth>
          }
        />
        <Route
          path="/playlists"
          element={
            <RequireAuth>
              <Playlists />
            </RequireAuth>
          }
        />
        <Route
          path="/config"
          element={
            <RequireAuth>
              <Settings />
            </RequireAuth>
          }
        />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  )
}
