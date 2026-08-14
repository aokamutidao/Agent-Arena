"use client";

import { createContext, useContext, useEffect, useState, ReactNode } from "react";
import { useAccount } from "wagmi";

interface User {
  id: string;
  address: string;
  username: string;
  ac_balance: number;
  usdc_balance?: number;        // 链上 USDC 余额（整数，单位 USDC）
  usdc_balance_raw?: string;    // 链上 USDC 余额（原始值，wei 精度）
  ac_on_chain_balance?: number; // 链上 AC 余额
  created_at: string;
}

interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (address: string, message: string, signature: string) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
  updateBalance: (balance: number) => void;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const { address: connectedAddress } = useAccount();

  // 从 localStorage 恢复登录状态
  useEffect(() => {
    const savedToken = localStorage.getItem("auth_token");
    const savedUser = localStorage.getItem("auth_user");
    if (savedToken && savedUser) {
      setToken(savedToken);
      setUser(JSON.parse(savedUser));
    }
  }, []);

  // 监听钱包地址变化：如果当前连接的钱包与登录的钱包不同，清除登录状态
  useEffect(() => {
    if (!connectedAddress || !user) return;

    // 比较地址（都转为小写）
    const savedAddress = user.address.toLowerCase();
    const currentAddress = connectedAddress.toLowerCase();

    if (savedAddress !== currentAddress) {
      console.log("[Auth] 钱包地址变化，清除登录状态");
      console.log(`  旧地址: ${savedAddress}`);
      console.log(`  新地址: ${currentAddress}`);
      logout();
    }
  }, [connectedAddress, user]);

  const login = async (address: string, message: string, signature: string) => {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    const res = await fetch(`${apiUrl}/api/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ address, message, signature }),
    });

    if (!res.ok) {
      const err = await res.json();
      throw new Error(err.error || "login failed");
    }

    const data = await res.json();
    setToken(data.token);
    setUser(data.user);
    localStorage.setItem("auth_token", data.token);
    localStorage.setItem("auth_user", JSON.stringify(data.user));
  };

  const logout = () => {
    setToken(null);
    setUser(null);
    localStorage.removeItem("auth_token");
    localStorage.removeItem("auth_user");
  };

  const updateBalance = (balance: number) => {
    if (user) {
      const updatedUser = { ...user, ac_balance: balance, ac_on_chain_balance: balance };
      setUser(updatedUser);
      localStorage.setItem("auth_user", JSON.stringify(updatedUser));
    }
  };

  const refreshUser = async () => {
    const savedToken = localStorage.getItem("auth_token");
    if (!savedToken) return;
    try {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
      const res = await fetch(`${apiUrl}/api/auth/profile`, {
        headers: { Authorization: `Bearer ${savedToken}` },
      });
      if (res.ok) {
        const data = await res.json();
        // 后端返回的 profile 结构可能包含链上字段；仅更新核心用户信息
        const updatedUser: User = {
          id: data.id,
          address: data.address,
          username: data.username,
          ac_balance: data.ac_balance,
          created_at: data.created_at,
        };
        if (typeof data.usdc_balance === "number") {
          updatedUser.usdc_balance = data.usdc_balance;
        }
        if (typeof data.usdc_balance_raw === "string") {
          updatedUser.usdc_balance_raw = data.usdc_balance_raw;
        }
        if (typeof data.ac_on_chain_balance === "number") {
          updatedUser.ac_on_chain_balance = data.ac_on_chain_balance;
        }
        setUser(updatedUser);
        localStorage.setItem("auth_user", JSON.stringify(updatedUser));
      }
    } catch (err) {
      console.error("refreshUser failed:", err);
    }
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        login,
        logout,
        isAuthenticated: !!token,
        updateBalance,
        refreshUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
