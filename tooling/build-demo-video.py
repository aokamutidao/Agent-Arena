#!/usr/bin/env python3
"""Build an offline narrated demo from real, timestamped UI captures.

Only Python's standard library, macOS say, ffmpeg and ffprobe are required.
This program does not operate a browser or record the screen.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import math
import re
import shutil
import subprocess
import sys
import textwrap
import wave
from fractions import Fraction
from pathlib import Path


SAMPLE_RATE = 48_000
BACKGROUND = "0x101827"
ID_PATTERN = re.compile(r"^[A-Za-z0-9][A-Za-z0-9_-]*$")


class DemoError(Exception):
    """An actionable input or media-tool failure."""


def run(command: list[str], cwd: Path | None = None) -> str:
    result = subprocess.run(command, capture_output=True, text=True, cwd=cwd)
    if result.returncode:
        details = (result.stderr or result.stdout).strip()[-5000:]
        raise DemoError(f"{Path(command[0]).name} exited {result.returncode}: {details}")
    return result.stdout


def find_tool(name: str, override: str | None = None) -> str:
    if override:
        candidate = Path(override).expanduser().resolve()
        if not candidate.is_file():
            raise DemoError(f"Tool does not exist: {candidate}")
        return str(candidate)
    located = shutil.which(name)
    if located:
        return located
    for prefix in ("/opt/homebrew/bin", "/usr/local/bin", "/usr/bin"):
        candidate = Path(prefix) / name
        if candidate.is_file():
            return str(candidate)
    raise DemoError(f"Required tool '{name}' is unavailable; install it separately or pass --{name}.")


def read_json(path: Path) -> dict:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
    except (OSError, UnicodeError, json.JSONDecodeError) as exc:
        raise DemoError(f"Cannot read JSON {path}: {exc}") from exc
    if not isinstance(value, dict):
        raise DemoError(f"Expected a JSON object: {path}")
    return value


def write_json(path: Path, value: dict) -> None:
    temporary = path.with_suffix(path.suffix + ".tmp")
    temporary.write_text(json.dumps(value, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    temporary.replace(path)


def seconds(value: object, label: str, minimum: float = 0) -> float:
    if isinstance(value, bool) or not isinstance(value, (int, float)):
        raise DemoError(f"{label} must be a number of seconds.")
    number = float(value)
    if not math.isfinite(number) or number < minimum:
        raise DemoError(f"{label} must be finite and >= {minimum}.")
    return number


def validate_narration(path: Path) -> dict:
    value = read_json(path)
    chapters = value.get("chapters")
    if not isinstance(chapters, list) or not chapters:
        raise DemoError("Narration requires a nonempty chapters array.")
    seen: set[str] = set()
    for chapter in chapters:
        if not isinstance(chapter, dict):
            raise DemoError("Each chapter must be an object.")
        chapter_id = chapter.get("id")
        if not isinstance(chapter_id, str) or not ID_PATTERN.fullmatch(chapter_id):
            raise DemoError("Chapter id must contain only ASCII letters, numbers, '_' and '-'.")
        if chapter_id in seen:
            raise DemoError(f"Duplicate chapter id: {chapter_id}")
        seen.add(chapter_id)
        if not isinstance(chapter.get("title"), str):
            raise DemoError(f"Chapter {chapter_id}: title must be a string.")
        cues = chapter.get("cues")
        if not isinstance(cues, list) or not cues:
            raise DemoError(f"Chapter {chapter_id}: cues must be nonempty.")
        for index, cue in enumerate(cues, 1):
            if not isinstance(cue, dict):
                raise DemoError(f"Chapter {chapter_id}, cue {index}: expected an object.")
            for language in ("en", "zh"):
                if not isinstance(cue.get(language), str) or not cue[language].strip():
                    raise DemoError(f"Chapter {chapter_id}, cue {index}: {language} must be nonempty text.")
    return value


def probe_duration(path: Path, ffprobe: str) -> float:
    output = run([ffprobe, "-v", "error", "-show_entries", "format=duration", "-of", "json", str(path)])
    try:
        duration = float(json.loads(output)["format"]["duration"])
    except (ValueError, KeyError, TypeError, json.JSONDecodeError) as exc:
        raise DemoError(f"No readable media duration in {path}. macOS say may require execution outside the sandbox.") from exc
    if not math.isfinite(duration) or duration <= 0.1:
        raise DemoError(f"Invalid/empty audio {path} ({duration}s). macOS say may require execution outside the sandbox.")
    return duration


def wave_frames(path: Path) -> int:
    try:
        with wave.open(str(path), "rb") as source:
            if (source.getnchannels(), source.getsampwidth(), source.getframerate()) != (1, 2, SAMPLE_RATE):
                raise DemoError(f"Expected mono 16-bit {SAMPLE_RATE}Hz WAV: {path}")
            frames = source.getnframes()
    except (OSError, wave.Error) as exc:
        raise DemoError(f"Invalid WAV {path}: {exc}") from exc
    if frames <= SAMPLE_RATE // 10:
        raise DemoError(f"Audio is too short: {path}")
    return frames


def concat_wave(destination: Path, sources: list[Path], silence_frames: int = 0) -> None:
    temporary = destination.with_suffix(".tmp.wav")
    with wave.open(str(temporary), "wb") as target:
        target.setnchannels(1)
        target.setsampwidth(2)
        target.setframerate(SAMPLE_RATE)
        for path in sources:
            wave_frames(path)
            with wave.open(str(path), "rb") as source:
                while chunk := source.readframes(SAMPLE_RATE):
                    target.writeframesraw(chunk)
            if silence_frames:
                target.writeframesraw(b"\0\0" * silence_frames)
    temporary.replace(destination)


def prepare_audio(args: argparse.Namespace, narration: dict, output: Path) -> dict:
    say = find_tool("say", args.say)
    ffmpeg = find_tool("ffmpeg", args.ffmpeg)
    ffprobe = find_tool("ffprobe", args.ffprobe)
    cue_directory = output / "audio" / "cues"
    chapter_directory = output / "audio" / "chapters"
    cue_directory.mkdir(parents=True, exist_ok=True)
    chapter_directory.mkdir(parents=True, exist_ok=True)
    silence_frames = round(args.silence * SAMPLE_RATE)
    timeline: dict = {
        "version": 1,
        "voice": args.voice,
        "rate": args.rate,
        "silence_seconds": silence_frames / SAMPLE_RATE,
        "chapters": [],
    }
    absolute_frames = 0
    chapter_audio: list[Path] = []
    for chapter in narration["chapters"]:
        chapter_id = chapter["id"]
        print(f"Preparing {chapter_id}: {chapter['title']}", flush=True)
        relative_frames = 0
        timed_cues: list[dict] = []
        cue_audio: list[Path] = []
        for index, cue in enumerate(chapter["cues"], 1):
            fingerprint = hashlib.sha256(json.dumps({
                "version": 1, "voice": args.voice, "rate": args.rate, "en": cue["en"],
            }, ensure_ascii=False, sort_keys=True).encode("utf-8")).hexdigest()
            stem = f"{chapter_id}-{index:03d}-{fingerprint[:20]}"
            text_file = cue_directory / f"{stem}.txt"
            aiff = cue_directory / f"{stem}.aiff"
            wav = cue_directory / f"{stem}.wav"
            valid_aiff = False
            if aiff.is_file():
                try:
                    probe_duration(aiff, ffprobe)
                    valid_aiff = True
                except DemoError:
                    pass
            if not valid_aiff:
                text_file.write_text(cue["en"], encoding="utf-8")
                temporary = aiff.with_suffix(".tmp.aiff")
                run([say, "-v", args.voice, "-r", str(args.rate), "-f", str(text_file), "-o", str(temporary)])
                probe_duration(temporary, ffprobe)
                temporary.replace(aiff)
            valid_wav = False
            if wav.is_file() and valid_aiff:
                try:
                    wave_frames(wav)
                    valid_wav = True
                except DemoError:
                    pass
            if not valid_wav:
                temporary = wav.with_suffix(".tmp.wav")
                run([ffmpeg, "-nostdin", "-hide_banner", "-loglevel", "error", "-y", "-i", str(aiff),
                     "-ar", str(SAMPLE_RATE), "-ac", "1", "-c:a", "pcm_s16le", str(temporary)])
                wave_frames(temporary)
                temporary.replace(wav)
            count = wave_frames(wav)
            cue_id = f"{chapter_id}-{index:03d}"
            timed_cues.append({
                "id": cue_id, "en": cue["en"], "zh": cue["zh"],
                "start": round((absolute_frames + relative_frames) / SAMPLE_RATE, 6),
                "end": round((absolute_frames + relative_frames + count) / SAMPLE_RATE, 6),
                "relative_start": round(relative_frames / SAMPLE_RATE, 6),
                "relative_end": round((relative_frames + count) / SAMPLE_RATE, 6),
                "duration": round(count / SAMPLE_RATE, 6),
                "audio_file": str(wav.relative_to(output)),
            })
            relative_frames += count + silence_frames
            cue_audio.append(wav)
        chapter_file = chapter_directory / f"{chapter_id}.wav"
        concat_wave(chapter_file, cue_audio, silence_frames)
        chapter_audio.append(chapter_file)
        timeline["chapters"].append({
            "id": chapter_id, "title": chapter["title"],
            "start": round(absolute_frames / SAMPLE_RATE, 6),
            "end": round((absolute_frames + relative_frames) / SAMPLE_RATE, 6),
            "duration": round(relative_frames / SAMPLE_RATE, 6),
            "audio_file": str(chapter_file.relative_to(output)), "cues": timed_cues,
        })
        absolute_frames += relative_frames
        print(f"  {len(timed_cues)} cues, {relative_frames / SAMPLE_RATE:.2f}s", flush=True)
    total_audio = output / "audio" / "narration.wav"
    concat_wave(total_audio, chapter_audio)
    run([ffmpeg, "-nostdin", "-hide_banner", "-loglevel", "error", "-y", "-i", str(total_audio),
         "-c:a", "aac", "-b:a", "160k", "-movflags", "+faststart", str(output / "audio.m4a")])
    probe_duration(output / "audio.m4a", ffprobe)
    timeline["total_duration"] = round(absolute_frames / SAMPLE_RATE, 6)
    timeline["audio_file"] = "audio.m4a"
    timeline["audio_wav_file"] = str(total_audio.relative_to(output))
    write_json(output / "timed-narration.json", timeline)
    print(f"Audio ready: {timeline['total_duration']:.2f}s; {output / 'timed-narration.json'}", flush=True)
    return timeline


def validate_timeline(timeline: dict, output: Path) -> None:
    if timeline.get("version") != 1:
        raise DemoError("Expected timed-narration.json version 1; run --prepare-audio first.")
    chapters = timeline.get("chapters")
    if not isinstance(chapters, list) or not chapters:
        raise DemoError("The timing JSON has no chapters.")
    previous_end = 0.0
    seen: set[str] = set()
    for chapter in chapters:
        if not isinstance(chapter, dict):
            raise DemoError("Invalid chapter in timing JSON.")
        chapter_id = chapter.get("id")
        if not isinstance(chapter_id, str) or not ID_PATTERN.fullmatch(chapter_id) or chapter_id in seen:
            raise DemoError("Invalid or duplicate chapter id in timing JSON.")
        seen.add(chapter_id)
        start = seconds(chapter.get("start"), f"{chapter_id}.start")
        end = seconds(chapter.get("end"), f"{chapter_id}.end")
        duration = seconds(chapter.get("duration"), f"{chapter_id}.duration", 0.1)
        if abs(start - previous_end) > 0.0001 or abs(end - start - duration) > 0.0001:
            raise DemoError(f"Chapter {chapter_id}: noncontiguous or inconsistent timing.")
        previous_end = end
        cues = chapter.get("cues")
        if not isinstance(cues, list) or not cues:
            raise DemoError(f"Chapter {chapter_id}: missing timed cues.")
        cue_end = start
        for cue in cues:
            if not isinstance(cue, dict):
                raise DemoError(f"Chapter {chapter_id}: invalid timed cue.")
            cue_start = seconds(cue.get("start"), "cue.start")
            next_end = seconds(cue.get("end"), "cue.end")
            if cue_start < cue_end - 0.0001 or next_end <= cue_start or next_end > end + 0.0001:
                raise DemoError(f"Chapter {chapter_id}: overlapping/out-of-bounds subtitle cue.")
            if any(not isinstance(cue.get(language), str) for language in ("en", "zh")):
                raise DemoError(f"Chapter {chapter_id}: missing en/zh subtitle text.")
            cue_end = next_end
    total = seconds(timeline.get("total_duration"), "total_duration", 0.1)
    if abs(total - previous_end) > 0.0001:
        raise DemoError("total_duration disagrees with chapter timing.")
    audio_file = timeline.get("audio_file")
    if not isinstance(audio_file, str) or not resolve_inside(output, audio_file).is_file():
        raise DemoError("The timed narration audio is missing; run --prepare-audio first.")


def resolve_inside(directory: Path, filename: str) -> Path:
    if not filename or "\n" in filename or "\r" in filename:
        raise DemoError("Media filename must be nonempty and contain no newlines.")
    path = (directory / filename).resolve()
    if not path.is_relative_to(directory.resolve()):
        raise DemoError(f"Media filename escapes its directory: {filename}")
    return path


def natural_key(path: Path) -> list:
    return [int(part) if part.isdigit() else part.lower() for part in re.split(r"(\d+)", path.name)]


def capture_frames(captures: Path, chapter_id: str, capture_fps: float) -> tuple[list[tuple[Path, float]], float]:
    directory = captures / chapter_id
    manifest = directory / "manifest.json"
    image_suffixes = {".png", ".jpg", ".jpeg", ".webp", ".bmp"}
    frames: list[tuple[Path, float]] = []
    if manifest.is_file():
        value = read_json(manifest)
        items = value.get("frames")
        if not isinstance(items, list) or not items:
            raise DemoError(f"Capture manifest has no frames: {manifest}")
        previous = -1.0
        for item in items:
            if not isinstance(item, dict) or not isinstance(item.get("file"), str):
                raise DemoError(f"Invalid frame in {manifest}")
            timestamp = seconds(item.get("t"), f"{manifest}: frame.t")
            if timestamp <= previous:
                raise DemoError(f"Frame timestamps must be strictly increasing: {manifest}")
            image = resolve_inside(directory, item["file"])
            if not image.is_file() or image.suffix.lower() not in image_suffixes:
                raise DemoError(f"Capture image missing or unsupported: {image}")
            frames.append((image, timestamp))
            previous = timestamp
        duration = seconds(value.get("duration", frames[-1][1] + 1 / capture_fps), f"{manifest}: duration")
        if duration < frames[-1][1]:
            raise DemoError(f"Capture duration precedes the final frame: {manifest}")
        if frames[0][1] > 0.001:
            print(f"Warning: {chapter_id} starts at {frames[0][1]:.3f}s; its first image will also cover the initial gap.", file=sys.stderr)
        return frames, max(duration, frames[-1][1] + 1 / 30_000)
    if directory.is_dir():
        images = sorted((path for path in directory.iterdir() if path.is_file() and path.suffix.lower() in image_suffixes), key=natural_key)
        if images:
            if len(images) > 1:
                print(f"Warning: {chapter_id} has no manifest; inferring timestamps at {capture_fps:g}fps.", file=sys.stderr)
            return [(path.resolve(), index / capture_fps) for index, path in enumerate(images)], len(images) / capture_fps
    for suffix in sorted(image_suffixes):
        image = captures / f"{chapter_id}{suffix}"
        if image.is_file():
            return [(image.resolve(), 0.0)], 1 / capture_fps
    raise DemoError(f"No manifest or still image for chapter {chapter_id} in {captures}")


def concat_filename(path: Path) -> str:
    value = str(path)
    if "\n" in value or "\r" in value:
        raise DemoError("ffconcat paths cannot contain newlines.")
    return "'" + value.replace("'", "'\\''") + "'"


def capture_concat_text(frames: list[tuple[Path, float]], duration: float) -> str:
    lines = ["ffconcat version 1.0"]
    for index, (image, timestamp) in enumerate(frames):
        # A late first capture is held from chapter start; later times remain intact.
        start = 0.0 if index == 0 else timestamp
        end = frames[index + 1][1] if index + 1 < len(frames) else duration
        lines.extend([f"file {concat_filename(image)}", "option framerate 1000", f"duration {max(end - start, 0.000001):.6f}"])
    # The concat demuxer needs a repeated last file to honour its hold duration.
    lines.extend([f"file {concat_filename(frames[-1][0])}", "option framerate 1000"])
    return "\n".join(lines) + "\n"


def write_capture_concat(path: Path, frames: list[tuple[Path, float]], duration: float) -> None:
    path.write_text(capture_concat_text(frames, duration), encoding="utf-8")


def reusable_chapter(video: Path, concat: Path, expected_concat: str,
                     frames: list[tuple[Path, float]], fps: int, target_duration: float,
                     ffprobe: str) -> bool:
    """Reuse a final chapter only when both its inputs and video format match."""
    if not video.is_file() or not concat.is_file():
        return False
    try:
        if concat.read_text(encoding="utf-8") != expected_concat:
            return False
        video_mtime = video.stat().st_mtime_ns
        if any(image.stat().st_mtime_ns > video_mtime for image, _ in frames):
            return False
        metadata = json.loads(run([
            ffprobe, "-v", "error", "-select_streams", "v:0", "-show_entries",
            "stream=codec_name,pix_fmt,width,height,avg_frame_rate,duration:format=duration",
            "-of", "json", str(video),
        ]))
        stream = metadata["streams"][0]
        duration = float(stream.get("duration", metadata.get("format", {}).get("duration", "nan")))
        actual_fps = float(Fraction(stream["avg_frame_rate"]))
        return (
            stream.get("codec_name") == "h264" and stream.get("pix_fmt") in ("yuv420p", "yuvj420p")
            and stream.get("width") == 1920 and stream.get("height") == 1080
            and abs(actual_fps - fps) < 0.000001
            and math.isfinite(duration) and abs(duration - target_duration) < 0.002
        )
    except (DemoError, OSError, UnicodeError, ValueError, KeyError, TypeError, IndexError, ZeroDivisionError):
        return False


def subtitle_time(value: float, ass: bool = False) -> str:
    units = 100 if ass else 1000
    total = max(0, round(value * units))
    hours, total = divmod(total, 3600 * units)
    minutes, total = divmod(total, 60 * units)
    secs, remainder = divmod(total, units)
    if ass:
        return f"{hours}:{minutes:02d}:{secs:02d}.{remainder:02d}"
    return f"{hours:02d}:{minutes:02d}:{secs:02d},{remainder:03d}"


def ass_text(text: str) -> str:
    # Literal braces/backslashes cannot introduce ASS override commands.
    return text.replace("\\", "＼").replace("{", "｛").replace("}", "｝").replace("\r", "").replace("\n", " ")


def chinese_wrap(text: str, width: int = 48) -> list[str]:
    """Keep closing punctuation with its phrase and ASCII words intact."""
    lines: list[str] = []
    current = ""
    used = 0.0
    closing = set("，。！？；：、）】》」』〕〉］｝,.!?;:)]}…’”%")
    tokens = re.findall(r"[A-Za-z0-9]+(?:[.'’:/_-][A-Za-z0-9]+)*|\s+|.", " ".join(text.split()))
    for token in tokens:
        weight = sum(0.5 if ord(char) < 128 else 1.0 for char in token)
        if token in closing:
            # A punctuation mark may exceed the target by one character; it must
            # not become an isolated line or lead the next phrase.
            if current:
                current += token
                used += weight
            elif lines:
                lines[-1] += token
            else:
                current, used = token, weight
            continue
        if used + weight > width and current:
            lines.append(current.strip())
            current, used = "", 0.0
        if token.isspace() and not current:
            continue
        current += token
        used += weight
    if current.strip():
        lines.append(current.strip())
    return lines or [""]


def write_subtitles(output: Path, timeline: dict) -> None:
    cues = [cue for chapter in timeline["chapters"] for cue in chapter["cues"]]
    for language in ("bilingual", "en", "zh"):
        entries: list[str] = []
        for index, cue in enumerate(cues, 1):
            content = cue["en"] + "\n" + cue["zh"] if language == "bilingual" else cue[language]
            entries.append(f"{index}\n{subtitle_time(cue['start'])} --> {subtitle_time(cue['end'])}\n{content.strip()}\n")
        (output / f"subtitles-{language}.srt").write_text("\n".join(entries), encoding="utf-8")
    header = """[Script Info]
ScriptType: v4.00+
PlayResX: 1920
PlayResY: 1080
WrapStyle: 2
ScaledBorderAndShadow: yes

[V4+ Styles]
Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding
Style: English,Arial,32,&H00FFFFFF,&H00FFFFFF,&H00101827,&H00101827,0,0,0,0,100,100,0,0,1,1,0,8,160,160,0,1
Style: Chinese,PingFang SC,30,&H00D8E9FF,&H00D8E9FF,&H00101827,&H00101827,0,0,0,0,100,100,0,0,1,1,0,8,160,160,0,1

[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
"""
    dialogues: list[str] = []
    for cue in cues:
        en_lines = textwrap.wrap(" ".join(cue["en"].split()), width=86, break_long_words=True, break_on_hyphens=False) or [""]
        zh_lines = chinese_wrap(cue["zh"])
        en_size, zh_size = 32, 30
        # Long cues remain readable without being clipped into the page above.
        if len(en_lines) > 2 or len(zh_lines) > 2 or len(en_lines) + len(zh_lines) > 4:
            print(f"Warning: {cue.get('id', 'cue')} is long; reducing subtitle size. Shorter narration cues are preferable.", file=sys.stderr)
            en_lines = textwrap.wrap(" ".join(cue["en"].split()), width=104, break_long_words=True, break_on_hyphens=False) or [""]
            zh_lines = chinese_wrap(cue["zh"], width=60)
            en_size, zh_size = 27, 25
        height = len(en_lines) * (en_size + 8) + len(zh_lines) * (zh_size + 8) + 6
        if height > 164:
            raise DemoError(f"Cue {cue.get('id')} is too long for the subtitle strip; split it into shorter cues.")
        start, end = subtitle_time(cue["start"], True), subtitle_time(cue["end"], True)
        y = 916
        for line in en_lines:
            dialogues.append(f"Dialogue: 0,{start},{end},English,,0,0,0,,{{\\an8\\pos(960,{y})\\fs{en_size}}}{ass_text(line)}")
            y += en_size + 8
        y += 6
        for line in zh_lines:
            dialogues.append(f"Dialogue: 0,{start},{end},Chinese,,0,0,0,,{{\\an8\\pos(960,{y})\\fs{zh_size}}}{ass_text(line)}")
            y += zh_size + 8
    (output / "subtitles.ass").write_text(header + "\n".join(dialogues) + "\n", encoding="utf-8")


def render_video(args: argparse.Namespace, output: Path, timeline: dict) -> None:
    validate_timeline(timeline, output)
    ffmpeg = find_tool("ffmpeg", args.ffmpeg)
    ffprobe = find_tool("ffprobe", args.ffprobe)
    captures = args.captures.expanduser().resolve()
    if not captures.is_dir():
        raise DemoError(f"Capture directory does not exist: {captures}")
    audio = resolve_inside(output, timeline["audio_file"])
    actual_audio = probe_duration(audio, ffprobe)
    if abs(actual_audio - timeline["total_duration"]) > 0.1:
        raise DemoError("Narration audio duration disagrees with its timing JSON; run --prepare-audio again.")
    write_subtitles(output, timeline)
    video_directory = output / "video" / "chapters"
    video_directory.mkdir(parents=True, exist_ok=True)
    chapter_files: list[Path] = []
    for chapter in timeline["chapters"]:
        chapter_id = chapter["id"]
        frames, captured_duration = capture_frames(captures, chapter_id, args.capture_fps)
        frame_count = round(chapter["end"] * args.fps) - round(chapter["start"] * args.fps)
        if frame_count < 1:
            raise DemoError(f"Chapter {chapter_id} is shorter than one output video frame.")
        target_duration = frame_count / args.fps
        if captured_duration > target_duration + 1 / args.fps:
            print(f"Warning: {chapter_id} capture is {captured_duration:.2f}s but narration is {target_duration:.2f}s; its tail will be trimmed.", file=sys.stderr)
        concat_file = video_directory / f"{chapter_id}.ffconcat"
        chapter_file = video_directory / f"{chapter_id}.mp4"
        expected_concat = capture_concat_text(frames, captured_duration)
        if args.reuse_chapters and reusable_chapter(
            chapter_file, concat_file, expected_concat, frames, args.fps, target_duration, ffprobe,
        ):
            print(f"Reusing verified {chapter_id}: {target_duration:.2f}s (unchanged captures)", flush=True)
            chapter_files.append(chapter_file)
            continue
        print(f"Rendering {chapter_id}: {len(frames)} captures, {target_duration:.2f}s", flush=True)
        concat_file.write_text(expected_concat, encoding="utf-8")
        filters = (
            "scale=1600:900:force_original_aspect_ratio=decrease:force_divisible_by=2,"
            f"pad=1600:900:(ow-iw)/2:(oh-ih)/2:color={BACKGROUND},setsar=1,"
            f"pad=1920:1080:160:0:color={BACKGROUND},fps={args.fps},"
            f"tpad=stop_mode=clone:stop_duration={target_duration:.6f},"
            f"trim=duration={target_duration:.6f},setpts=PTS-STARTPTS"
        )
        temporary = chapter_file.with_suffix(".tmp.mp4")
        run([ffmpeg, "-nostdin", "-hide_banner", "-loglevel", "error", "-y", "-f", "concat", "-safe", "0",
             "-i", str(concat_file), "-an", "-vf", filters, "-frames:v", str(frame_count),
             "-c:v", "libx264", "-preset", "fast", "-crf", "20", "-pix_fmt", "yuv420p", "-movflags", "+faststart", str(temporary)])
        probe_duration(temporary, ffprobe)
        temporary.replace(chapter_file)
        chapter_files.append(chapter_file)
    video_concat = output / "video" / "chapters.ffconcat"
    video_concat.write_text("ffconcat version 1.0\n" + "\n".join(f"file {concat_filename(path)}" for path in chapter_files) + "\n", encoding="utf-8")
    clean = output / "demo-clean.mp4"
    clean_temporary = output / "demo-clean.tmp.mp4"
    total = timeline["total_duration"]
    run([ffmpeg, "-nostdin", "-hide_banner", "-loglevel", "error", "-y", "-f", "concat", "-safe", "0", "-i", str(video_concat),
         "-i", str(audio), "-map", "0:v:0", "-map", "1:a:0", "-c", "copy", "-t", f"{total:.6f}", "-movflags", "+faststart", str(clean_temporary)])
    probe_duration(clean_temporary, ffprobe)
    clean_temporary.replace(clean)
    captioned = output / "demo-captioned.mp4"
    captioned_temporary = output / "demo-captioned.tmp.mp4"
    # A fixed basename in cwd avoids the filter parser's two levels of path escaping.
    run([ffmpeg, "-nostdin", "-hide_banner", "-loglevel", "error", "-y", "-i", str(clean),
         "-vf", "subtitles=filename=subtitles.ass", "-c:v", "libx264", "-preset", "fast", "-crf", "20", "-pix_fmt", "yuv420p",
         "-c:a", "copy", "-movflags", "+faststart", str(captioned_temporary)], cwd=output)
    final_duration = probe_duration(captioned_temporary, ffprobe)
    captioned_temporary.replace(captioned)
    run([ffmpeg, "-nostdin", "-hide_banner", "-loglevel", "error", "-y", "-ss", f"{min(1.0, total / 2):.3f}",
         "-i", str(captioned), "-frames:v", "1", "-update", "1", str(output / "cover.png")])
    write_json(output / "render-summary.json", {
        "version": 1, "width": 1920, "height": 1080, "fps": args.fps,
        "narration_duration": total, "video_duration": final_duration,
        "clean_video": clean.name, "captioned_video": captioned.name,
        "audio": audio.name, "cover": "cover.png",
        "subtitles": ["subtitles.ass", "subtitles-bilingual.srt", "subtitles-en.srt", "subtitles-zh.srt"],
    })
    print(f"Video ready: {captioned} ({final_duration:.2f}s); clean version: {clean}", flush=True)


def parser() -> argparse.ArgumentParser:
    result = argparse.ArgumentParser(
        description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""Narration JSON:
  {"chapters":[{"id":"01-intro","title":"Overview",
    "cues":[{"en":"Welcome to Agent Arena.","zh":"欢迎使用 Agent Arena。"}]}]}

Capture layout (no browser automation is performed):
  captures/01-intro/manifest.json
  {"duration":12.8,"frames":[{"file":"frames001.png","t":0},
    {"file":"frames002.png","t":0.2}]}
  Frame times are real elapsed seconds. Paths are relative to the chapter.
  A single frame in a manifest, or captures/01-intro.png, is also supported.
  A chapter directory containing only numbered images uses --capture-fps.

Examples (run from the repository root):
  python3 tooling/build-demo-video.py --prepare-audio --narration narration.json
  python3 tooling/build-demo-video.py --render --captures captures
  python3 tooling/build-demo-video.py --prepare-audio --render \\
    --narration narration.json --captures captures --rate 135

timed-narration.json v1 uses seconds: chapter and cue start/end are absolute;
cue relative_start/relative_end are relative to the chapter. Cue end excludes
the following silence; chapter duration includes it. Audio paths are relative
to --output. Narration is English; subtitles retain both original languages.
If say returns an empty AIFF in a sandbox, run the same command with approved
local audio-file access. No sound is played. Media files should remain outside Git.
""",
    )
    result.add_argument("--prepare-audio", action="store_true", help="Generate/cache cue speech and the complete timing JSON.")
    result.add_argument("--render", action="store_true", help="Create clean/captioned H.264 MP4s, bilingual/en/zh SRT, ASS and a cover.")
    result.add_argument("--reuse-chapters", action="store_true", help="Reuse verified final chapter MP4s with unchanged captures and matching dimensions/fps/duration; default re-encodes every chapter.")
    result.add_argument("--narration", type=Path, help="UTF-8 narration JSON (required by --prepare-audio).")
    result.add_argument("--captures", type=Path, help="Directory with chapter captures (required by --render).")
    result.add_argument("--output", type=Path, default=Path("deliverables/demo-video"), help="Output/cache directory (default: %(default)s).")
    result.add_argument("--voice", default="Samantha", help="macOS English voice (default: %(default)s).")
    result.add_argument("--rate", type=int, default=135, help="Speech words per minute (default: %(default)s).")
    result.add_argument("--silence", type=float, default=0.25, help="Seconds of silence after every cue (default: %(default)s).")
    result.add_argument("--fps", type=int, default=30, help="Final video frames per second (default: %(default)s).")
    result.add_argument("--capture-fps", type=float, default=5, help="Fallback rate for numbered images without a manifest (default: %(default)s).")
    result.add_argument("--ffmpeg", help="Explicit ffmpeg executable path.")
    result.add_argument("--ffprobe", help="Explicit ffprobe executable path.")
    result.add_argument("--say", help="Explicit macOS say executable path.")
    return result


def main() -> int:
    arguments = parser()
    args = arguments.parse_args()
    if not args.prepare_audio and not args.render:
        arguments.error("Choose --prepare-audio, --render, or both.")
    if args.prepare_audio and not args.narration:
        arguments.error("--prepare-audio requires --narration.")
    if args.render and not args.captures:
        arguments.error("--render requires --captures.")
    if not 1 <= args.rate <= 500 or not 1 <= args.fps <= 60:
        arguments.error("--rate must be 1..500; --fps must be 1..60.")
    try:
        seconds(args.silence, "--silence")
        seconds(args.capture_fps, "--capture-fps", 0.01)
        output = args.output.expanduser().resolve()
        output.mkdir(parents=True, exist_ok=True)
        timeline = prepare_audio(args, validate_narration(args.narration), output) if args.prepare_audio else None
        if args.render:
            timeline = timeline or read_json(output / "timed-narration.json")
            render_video(args, output, timeline)
        return 0
    except (DemoError, OSError) as exc:
        print(f"Error: {exc}", file=sys.stderr)
        return 1
    except KeyboardInterrupt:
        print("Interrupted. Existing valid cue audio remains cached.", file=sys.stderr)
        return 130


if __name__ == "__main__":
    raise SystemExit(main())
