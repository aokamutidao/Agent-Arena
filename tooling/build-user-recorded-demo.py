#!/usr/bin/env python3
"""Assemble the supplied screen recordings, actual UI captures and narrated cards.

This offline editor never controls a browser or performs a business action. It
requires Python's standard library plus ffmpeg/ffprobe. Generate cue WAVs first
with build-demo-video.py --prepare-audio; this script fits them to the edit.
"""

from __future__ import annotations

import argparse
import concurrent.futures
import hashlib
import json
import math
import re
import shutil
import subprocess
import sys
import textwrap
import unicodedata
import wave
from pathlib import Path

VERSION = 2
ROOT = Path(__file__).resolve().parents[1]
FPS = 30
RATE = 48000
WIDTH, HEIGHT = 1920, 1080
BACKGROUND = "0x101827"
CROP = "3408:1788:8:248"
CARDS = {"title": "01-title", "roles": "02-role-flow", "architecture": "03-architecture",
         "contracts": "04-contracts", "engine": "05-agent-evidence", "security": "06-security", "end": "07-end"}
SOURCES = {1: "1赛事准备.mov", 2: "2 开始对战.mov", 3: "3审计检查.mov", 4: "4 领取奖励.mov"}


class EditError(Exception):
    pass


def run(argv: list[str], cwd: Path | None = None) -> str:
    result = subprocess.run(argv, cwd=cwd, capture_output=True, text=True)
    if result.returncode:
        raise EditError(f"{Path(argv[0]).name} failed ({result.returncode}):\n{result.stderr[-6000:]}")
    return result.stdout


def read_json(path: Path) -> dict:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
    except (OSError, ValueError) as exc:
        raise EditError(f"Cannot read {path}: {exc}") from exc
    if not isinstance(value, dict):
        raise EditError(f"JSON object expected: {path}")
    return value


def write_json(path: Path, value: dict) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    temp = path.with_suffix(path.suffix + ".tmp")
    temp.write_text(json.dumps(value, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    temp.replace(path)


def tool(name: str, override: str | None) -> str:
    found = override or shutil.which(name)
    if not found:
        raise EditError(f"Required tool {name} is unavailable.")
    return str(Path(found).expanduser().resolve())


def stat_key(path: Path) -> dict:
    if not path.is_file():
        raise EditError(f"Missing input: {path}")
    stat = path.stat()
    return {"path": str(path.resolve()), "size": stat.st_size, "mtime_ns": stat.st_mtime_ns}


def fingerprint(value: object) -> str:
    return hashlib.sha256(json.dumps(value, ensure_ascii=False, sort_keys=True).encode()).hexdigest()


def ffconcat_path(path: Path) -> str:
    name = str(path.resolve())
    if "\n" in name or "\r" in name:
        raise EditError("Newline in media path is unsupported.")
    return "'" + name.replace("'", "'\\''") + "'"


def number(value: object, label: str, minimum: float = 0.0) -> float:
    if isinstance(value, bool) or not isinstance(value, (int, float)):
        raise EditError(f"{label} must be numeric.")
    value = float(value)
    if not math.isfinite(value) or value < minimum:
        raise EditError(f"Invalid {label}: {value}")
    return value


def probe(path: Path, args: argparse.Namespace) -> dict:
    return json.loads(run([args.ffprobe, "-v", "error", "-show_streams", "-show_format", "-of", "json", str(path)]))


def duration(path: Path, args: argparse.Namespace) -> float:
    return float(probe(path, args)["format"]["duration"])


def make_plan(args: argparse.Namespace) -> dict:
    narration = read_json(args.narration)
    chapters = narration.get("chapters")
    if not isinstance(chapters, list) or not chapters:
        raise EditError("Narration needs chapters.")
    scenes, seen, start = [], set(), 0.0
    for chapter in chapters:
        scene = dict(chapter)
        scene_id = str(scene.get("id", ""))
        if not re.fullmatch(r"[A-Za-z0-9][A-Za-z0-9_-]*", scene_id) or scene_id in seen:
            raise EditError(f"Invalid/duplicate scene id: {scene_id}")
        seen.add(scene_id)
        if not isinstance(scene.get("title"), str) or not scene.get("cues"):
            raise EditError(f"Scene {scene_id} needs a title and cues.")
        for cue in scene["cues"]:
            if not all(isinstance(cue.get(lang), str) and cue[lang].strip() for lang in ("en", "zh")):
                raise EditError(f"Scene {scene_id}: every cue needs English and Chinese.")
        planned = number(scene.get("duration"), "scene duration", 0.1)
        visual = scene.get("visual")
        content_duration = 0.0
        if visual == "source":
            file_id = scene.get("file")
            if file_id not in SOURCES:
                raise EditError(f"Unknown source id {file_id}")
            scene["source_file"] = str(args.sources / SOURCES[file_id])
            if not scene.get("ranges"):
                raise EditError(f"Scene {scene_id}: source needs ranges.")
            ranges = []
            for span in scene["ranges"]:
                if not isinstance(span, list) or len(span) != 2:
                    raise EditError("Each source range must be [in, out].")
                first, last = number(span[0], "range in"), number(span[1], "range out")
                if last <= first:
                    raise EditError(f"Empty/reversed range: {span}")
                ranges.append([first, last])
                content_duration += last - first
            scene["ranges"] = ranges
        elif visual in CARDS:
            scene["card_file"] = str(args.cards / (CARDS[visual] + ".png"))
        elif visual in ("home", "funds", "guide", "audit", "refund"):
            scene["capture_manifest"] = str(args.intro / visual / "manifest.json")
            capture_path = Path(scene["capture_manifest"])
            if capture_path.is_file():
                content_duration = number(read_json(capture_path).get("duration"), "capture duration", 0.01)
        else:
            raise EditError(f"Unknown visual: {visual}")
        # An edit cannot silently shorten a chosen source interval or capture.
        frame_count = math.ceil(max(planned, content_duration) * FPS - 1e-8)
        scene.update({"duration": frame_count / FPS, "planned_duration": planned,
                      "frames": frame_count, "start": start, "end": start + frame_count / FPS,
                      "content_duration": content_duration,
                      "freeze_seconds": max(0.0, frame_count / FPS - content_duration)})
        start = scene["end"]
        scenes.append(scene)
    return {"version": VERSION, "title": narration.get("title", "Agent Arena"),
            "narration_source": str(args.narration), "narration_fingerprint": fingerprint(narration),
            "source_crop": CROP, "source_playback_speed": 1.0,
            "source_audio": "removed", "output_size": [WIDTH, HEIGHT], "fps": FPS,
            "caption_strip": {"top": 900, "height": 180}, "chapters": scenes,
            "total_duration": start,
            "provenance": "User screen recordings and timestamped real browser captures; architecture cards are explanatory graphics. No business action or UI state is synthesized."}


def wav_frames(path: Path) -> int:
    try:
        with wave.open(str(path), "rb") as source:
            if (source.getnchannels(), source.getsampwidth(), source.getframerate()) != (1, 2, RATE):
                raise EditError(f"Expected mono 16-bit {RATE}Hz cue: {path}")
            return source.getnframes()
    except (OSError, wave.Error) as exc:
        raise EditError(f"Invalid WAV {path}: {exc}") from exc


def silence(target: wave.Wave_write, frames: int) -> None:
    while frames:
        part = min(frames, RATE)
        target.writeframesraw(b"\0\0" * part)
        frames -= part


def append_wav(target: wave.Wave_write, path: Path) -> None:
    with wave.open(str(path), "rb") as source:
        while block := source.readframes(RATE):
            target.writeframesraw(block)


def new_wav(path: Path) -> wave.Wave_write:
    target = wave.open(str(path), "wb")
    target.setnchannels(1)
    target.setsampwidth(2)
    target.setframerate(RATE)
    return target


def audio_blocks(scene: dict) -> list[tuple[int, int, int, int, bool]]:
    """Partition cue indices at explicit picture-sync anchors, in audio samples."""
    boundaries = [(0, 0, False)]
    for index, cue in enumerate(scene["cues"]):
        if "at" not in cue:
            continue
        at = round(number(cue["at"], "cue anchor") * RATE)
        if index == 0:
            boundaries[0] = (0, at, True)
        else:
            boundaries.append((index, at, True))
    boundaries.append((len(scene["cues"]), round(scene["duration"] * RATE), False))
    blocks = []
    for (lo, start, anchored), (hi, end, _) in zip(boundaries, boundaries[1:]):
        if hi <= lo or end <= start:
            raise EditError(f"Scene {scene['id']}: anchors must increase and remain within the scene.")
        blocks.append((lo, hi, start, end, anchored))
    return blocks


def fit_audio(args: argparse.Namespace, plan: dict) -> dict:
    raw = read_json(args.voice / "timed-narration.json")
    chapters = {str(ch["id"]): ch for ch in raw.get("chapters", [])}
    audio_dir = args.work / "audio"
    audio_dir.mkdir(parents=True, exist_ok=True)
    fitted = {"version": VERSION, "voice": raw.get("voice"), "rate": raw.get("rate"),
              "total_duration": plan["total_duration"], "chapters": []}
    chapter_paths = []
    for scene in plan["chapters"]:
        scene_id = scene["id"]
        original = chapters.get(scene_id)
        if not original or len(original.get("cues", [])) != len(scene["cues"]):
            raise EditError(f"Raw narration for scene {scene_id} is missing or stale.")
        paths, counts = [], []
        for expected, cue in zip(scene["cues"], original["cues"]):
            if expected["en"] != cue.get("en") or expected["zh"] != cue.get("zh"):
                raise EditError(f"Raw speech text no longer matches scene {scene_id}; regenerate voice.")
            path = args.voice / cue["audio_file"]
            paths.append(path)
            counts.append(wav_frames(path))
        total = round(scene["duration"] * RATE)
        gaps = max(0, len(paths) - 1)
        min_gap = round(0.2 * RATE)
        blocks = audio_blocks(scene)
        speed = 1.0
        for lo, hi, begin, end, anchored in blocks:
            block_gaps = hi - lo - 1 + int(hi < len(paths))
            budget = end - begin - block_gaps * min_gap
            if budget <= 0:
                raise EditError(f"Scene {scene_id}: anchored block has no speech capacity.")
            if sum(counts[lo:hi]) > budget:
                speed = max(speed, sum(counts[lo:hi]) / (budget - round(0.06 * RATE)))
        if speed > 1.0:
            if speed > 1.12 + 1e-9:
                raise EditError(f"Scene {scene_id}: speech needs {sum(counts)/RATE:.2f}s plus gaps in {scene['duration']:.2f}s; required atempo {speed:.3f} exceeds 1.12. Shorten narration or lengthen scene.")
            adjusted, adjusted_counts = [], []
            for index, path in enumerate(paths):
                key = fingerprint({"version": VERSION, "input": stat_key(path), "speed": speed})
                output = audio_dir / f"{scene_id}-{index+1:03d}-{key[:16]}.wav"
                if not output.is_file():
                    temporary = output.with_suffix(".tmp.wav")
                    run([args.ffmpeg, "-nostdin", "-v", "error", "-y", "-i", str(path),
                         "-af", f"atempo={speed:.9f}", "-ar", str(RATE), "-ac", "1", "-c:a", "pcm_s16le", str(temporary)])
                    temporary.replace(output)
                adjusted.append(output)
                adjusted_counts.append(wav_frames(output))
            paths, counts = adjusted, adjusted_counts
        starts = [0] * len(paths)
        gap_sizes = []
        for lo, hi, begin, end, anchored in blocks:
            block_gaps = hi - lo - 1 + int(hi < len(paths))
            extra = end - begin - sum(counts[lo:hi]) - block_gaps * min_gap
            if extra < 0:
                raise EditError(f"Scene {scene_id}: anchored speech block still overruns by {-extra/RATE:.3f}s after atempo.")
            gap_size = min_gap + (min(round(0.8 * RATE), extra // block_gaps) if block_gaps else 0)
            remaining = end - begin - sum(counts[lo:hi]) - block_gaps * gap_size
            # Fixed picture-sync cues start exactly at their anchor. For an
            # unanchored prefix, spare space precedes speech so its final pause
            # before the next anchored cue is at most one second.
            lead = 0 if anchored else (remaining if hi < len(paths) else remaining // 2)
            cursor = begin + lead
            for index in range(lo, hi):
                starts[index] = cursor
                cursor += counts[index] + gap_size
            gap_sizes.append(gap_size / RATE)
        for index, cue in enumerate(scene["cues"]):
            if "at" in cue and abs(starts[index] / RATE - cue["at"]) > 0.5:
                raise EditError(f"Scene {scene_id} cue {index+1}: picture-sync anchor conflicts by more than 0.5s.")
        cursor, cues = 0, []
        output = audio_dir / f"scene-{scene_id}.wav"
        with new_wav(output) as target:
            for index, (path, count, cue) in enumerate(zip(paths, counts, scene["cues"])):
                if starts[index] < cursor:
                    raise EditError(f"Scene {scene_id}: narration cues overlap.")
                silence(target, starts[index] - cursor)
                cursor = starts[index]
                append_wav(target, path)
                cues.append({"id": f"{scene_id}-{index+1:03d}", **cue,
                             "relative_start": cursor / RATE, "relative_end": (cursor + count) / RATE,
                             "start": scene["start"] + cursor / RATE,
                             "end": scene["start"] + (cursor + count) / RATE,
                             "duration": count / RATE, "audio_file": str(path.resolve())})
                cursor += count
            if cursor > total:
                raise EditError(f"Scene {scene_id}: narration overruns scene duration.")
            silence(target, total - cursor)
        if wav_frames(output) != total:
            raise EditError(f"Scene {scene_id}: audio sample count mismatch.")
        chapter_paths.append(output)
        fitted["chapters"].append({"id": scene_id, "title": scene["title"],
                                   "start": scene["start"], "end": scene["end"], "duration": scene["duration"],
                                   "atempo": speed, "lead_seconds": starts[0] / RATE,
                                   "block_gap_seconds": gap_sizes,
                                   "tail_seconds": (total - cursor) / RATE,
                                   "audio_file": str(output.resolve()), "cues": cues})
        print(f"Audio {scene_id}: {scene['duration']:.2f}s, speech {sum(counts)/RATE:.2f}s, atempo {speed:.3f}", flush=True)
    narration = audio_dir / "narration.wav"
    with new_wav(narration) as target:
        for path in chapter_paths:
            append_wav(target, path)
    run([args.ffmpeg, "-nostdin", "-v", "error", "-y", "-i", str(narration),
         "-af", "loudnorm=I=-16:TP=-1.5:LRA=11", "-ar", str(RATE),
         "-c:a", "aac", "-b:a", "192k", "-movflags", "+faststart", str(args.output / "narration.m4a")])
    write_json(args.output / "timed-narration.json", fitted)
    write_subtitles(args.output, fitted)
    write_chapters(args.output, plan)
    return fitted


def timestamp(value: float, ass: bool = False) -> str:
    scale = 100 if ass else 1000
    ticks = round(value * scale)
    hours, ticks = divmod(ticks, 3600 * scale)
    minutes, ticks = divmod(ticks, 60 * scale)
    secs, fraction = divmod(ticks, scale)
    return f"{hours}:{minutes:02}:{secs:02}.{fraction:02}" if ass else f"{hours:02}:{minutes:02}:{secs:02},{fraction:03}"


def wrap_caption(text: str, language: str) -> list[str]:
    text = " ".join(text.split())
    if language == "en":
        lines = textwrap.wrap(text, width=104, break_long_words=False, break_on_hyphens=False)
    else:
        lines, line, units = [], "", 0.0
        for char in text:
            width = 1.0 if unicodedata.east_asian_width(char) in ("W", "F") else 0.55
            if units + width > 57 and line:
                lines.append(line)
                line, units = "", 0.0
            line += char
            units += width
        if line:
            lines.append(line)
    if len(lines) > 2:
        raise EditError(f"Caption exceeds two {language} lines; shorten/split cue: {text}")
    return lines


def ass_text(text: str) -> str:
    return text.replace("\\", "\\\\").replace("{", "\\{").replace("}", "\\}")


def ass_header() -> str:
    return """[Script Info]
ScriptType: v4.00+
PlayResX: 1920
PlayResY: 1080
WrapStyle: 2
ScaledBorderAndShadow: yes

[V4+ Styles]
Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding
Style: English,Arial,32,&H00FFFFFF,&H00FFFFFF,&H00101827,&H00101827,0,0,0,0,100,100,0,0,1,0,0,8,60,60,0,1
Style: Chinese,PingFang SC,30,&H00E8EBEF,&H00FFFFFF,&H00101827,&H00101827,0,0,0,0,100,100,0,0,1,0,0,8,60,60,0,1
Style: Title,Arial,36,&H00F5F7FA,&H00FFFFFF,&H00101827,&H00101827,-1,0,0,0,100,100,0,0,1,0,0,7,64,64,0,1

[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
"""


def write_subtitles(output: Path, timeline: dict) -> None:
    cues = [cue for chapter in timeline["chapters"] for cue in chapter["cues"]]
    srt = {"en": [], "zh": [], "bilingual": []}
    ass = []
    for index, cue in enumerate(cues, 1):
        wrapped = {lang: wrap_caption(cue[lang], lang) for lang in ("en", "zh")}
        timing = f"{timestamp(cue['start'])} --> {timestamp(cue['end'])}"
        for lang in srt:
            text = "\n".join(wrapped[lang] if lang != "bilingual" else wrapped["en"] + wrapped["zh"])
            srt[lang].append(f"{index}\n{timing}\n{text}\n")
        for lang, style, ypos in (("en", "English", 916), ("zh", "Chinese", 997)):
            text = r"\N".join(ass_text(line) for line in wrapped[lang])
            ass.append(f"Dialogue: 0,{timestamp(cue['start'], True)},{timestamp(cue['end'], True)},{style},,0,0,0,,{{\\pos(960,{ypos})\\q2}}{text}")
    for lang, entries in srt.items():
        (output / f"subtitles-{lang}.srt").write_text("\n".join(entries), encoding="utf-8")
    (output / "subtitles.ass").write_text(ass_header() + "\n".join(ass) + "\n", encoding="utf-8")


def write_chapters(output: Path, plan: dict) -> None:
    lines = ["# Agent Arena — chapter timestamps", "", "English narration with Chinese and English subtitles.", ""]
    for scene in plan["chapters"]:
        seconds = int(scene["start"])
        lines.append(f"- {seconds//60:02}:{seconds%60:02} — {scene['title']}")
    (output / "chaptertimestamps.md").write_text("\n".join(lines) + "\n", encoding="utf-8")


def capture_inputs(scene: dict, work: Path) -> tuple[Path, list[dict]]:
    manifest = Path(scene["capture_manifest"])
    data = read_json(manifest)
    frames = data.get("frames")
    if not frames:
        raise EditError(f"No capture frames: {manifest}")
    length = number(data.get("duration"), "capture duration", 0.01)
    resolved, dependencies = [], [stat_key(manifest)]
    previous = -1.0
    for frame in frames:
        t = number(frame.get("t"), "frame timestamp")
        if t <= previous or t >= length:
            raise EditError(f"Nonmonotonic/out-of-range capture timestamps: {manifest}")
        path = (manifest.parent / frame["file"]).resolve()
        if not path.is_relative_to(manifest.parent.resolve()):
            raise EditError(f"Capture frame escapes capture directory: {path}")
        resolved.append((path, t))
        dependencies.append(stat_key(path))
        previous = t
    if resolved[0][1] > 0.001:
        raise EditError(f"Capture must begin at t=0: {manifest}")
    lines = ["ffconcat version 1.0"]
    for index, (path, t) in enumerate(resolved):
        end = resolved[index + 1][1] if index + 1 < len(resolved) else length
        lines.extend([f"file {ffconcat_path(path)}", f"duration {end - t:.9f}"])
    lines.append(f"file {ffconcat_path(resolved[-1][0])}")
    concat = work / "capture.ffconcat"
    concat.write_text("\n".join(lines) + "\n", encoding="utf-8")
    return concat, dependencies


def validate_video(path: Path, expected: float, args: argparse.Namespace, audio: bool) -> dict:
    data = probe(path, args)
    videos = [s for s in data["streams"] if s["codec_type"] == "video"]
    audios = [s for s in data["streams"] if s["codec_type"] == "audio"]
    if len(videos) != 1:
        raise EditError(f"Expected exactly one video stream: {path}")
    video = videos[0]
    if (video["width"], video["height"], video.get("pix_fmt")) != (WIDTH, HEIGHT, "yuv420p"):
        raise EditError(f"Wrong video geometry/pixel format: {path}")
    if video.get("codec_name") != "h264" or video.get("avg_frame_rate") != "30/1":
        raise EditError(f"Expected H264, constant 30fps: {path}")
    vd = float(video.get("duration", data["format"]["duration"]))
    if abs(vd - expected) > 2 / FPS + 0.002:
        raise EditError(f"Unexpected video duration {vd:.3f}, expected {expected:.3f}: {path}")
    result = {"path": str(path), "width": WIDTH, "height": HEIGHT, "fps": 30, "duration": vd}
    if audio:
        if len(audios) != 1:
            raise EditError(f"Expected exactly one narration stream: {path}")
        ad = float(audios[0].get("duration", data["format"]["duration"]))
        offset = abs(float(video.get("start_time", 0)) - float(audios[0].get("start_time", 0)))
        if abs(ad - vd) > 0.12 or offset > 0.05:
            raise EditError(f"Audio/video sync mismatch: {path}, durations {ad}/{vd}, start offset {offset}")
        result.update({"audio_duration": ad, "av_duration_delta": abs(ad - vd), "av_start_delta": offset})
    elif audios:
        raise EditError(f"Source audio was not removed: {path}")
    return result


def render_scene(scene: dict, args: argparse.Namespace) -> Path:
    work = args.work / "scenes" / scene["id"]
    work.mkdir(parents=True, exist_ok=True)
    output = work / "scene.mp4"
    cache = work / "cache.json"
    dependencies, input_args, filters = [], [], []
    visual = scene["visual"]
    canvas = ("scale=1820:826:force_original_aspect_ratio=decrease:force_divisible_by=2,"
              f"setsar=1,pad=1920:900:(ow-iw)/2:60+(840-ih)/2:color={BACKGROUND},"
              f"pad=1920:1080:0:0:color={BACKGROUND},fps={FPS},format=yuv420p")
    if visual == "source":
        source = Path(scene["source_file"])
        dependencies.append(stat_key(source))
        metadata = probe(source, args)
        video = next(stream for stream in metadata["streams"] if stream["codec_type"] == "video")
        if video["width"] < 3416 or video["height"] < 2036:
            raise EditError(f"Source too small for requested crop: {source}")
        length = float(metadata["format"]["duration"])
        for index, (start, end) in enumerate(scene["ranges"]):
            if end > length + 0.02:
                raise EditError(f"Source range exceeds {source}: {end} > {length}")
            input_args.extend(["-ss", f"{start:.9f}", "-t", f"{end-start:.9f}", "-i", str(source)])
            filters.append(f"[{index}:v]setpts=PTS-STARTPTS,crop={CROP},{canvas}[part{index}]")
        count = len(scene["ranges"])
        if count > 1:
            filters.append("".join(f"[part{i}]" for i in range(count)) + f"concat=n={count}:v=1:a=0[content]")
        else:
            filters.append("[part0]null[content]")
    elif visual in CARDS:
        card = Path(scene["card_file"])
        dependencies.append(stat_key(card))
        input_args.extend(["-loop", "1", "-framerate", str(FPS), "-i", str(card)])
        filters.append(f"[0:v]scale=1920:900,setsar=1,pad=1920:1080:0:0:color={BACKGROUND},fps={FPS},format=yuv420p[content]")
    else:
        concat, dependencies = capture_inputs(scene, work)
        input_args.extend(["-safe", "0", "-f", "concat", "-i", str(concat)])
        capture_canvas = canvas.replace(
            "force_divisible_by=2,", "force_divisible_by=2:in_range=pc:out_range=tv,"
        ).replace("setsar=1,", "setparams=range=limited,setsar=1,")
        filters.append(f"[0:v]setpts=PTS-STARTPTS,{capture_canvas}[content]")
    # Unlimited clone extension; output frame count equals the planned scene.
    # make_plan ensures that frame count can contain every selected source tail.
    final = "[content]tpad=stop_mode=clone:stop=-1"
    if visual not in CARDS:
        title = ass_header() + f"Dialogue: 0,0:00:00.00,{timestamp(scene['duration'], True)},Title,,0,0,0,,{{\\pos(64,13)}}{ass_text(scene['title'])}\n"
        (work / "title.ass").write_text(title, encoding="utf-8")
        final += ",ass=filename=title.ass"
    final += "[out]"
    filters.append(final)
    key = fingerprint({"version": VERSION, "script": stat_key(Path(__file__)), "scene": scene,
                       "dependencies": dependencies, "crop": CROP, "canvas": canvas,
                       "ffmpeg": stat_key(Path(args.ffmpeg)), "encoding": ["libx264", "veryfast", 19]})
    if output.is_file() and cache.is_file() and read_json(cache).get("key") == key:
        validate_video(output, scene["duration"], args, audio=False)
        print(f"Cached scene {scene['id']}: {scene['title']}", flush=True)
        return output
    temporary = work / "scene.tmp.mp4"
    run([args.ffmpeg, "-nostdin", "-v", "error", "-y", "-threads", "2", *input_args,
         "-filter_complex_threads", "1", "-filter_complex", ";".join(filters), "-map", "[out]", "-an",
         "-frames:v", str(scene["frames"]), "-c:v", "libx264", "-preset", "veryfast", "-crf", "19",
         "-pix_fmt", "yuv420p", "-r", str(FPS), "-threads", "2", "-video_track_timescale", "30000",
         *(["-color_range", "tv"] if visual != "source" and visual not in CARDS else []),
         "-movflags", "+faststart", str(temporary)], cwd=work)
    validation = validate_video(temporary, scene["duration"], args, audio=False)
    temporary.replace(output)
    write_json(cache, {"key": key, "dependencies": dependencies, "validation": validation})
    print(f"Rendered scene {scene['id']}: {scene['title']}", flush=True)
    return output


def render(args: argparse.Namespace, plan: dict, timeline: dict) -> None:
    with concurrent.futures.ThreadPoolExecutor(max_workers=args.jobs) as executor:
        futures = [executor.submit(render_scene, scene, args) for scene in plan["chapters"]]
        # Preserve edit order even though workers finish independently.
        paths = [future.result() for future in futures]
    concat = args.work / "scenes.ffconcat"
    concat.write_text("ffconcat version 1.0\n" + "\n".join(f"file {ffconcat_path(p)}" for p in paths) + "\n", encoding="utf-8")
    clean = args.output / "Agent-Arena-Demo-clean.mp4"
    run([args.ffmpeg, "-nostdin", "-v", "error", "-y", "-safe", "0", "-f", "concat", "-i", str(concat),
         "-i", str(args.output / "narration.m4a"), "-map", "0:v:0", "-map", "1:a:0", "-c", "copy",
         "-movflags", "+faststart", str(clean)])
    clean_check = validate_video(clean, plan["total_duration"], args, audio=True)
    final = args.output / "Agent-Arena-Demo-EN-ZH.mp4"
    run([args.ffmpeg, "-nostdin", "-v", "error", "-y", "-threads", "2", "-i", str(clean),
         "-vf", "ass=filename=subtitles.ass", "-filter_threads", "1", "-c:v", "libx264", "-preset", "veryfast",
         "-crf", "20", "-pix_fmt", "yuv420p", "-r", str(FPS), "-threads", "2", "-c:a", "copy",
         "-movflags", "+faststart", str(final)], cwd=args.output)
    final_check = validate_video(final, plan["total_duration"], args, audio=True)
    run([args.ffmpeg, "-nostdin", "-v", "error", "-y", "-ss", "1", "-i", str(final),
         "-frames:v", "1", "-update", "1", str(args.output / "cover.png")])
    manifest = {**plan, "audio_timeline": timeline, "validation": {"clean": clean_check, "captioned": final_check},
                "outputs": ["Agent-Arena-Demo-EN-ZH.mp4", "Agent-Arena-Demo-clean.mp4", "narration.m4a",
                            "subtitles-en.srt", "subtitles-zh.srt", "subtitles-bilingual.srt", "subtitles.ass",
                            "timed-narration.json", "chaptertimestamps.md", "cover.png"],
                "scene_caches": [read_json(args.work / "scenes" / s["id"] / "cache.json") for s in plan["chapters"]]}
    write_json(args.output / "manifest.json", manifest)
    print(f"Finished {plan['total_duration']:.2f}s: {final}", flush=True)


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    mode = parser.add_mutually_exclusive_group(required=True)
    mode.add_argument("--plan", action="store_true", help="Validate narration/edit structure and write edit-plan.json; no media rendering")
    mode.add_argument("--audio", action="store_true", help="Fit pre-generated cue WAVs to scenes and create narration/subtitles")
    mode.add_argument("--render", action="store_true", help="Fit narration, render scenes, assemble both videos and validate")
    mode.add_argument("--all", action="store_true", help="Plan, fit audio, render and validate; speech cue WAVs must already exist")
    parser.add_argument("--narration", type=Path, default=ROOT / "docs/course-project/demo-user-recording/NARRATION.en-zh.json")
    parser.add_argument("--voice", type=Path, default=ROOT / ".local/user-demo-edit/voice")
    parser.add_argument("--sources", type=Path, default=Path.home() / "Desktop")
    parser.add_argument("--intro", type=Path, default=ROOT / ".local/user-demo-edit/intro")
    parser.add_argument("--cards", type=Path, default=ROOT / ".local/user-demo-edit/cards")
    parser.add_argument("--work", type=Path, default=ROOT / ".local/user-demo-edit/render")
    parser.add_argument("--output", type=Path, default=ROOT / "deliverables/demo-video-user")
    parser.add_argument("--jobs", type=int, choices=(1, 2), default=2)
    parser.add_argument("--ffmpeg")
    parser.add_argument("--ffprobe")
    args = parser.parse_args()
    try:
        for key in ("narration", "voice", "sources", "intro", "cards", "work", "output"):
            setattr(args, key, getattr(args, key).expanduser().resolve())
        args.ffmpeg, args.ffprobe = tool("ffmpeg", args.ffmpeg), tool("ffprobe", args.ffprobe)
        args.output.mkdir(parents=True, exist_ok=True)
        args.work.mkdir(parents=True, exist_ok=True)
        plan = make_plan(args)
        write_json(args.output / "edit-plan.json", plan)
        print(f"Edit: {len(plan['chapters'])} scenes, {plan['total_duration']:.2f} seconds", flush=True)
        if args.plan:
            for scene in plan["chapters"]:
                print(f"{scene['id']}: {scene['start']:7.2f}–{scene['end']:7.2f} {scene['visual']:12} {scene['title']}")
            return 0
        timeline = fit_audio(args, plan)
        if args.render or args.all:
            render(args, plan, timeline)
        return 0
    except (EditError, OSError, ValueError, KeyError, StopIteration) as exc:
        print(f"Build error: {exc}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
