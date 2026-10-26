#!/usr/bin/env node
/*
 * 视觉契约检查（对应 specs/visual-redesign/specs 01、02）
 *
 * 两类检查：
 *   A. 构建产物含设计令牌与组件类（防止 globals.css 静默失效重演）
 *   B. 全量扫描原子色 / 功能性 emoji 残留，遗留文件走「基线棘轮」
 *
 * 为什么是全量扫描而不是白名单：
 *   初版只检查硬编码的 10 个主链路文件，新增页面默认免检。
 *   结果 admin/settings/page.tsx 带着 23 处原子色（含 text-blue-600）
 *   一路绿灯合入 —— 护栏对它本该防住的东西完全无感。
 *   现在改为：未登记的文件必须零残留；已登记的遗留文件只准减不准增。
 *
 * 用法：
 *   node scripts/check-design.mjs                    仅静态检查（B）
 *   node scripts/check-design.mjs --css              额外检查构建产物（A，需先 build）
 *   node scripts/check-design.mjs --update-baseline  重新生成遗留基线（清理旧色后运行）
 *
 * 退出码非 0 表示有未通过项，供 CI 使用。
 */

import { readFileSync, readdirSync, existsSync, writeFileSync } from 'node:fs'
import { join, dirname, sep } from 'node:path'
import { fileURLToPath } from 'node:url'

const webRoot = join(dirname(fileURLToPath(import.meta.url)), '..')
const src = join(webRoot, 'src')
const baselinePath = join(webRoot, 'scripts/design-baseline.json')

let failures = 0
const fail = (msg) => { console.error('  ✗ ' + msg); failures++ }
const ok = (msg) => console.log('  ✓ ' + msg)
const note = (msg) => console.log('  · ' + msg)

// 原子色类名：出现在 className 里的 Tailwind 调色板色
const ATOMIC_COLOR = /(?:bg|text|border|ring|from|to|via|shadow|fill|stroke|accent|divide|outline)-(blue|gray|rose|amber|emerald|indigo|purple|yellow|red|green|slate|zinc)-[0-9]{2,3}/g

// 功能性 emoji（心情 emoji 😊😌😢😠😰 属于内容，不在此列）
const FUNC_EMOJI = /[\u{1F44D}\u{1F4AC}\u{1F52E}\u{1F4A7}\u{1F4CC}\u{2B50}\u{1F3AD}\u{1F4ED}\u{1F680}\u{1F91D}\u{1F4D6}\u{1F4A1}\u{1F3C6}\u{1F9D8}\u{1F504}\u{1F4DD}\u{1F6E1}\u{1F511}]/gu

// ── 扫描所有 .tsx，返回 { relPath: {colors, emoji} } ──
const walkTsx = (dir) =>
  readdirSync(dir, { withFileTypes: true }).flatMap((e) =>
    e.isDirectory()
      ? walkTsx(join(dir, e.name))
      : e.name.endsWith('.tsx')
        ? [join(dir, e.name)]
        : []
  )

function scan() {
  const out = {}
  for (const abs of walkTsx(src).sort()) {
    const rel = abs.slice(src.length + 1).split(sep).join('/')
    const code = readFileSync(abs, 'utf8')
    const colors = [...code.matchAll(ATOMIC_COLOR)].map((m) => m[0])
    const emoji = [...code.matchAll(FUNC_EMOJI)].map((m) => m[0])
    out[rel] = { colors, emoji }
  }
  return out
}

const scanned = scan()

// ── --update-baseline：把当前残留量固化为基线 ──
if (process.argv.includes('--update-baseline')) {
  const next = {}
  for (const [rel, { colors, emoji }] of Object.entries(scanned)) {
    if (colors.length || emoji.length) {
      next[rel] = { colors: colors.length, emoji: emoji.length }
    }
  }
  writeFileSync(baselinePath, JSON.stringify(next, null, 2) + '\n', 'utf8')
  const totalC = Object.values(next).reduce((a, b) => a + b.colors, 0)
  const totalE = Object.values(next).reduce((a, b) => a + b.emoji, 0)
  console.log(`\n✓ 基线已更新：${Object.keys(next).length} 个文件，原子色 ${totalC} 处，功能 emoji ${totalE} 处`)
  console.log(`  ${baselinePath}\n`)
  process.exit(0)
}

const baseline = existsSync(baselinePath)
  ? JSON.parse(readFileSync(baselinePath, 'utf8'))
  : {}

console.log('\n[B] 原子色 / 功能 emoji 残留（spec 02，全量扫描 + 遗留基线棘轮）')

let clean = 0
let shrunk = 0
for (const [rel, { colors, emoji }] of Object.entries(scanned)) {
  const budget = baseline[rel] ?? { colors: 0, emoji: 0 }

  if (colors.length > budget.colors) {
    const uniq = [...new Set(colors)].slice(0, 5).join(', ')
    fail(
      budget.colors === 0
        ? `${rel} 原子色 ${colors.length} 处（新文件不允许任何残留）: ${uniq}${colors.length > 5 ? '…' : ''}`
        : `${rel} 原子色 ${colors.length} 处，超出基线 ${budget.colors}: ${uniq}`
    )
  } else if (colors.length < budget.colors) {
    note(`${rel} 原子色 ${colors.length} 处（已低于基线 ${budget.colors}，请跑 --update-baseline 下调）`)
    shrunk++
  }

  if (emoji.length > budget.emoji) {
    fail(`${rel} 功能 emoji ${emoji.length} 处（基线 ${budget.emoji}）: ${[...new Set(emoji)].join(' ')}`)
  } else if (emoji.length < budget.emoji) {
    shrunk++
  }

  if (!colors.length && !emoji.length) clean++
}
ok(`${clean} / ${Object.keys(scanned).length} 个文件零残留`)
if (shrunk) note(`${shrunk} 项低于基线，建议运行 --update-baseline 收紧棘轮`)

// ── A. 构建产物令牌检查 ──
if (process.argv.includes('--css')) {
  console.log('\n[A] 构建产物含设计令牌与组件类（spec 01）')
  const cssDir = join(webRoot, '.next/static/chunks')
  let css = ''
  if (existsSync(cssDir)) {
    for (const f of readdirSync(cssDir)) {
      if (f.endsWith('.css')) css += readFileSync(join(cssDir, f), 'utf8')
    }
  }
  if (!css) {
    fail('未找到 CSS 产物，请先运行 next build')
  } else {
    const before = failures
    const need = [
      ['--color-accent', /--color-accent:/],
      ['--color-canvas', /--color-canvas:/],
      ['.card', /\.card[{,: ]/],
      ['.btn-primary', /\.btn-primary[{,: ]/],
      ['.btn-secondary', /\.btn-secondary[{,: ]/],
      ['.input-field', /\.input-field[{,: ]/],
    ]
    for (const [name, re] of need) (re.test(css) ? ok : fail)(name)
    for (const t of ['emotion', 'insight', 'case', 'action', 'achievement']) {
      for (const layer of ['soft', 'line', 'ink']) {
        const key = `--color-${t}-${layer}`
        if (!css.includes(key + ':')) fail(`标签令牌 ${key} 缺失`)
      }
    }
    if (failures === before) ok('15 个标签令牌齐全')
  }
}

console.log('')
if (failures) {
  console.error(`✗ 视觉契约检查未通过：${failures} 项\n`)
  process.exit(1)
}
console.log('✓ 视觉契约检查全部通过\n')
