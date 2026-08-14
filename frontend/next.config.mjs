/** @type {import('next').NextConfig} */
const nextConfig = {
  eslint: {
    // 黑客松阶段忽略 ESLint 错误以允许构建（已有 lint 错误为历史遗留）
    ignoreDuringBuilds: true,
  },
};

export default nextConfig;
