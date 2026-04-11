import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [react(), tailwindcss()],
    server: {
        host: true,
        proxy: {
            // Routes for coffee API → port 8081
            '/api/v1/coffee': {
                target: 'http://localhost:8081',
                changeOrigin: true,
                rewrite: (path) => path.replace(/^\/api\/v1\/coffee/, '/api/v1/coffee'),
            },

            // Routes for auth API → port 8080
            '/api/v1/auth': {
                target: 'http://localhost:8080',
                changeOrigin: true,
                rewrite: (path) => path.replace(/^\/api\/v1\/auth/, '/api/v1/auth'),
            },
        },
    },
})

