/** @type {import('next').NextConfig} */
const nextConfig = {
  // Produce a self-contained server build for a small Docker image / Railway.
  output: "standalone",
  reactStrictMode: true,
};

export default nextConfig;
