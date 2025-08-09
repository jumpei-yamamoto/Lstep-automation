// Node 20+
// 目的: 'use client' な components/** で services/lib/apiClient import や fetch 直呼びを禁止
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const cfgPath = path.join(__dirname, "archcheck.config.yaml");
const cfgRaw = fs.readFileSync(cfgPath, "utf8");

function pickList(key) {
  const blockRe = new RegExp(`${key}:[\\s\\S]*?(?=\\n\\S|$)`, "m");
  const block = cfgRaw.match(blockRe)?.[0] || "";
  const itemRe = /-\s+["']?([^"'\n]+)["']?/g;
  const out = [];
  let m;
  while ((m = itemRe.exec(block)) !== null) out.push(m[1]);
  return out;
}
function pickScalar(key, def) {
  const re = new RegExp(`${key}:\\s*["']?([^\\n"']+)["']?`);
  const m = cfgRaw.match(re);
  return m ? m[1] : def;
}

const cli = process.argv.slice(2);
function arg(name, def) {
  const i = cli.findIndex((a) => a === name);
  if (i >= 0 && cli[i + 1]) return cli[i + 1];
  const kv = cli.find((a) => a.startsWith(name + "="));
  if (kv) return kv.split("=")[1];
  return def;
}

const frontendRoot = arg(
  "--frontend-root",
  pickScalar("frontend_root", "frontend")
);
const jsonOut = arg("--json-out", ".github/bots/out/archcheck-frontend.json");

const componentsDisallowImports = pickList("components_disallow_imports");
const componentsDisallowGlobals = pickList("components_disallow_globals");
const localStorageSeverity = pickScalar("localstorage_severity", "warn");

function listFiles(dir, exts = [".tsx", ".ts", ".jsx", ".js"]) {
  const res = [];
  if (!fs.existsSync(dir)) return res;
  function walk(d) {
    for (const ent of fs.readdirSync(d, { withFileTypes: true })) {
      const p = path.join(d, ent.name);
      if (ent.isDirectory()) walk(p);
      else if (exts.includes(path.extname(ent.name))) res.push(p);
    }
  }
  walk(dir);
  return res;
}

const componentsDir = path.join(frontendRoot, "components");
const files = listFiles(componentsDir);
const importRes = componentsDisallowImports.map((s) => new RegExp(s));
const globalRes = componentsDisallowGlobals.map(
  (s) => new RegExp(`\\b${s}\\s*\\(`)
);

const violations = [];
const warnings = [];

const hasUseClient = (src) => {
  const first5 = src
    .split("\n")
    .slice(0, 5)
    .map((l) => l.trim());
  return first5.includes("'use client'") || first5.includes('"use client"');
};

for (const f of files) {
  const src = fs.readFileSync(f, "utf8");
  if (!hasUseClient(src)) continue;

  // import checks
  const importLines = src.match(/import\s+[^;]+;?/g) || [];
  for (const line of importLines) {
    const m = line.match(/from\s+['"]([^'"]+)['"]/);
    const spec = m?.[1] ?? "";
    if (!spec) continue;
    for (const r of importRes) {
      if (r.test(spec)) {
        violations.push({ file: f, kind: "forbidden-import", detail: spec });
      }
    }
  }

  // fetch direct call
  for (const r of globalRes) {
    if (r.test(src)) {
      violations.push({ file: f, kind: "forbidden-global", detail: "fetch()" });
    }
  }

  // localStorage usage
  if (/\blocalStorage\b/.test(src)) {
    const msg = { file: f, kind: "localStorage", detail: "localStorage usage" };
    if (localStorageSeverity === "error") violations.push(msg);
    else warnings.push(msg);
  }
}

const out = { violations, warnings };
fs.mkdirSync(path.dirname(jsonOut), { recursive: true });
fs.writeFileSync(jsonOut, JSON.stringify(out, null, 2));

if (warnings.length) {
  console.log("⚠️ frontend warnings:");
  warnings.forEach((w) => console.log(" -", w.file, w.detail));
}
if (violations.length) {
  console.error("❌ frontend architecture violations:", violations.length);
  process.exit(2);
}
console.log("✅ archcheck-frontend: no violations");
