#!/usr/bin/env python3
"""Analyze search service.go coverage by mapping to source lines."""
import subprocess, os

from pathlib import Path
SERVER_DIR = str(Path(__file__).resolve().parents[1] / "server")
cov_file = os.path.join(SERVER_DIR, "search_cov_debug.out")
src_file = os.path.join(SERVER_DIR, "internal", "search", "service.go")

subprocess.run(
    ["go", "test", "./internal/search/...", "-count=1", f"-coverprofile={cov_file}"],
    capture_output=True, text=True, cwd=SERVER_DIR, timeout=30
)

# Read source to get line context
src_lines = {}
with open(src_file, "r", encoding="utf-8", errors="replace") as f:
    for i, line in enumerate(f, 1):
        src_lines[i] = line.rstrip()

# Parse coverage profile
covered_count = 0
uncovered_count = 0

print("=" * 80)
print("UNCOVERED lines in search/service.go:")
print("=" * 80)
with open(cov_file, "r", encoding="utf-8", errors="replace") as f:
    for line in f:
        if "/service.go:" not in line or "/service_test.go:" in line:
            continue
        parts = line.strip().split()
        if len(parts) >= 3:
            try:
                s, c = int(parts[1]), int(parts[2])
                if c == 0:
                    uncovered_count += s
                    # Extract line number from service.go:44.36,51.2
                    line_range = parts[0].split("service.go:")[1]
                    start_line = int(line_range.split(".")[0])
                    src_text = src_lines.get(start_line, "?")
                    print(f"  L{start_line:3d}: {src_text[:80]}")
                else:
                    covered_count += s
            except:
                pass

total = covered_count + uncovered_count
print(f"\nTotal: {total} stmts, Covered: {covered_count}, Uncovered: {uncovered_count}")
print(f"Coverage: {covered_count/total*100:.1f}%")
print(f"Need to cover {int(total*0.8) - covered_count} more stmts to reach 80%")

os.remove(cov_file)
