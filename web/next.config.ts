import type { NextConfig } from "next";

const imageOrigins = (process.env.IMAGE_REMOTE_ORIGINS || '')
  .split(',').map((origin) => origin.trim()).filter(Boolean)
  .map((origin) => {
    const url = new URL(origin);
    if (!['http:', 'https:'].includes(url.protocol) || url.username || url.password ||
        url.pathname !== '/' || url.search || url.hash || url.hostname.includes('*')) {
      throw new Error('IMAGE_REMOTE_ORIGINS requires exact HTTP(S) origins');
    }
    return new URL('/**', url.origin);
  });

const nextConfig: NextConfig = {
  output: 'standalone',
  reactStrictMode: true,
  poweredByHeader: false,
  images: {
    remotePatterns: [
      ...imageOrigins,
      { protocol: 'http', hostname: 'localhost' },
    ],
  },
  compiler: {
    removeConsole: process.env.NODE_ENV === 'production',
  },
};

export default nextConfig;
