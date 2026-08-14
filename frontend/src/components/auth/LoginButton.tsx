"use client";

import { useState } from "react";
import { useAccount, useSignMessage } from "wagmi";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { useAuth } from "@/lib/auth";

export function LoginButton() {
  const { address, isConnected } = useAccount();
  const { signMessageAsync } = useSignMessage();
  const { login, logout, user, isAuthenticated } = useAuth();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleLogin = async () => {
    if (!address) return;

    setLoading(true);
    setError("");

    try {
      const timestamp = Date.now();
      const message = `Login to Agent Arena: ${address} at ${timestamp}`;
      const signature = await signMessageAsync({ message });

      await login(address, message, signature);
    } catch (err: any) {
      setError(err.message || "登录失败");
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    logout();
  };

  if (!isConnected) {
    return null; // RainbowKit 的 ConnectButton 会处理
  }

  if (isAuthenticated && user) {
    return (
      <div className="flex items-center gap-2">
        <span className="text-sm text-muted-foreground">
          {user.username}
        </span>
        <Button variant="outline" size="sm" onClick={handleLogout}>
          登出
        </Button>
      </div>
    );
  }

  return (
    <div>
      <Button onClick={handleLogin} disabled={loading} size="sm">
        {loading ? "签名中..." : "登录"}
      </Button>
      {error && <p className="text-xs text-red-500 mt-1">{error}</p>}
    </div>
  );
}

// 完整的登录页面组件
export function LoginPage() {
  const { address, isConnected } = useAccount();
  const { signMessageAsync } = useSignMessage();
  const { login, isAuthenticated } = useAuth();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleLogin = async () => {
    if (!address) return;

    setLoading(true);
    setError("");

    try {
      const timestamp = Date.now();
      const message = `Login to Agent Arena: ${address} at ${timestamp}`;
      const signature = await signMessageAsync({ message });

      await login(address, message, signature);
    } catch (err: any) {
      setError(err.message || "登录失败");
    } finally {
      setLoading(false);
    }
  };

  if (isAuthenticated) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>✅ 已登录</CardTitle>
          <CardDescription>你已经成功登录 Agent Arena</CardDescription>
        </CardHeader>
      </Card>
    );
  }

  if (!isConnected) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>请先连接钱包</CardTitle>
          <CardDescription>点击右上角的「连接钱包」按钮</CardDescription>
        </CardHeader>
      </Card>
    );
  }

  return (
    <Card className="max-w-md mx-auto mt-20">
      <CardHeader>
        <CardTitle>登录 Agent Arena</CardTitle>
        <CardDescription>
          使用钱包签名登录，无需密码
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="text-sm text-muted-foreground">
          钱包地址: <span className="font-mono">{address}</span>
        </div>
        <Button
          onClick={handleLogin}
          disabled={loading}
          className="w-full"
        >
          {loading ? "签名中..." : "🔐 签名登录"}
        </Button>
        {error && (
          <p className="text-sm text-red-500 text-center">{error}</p>
        )}
        <p className="text-xs text-muted-foreground text-center">
          签名仅用于验证身份，不会产生任何费用
        </p>
      </CardContent>
    </Card>
  );
}
