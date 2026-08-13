import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface LayoutState {
  sidebarCollapsed: boolean
  setSidebarCollapsed: (collapsed: boolean) => void
  toggleSidebar: () => void
}

export const useLayoutStore = create<LayoutState>()(
  persist(
    set => ({
      sidebarCollapsed: true,
      setSidebarCollapsed: collapsed => set({ sidebarCollapsed: collapsed }),
      toggleSidebar: () => set(state => ({ sidebarCollapsed: !state.sidebarCollapsed })),
    }),
    {
      name: 'admin-layout-storage',
      version: 1,
      partialize: state => ({ sidebarCollapsed: state.sidebarCollapsed }),
    }
  )
)
