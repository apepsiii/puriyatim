// Service Worker - Puri Yatim PWA
// Versi cache — naikkan setiap kali ada update aset statis
const CACHE_VERSION = 'v1.0.0';
const CACHE_STATIC  = `puriyatim-static-${CACHE_VERSION}`;
const CACHE_PAGES   = `puriyatim-pages-${CACHE_VERSION}`;
const CACHE_IMAGES  = `puriyatim-images-${CACHE_VERSION}`;

// Aset statis yang dicache saat install
const PRECACHE_ASSETS = [
  '/',
  '/offline',
  '/static/css/app.css',
  '/static/js/app.js',
  '/static/manifest.json',
  '/static/images/icons/icon-192x192.png',
  '/static/images/icons/icon-512x512.png',
];

// Halaman yang di-cache saat dikunjungi (network-first)
const CACHEABLE_PAGES = [
  '/',
  '/tentang',
  '/berita',
  '/galeri',
  '/jumat-berkah',
  '/program-donasi',
  '/zakat',
  '/doa-harian',
  '/dzikir',
];

// ─── Install ──────────────────────────────────────────────────────────────────
self.addEventListener('install', event => {
  event.waitUntil(
    caches.open(CACHE_STATIC)
      .then(cache => cache.addAll(PRECACHE_ASSETS))
      .then(() => self.skipWaiting())
  );
});

// ─── Activate ─────────────────────────────────────────────────────────────────
self.addEventListener('activate', event => {
  const validCaches = [CACHE_STATIC, CACHE_PAGES, CACHE_IMAGES];
  event.waitUntil(
    caches.keys()
      .then(keys => Promise.all(
        keys
          .filter(key => !validCaches.includes(key))
          .map(key => caches.delete(key))
      ))
      .then(() => self.clients.claim())
  );
});

// ─── Fetch ────────────────────────────────────────────────────────────────────
self.addEventListener('fetch', event => {
  const { request } = event;
  const url = new URL(request.url);

  // Abaikan request non-GET dan request ke API / admin
  if (request.method !== 'GET') return;
  if (url.pathname.startsWith('/admin')) return;
  if (url.pathname.startsWith('/api/')) return;

  // Aset statis → Cache First
  if (isStaticAsset(url)) {
    event.respondWith(cacheFirstStrategy(request, CACHE_STATIC));
    return;
  }

  // Gambar → Cache First dengan fallback
  if (isImage(url)) {
    event.respondWith(cacheFirstStrategy(request, CACHE_IMAGES));
    return;
  }

  // Font & CDN eksternal → Cache First
  if (isExternalCDN(url)) {
    event.respondWith(cacheFirstStrategy(request, CACHE_STATIC));
    return;
  }

  // Halaman HTML → Network First dengan fallback ke cache
  if (isNavigationRequest(request)) {
    event.respondWith(networkFirstStrategy(request));
    return;
  }
});

// ─── Strategies ───────────────────────────────────────────────────────────────

/**
 * Cache First: cek cache dulu, jika tidak ada ambil dari network dan simpan.
 */
async function cacheFirstStrategy(request, cacheName) {
  const cached = await caches.match(request);
  if (cached) return cached;

  try {
    const response = await fetch(request);
    if (response.ok) {
      const cache = await caches.open(cacheName);
      cache.put(request, response.clone());
    }
    return response;
  } catch {
    return new Response('', { status: 503, statusText: 'Service Unavailable' });
  }
}

/**
 * Network First: coba network dulu, jika gagal fallback ke cache.
 * Jika keduanya gagal, tampilkan halaman offline.
 */
async function networkFirstStrategy(request) {
  try {
    const response = await fetch(request);
    if (response.ok) {
      const cache = await caches.open(CACHE_PAGES);
      cache.put(request, response.clone());
    }
    return response;
  } catch {
    const cached = await caches.match(request);
    if (cached) return cached;

    // Fallback ke halaman offline
    const offlinePage = await caches.match('/offline');
    return offlinePage || new Response(
      `<!DOCTYPE html><html lang="id"><head><meta charset="UTF-8">
      <title>Offline - Puri Yatim</title></head><body>
      <h1>Tidak ada koneksi internet</h1>
      <p>Silakan periksa koneksi Anda dan coba lagi.</p>
      </body></html>`,
      { headers: { 'Content-Type': 'text/html' } }
    );
  }
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

function isStaticAsset(url) {
  return url.pathname.startsWith('/static/css') ||
         url.pathname.startsWith('/static/js') ||
         url.pathname.startsWith('/static/manifest.json');
}

function isImage(url) {
  return url.pathname.startsWith('/static/images') ||
         url.pathname.startsWith('/static/uploads') ||
         /\.(png|jpg|jpeg|gif|webp|svg|ico)$/i.test(url.pathname);
}

function isExternalCDN(url) {
  return url.hostname.includes('cdnjs.cloudflare.com') ||
         url.hostname.includes('fonts.googleapis.com') ||
         url.hostname.includes('fonts.gstatic.com') ||
         url.hostname.includes('cdn.tailwindcss.com');
}

function isNavigationRequest(request) {
  return request.mode === 'navigate' ||
         request.headers.get('Accept')?.includes('text/html');
}

// ─── Push Notification (opsional, siap diaktifkan nanti) ──────────────────────
self.addEventListener('push', event => {
  if (!event.data) return;
  const data = event.data.json();
  event.waitUntil(
    self.registration.showNotification(data.title || 'Puri Yatim', {
      body:    data.body    || '',
      icon:    data.icon    || '/static/images/icons/icon-192x192.png',
      badge:   data.badge   || '/static/images/icons/icon-96x96.png',
      data:    data.url     ? { url: data.url } : {},
      actions: data.actions || [],
    })
  );
});

self.addEventListener('notificationclick', event => {
  event.notification.close();
  const url = event.notification.data?.url || '/';
  event.waitUntil(clients.openWindow(url));
});
