#!/usr/bin/env python3
"""Find uncovered lines in any module's service.go"""
import subprocess, os, sys

from pathlib import Path
SERVER_DIR = str(Path(__file__).resolve().parents[1] / "server")
module = sys.argv[1] if len(sys.argv) > 1 else "search"
cov_file = os.path.join(SERVER_DIR, f"cov_{module}.out")

subprocess.run(
    ["go", "test", f"./internal/{module}/...", "-count=1", f"-coverprofile={cov_file}"],
    capture_output=True, text=True, cwd=SERVER_DIR, timeout=60
)

total = covered = 0
with open(cov_file, "r", encoding="utf-8", errors="replace") as f:
    for line in f:
        if "/service.go:" in line and "/service_test.go:" not in line:
            parts = line.strip().split()
            if len(parts) >= 3:
                try:
                    s, c = int(parts[1]), int(parts[2])
                    total += s
                    if c > 0: covered += s
                    else: print(f"UNCOV: {parts[0]} ({s} stmts)")
                except: pass

print(f"\n{module} service coverage: {covered}/{total} = {covered*100/total:.1f}%" if total else f"no service.go found")
os.remove(cov_file)
