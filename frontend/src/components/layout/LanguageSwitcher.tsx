"use client";

import { useI18n } from "@/lib/i18n";
import { Button } from "@/components/ui/button";

export function LanguageSwitcher() {
  const { lang, setLang } = useI18n();

  return (
    <Button
      variant="outline"
      size="sm"
      onClick={() => setLang(lang === "en" ? "zh" : "en")}
      className="min-w-[60px]"
    >
      {lang === "en" ? "中文" : "EN"}
    </Button>
  );
}
