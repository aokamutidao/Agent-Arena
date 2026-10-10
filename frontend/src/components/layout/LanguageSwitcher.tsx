"use client";

import { Languages } from "lucide-react";
import { useI18n } from "@/lib/i18n";

export function LanguageSwitcher() {
  const { lang, setLang } = useI18n();
  const target = lang === "zh" ? "en" : "zh";
  return (
    <button
      type="button"
      onClick={() => setLang(target)}
      aria-label={lang === "zh" ? "Switch to English" : "切换为中文"}
      title={lang === "zh" ? "Switch to English" : "切换为中文"}
      className="inline-flex shrink-0 items-center gap-1.5 whitespace-nowrap rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm font-medium text-slate-700 transition-colors hover:border-emerald-300 hover:bg-emerald-50"
    >
      <Languages className="h-4 w-4" aria-hidden="true" />
      <span lang={target === "zh" ? "zh-CN" : "en"}>
        {lang === "zh" ? "English" : "中文"}
      </span>
    </button>
  );
}
