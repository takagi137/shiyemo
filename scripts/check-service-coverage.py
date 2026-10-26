#!/usr/bin/env python3
"""Calculate service-level test coverage for shiyemo modules."""
import subprocess
import os
import glob

from pathlib import Path
SERVER_DIR = str(Path(__file__).resolve().parents[1] / "server")
MODULES = ["post", "search", "backup"]

# Clean up old coverage files
for f in glob.glob(os.path.join(SERVER_DIR, "cov_*.out")):
    try: os.remove(f)
    except: pass

results = []

for mod in MODULES:
    cov_file = os.path.join(SERVER_DIR, f"cov_{mod}.out")
    try:
        proc = subprocess.run(
            ["go", "test", f"./internal/{mod}/...", "-count=1", f"-coverprofile={cov_file}"],
            capture_output=True, text=True, timeout=60,
            cwd=SERVER_DIR
        )
    except subprocess.TimeoutExpired:
        results.append((mod, 0, 0, "TIMEOUT"))
        continue
    except Exception as e:
        results.append((mod, 0, 0, f"ERR"))
        continue

    if not os.path.exists(cov_file):
        # Try checking stderr for clues
        results.append((mod, 0, 0, "NO FILE"))
        continue

    total_stmts = 0
    covered_stmts = 0
    with open(cov_file, "r", encoding="utf-8", errors="replace") as f:
        for line in f:
            # coverage profile format:
            # pkg/file.go:startCol,endCol numStmt count
            if ("/service.go:" in line and 
                "/service_test.go:" not in line and
                "/service_repo.go:" not in line):
                parts = line.strip().split()
                if len(parts) >= 3:
                    try:
                        stmts = int(parts[1])
                        count = int(parts[2])
                        total_stmts += stmts
                        if count > 0:
                            covered_stmts += stmts
                    except ValueError:
                        pass

    if total_stmts > 0:
        pct = covered_stmts * 100 / total_stmts
        results.append((mod, covered_stmts, total_stmts, f"{pct:.1f}%"))
    else:
        results.append((mod, 0, 0, "NO SVC"))

# Cleanup
for f in glob.glob(os.path.join(SERVER_DIR, "cov_*.out")):
    try: os.remove(f)
    except: pass

# Print results sorted by coverage ascending
def sort_key(x):
    p = x[3]
    if p in ("NO FILE", "NO SVC", "TIMEOUT", "ERR"):
        return -1
    return float(p.rstrip('%'))

results.sort(key=sort_key)

print(f"{'Module':<20} {'Covered':>8} {'Total':>8} {'Service%':>10}  {'Status'}")
print("-" * 68)
for mod, covered, total, pct_str in results:
    if pct_str in ("NO FILE", "NO SVC", "TIMEOUT", "ERR"):
        status = f"⚠️  {pct_str}"
    else:
        pct = float(pct_str.rstrip('%'))
        if pct >= 80:
            status = "✅ ≥80%"
        elif pct >= 50:
            status = "⚠️  50-80%"
        else:
            status = "❌ <50%"
    print(f"{mod:<20} {covered:>8} {total:>8} {pct_str:>10}  {status}")

print()
print("=" * 68)
print("ACTION PLAN: Modules with service coverage < 80%")
print("=" * 68)
target_reached = []
need_work = []
for mod, covered, total, pct_str in results:
    if pct_str in ("NO FILE", "NO SVC", "TIMEOUT", "ERR"):
        need_work.append((mod, covered, total, pct_str))
        continue
    pct = float(pct_str.rstrip('%'))
    if pct >= 80:
        target_reached.append(mod)
    else:
        need_work.append((mod, covered, total, pct_str))

for mod, covered, total, pct_str in need_work:
    if pct_str == "NO SVC":
        print(f"  {mod}: no service.go in package (only handler/repo/model)")
    elif pct_str in ("NO FILE", "TIMEOUT", "ERR"):
        print(f"  {mod}: {pct_str}")
    else:
        pct = float(pct_str.rstrip('%'))
        target_stmts = int(total * 0.8)
        gap = target_stmts - covered
        print(f"  {mod}: {pct_str} → need ~{gap} more stmts to reach 80%")

print(f"\nAlready ≥80%: {', '.join(target_reached) if target_reached else 'none'}")
print(f"Need work: {len(need_work)} modules")
