"use client";
import type { Tournament } from "@/lib/course";
import { useI18n } from "@/lib/i18n";
import {
  getTournamentPresentation,
  type PresentationTone,
} from "@/lib/course-presentation";

const TONE_STYLES: Record<PresentationTone, string> = {
  neutral:
    "border-slate-200 bg-slate-100 text-slate-700 dark:border-slate-700 dark:bg-slate-800 dark:text-slate-200",
  success:
    "border-emerald-200 bg-emerald-50 text-emerald-800 dark:border-emerald-800 dark:bg-emerald-950 dark:text-emerald-200",
  warning:
    "border-amber-200 bg-amber-50 text-amber-800 dark:border-amber-800 dark:bg-amber-950 dark:text-amber-200",
  danger:
    "border-red-200 bg-red-50 text-red-800 dark:border-red-800 dark:bg-red-950 dark:text-red-200",
  info: "border-sky-200 bg-sky-50 text-sky-800 dark:border-sky-800 dark:bg-sky-950 dark:text-sky-200",
};

interface TournamentStatusBadgeProps {
  tournament: Tournament;
  now: number;
}

/** A text-labelled badge; color supplements the business state instead of replacing it. */
export function TournamentStatusBadge({
  tournament,
  now,
}: TournamentStatusBadgeProps) {
  const { lang } = useI18n();
  const { label, tone, hint } = getTournamentPresentation(
    tournament,
    now,
    lang,
  );
  return (
    <span
      className={`inline-flex max-w-full items-center gap-2 rounded-full border px-3 py-1 text-xs font-semibold leading-5 ${TONE_STYLES[tone]}`}
      title={hint}
      aria-label={`${label}${lang === "zh" ? "。" : ". "}${hint}`}
    >
      <span
        aria-hidden="true"
        className="h-1.5 w-1.5 shrink-0 rounded-full bg-current"
      />
      <span>{label}</span>
    </span>
  );
}
