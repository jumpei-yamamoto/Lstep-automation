// DangerJS: collect JSON outputs & format a single PR comment with fails/warns
import fs from "node:fs";
import path from "node:path";
import { warn, fail, markdown } from "danger";

const outDir = ".github/bots/out";
function readJSON(p) {
  try {
    return JSON.parse(fs.readFileSync(p, "utf8"));
  } catch {
    return null;
  }
}

const go = readJSON(path.join(outDir, "archcheck-go.json")) || {
  violations: [],
  warnings: [],
};
const fe = readJSON(path.join(outDir, "archcheck-frontend.json")) || {
  violations: [],
  warnings: [],
};
const sg = readJSON(path.join(outDir, "semgrep.json")) || { results: [] };

const grepTxtPath = path.join(outDir, "grep-violations.txt");
const grepTxt = fs.existsSync(grepTxtPath)
  ? fs.readFileSync(grepTxtPath, "utf8")
  : "";

let failCount = 0;
let warnCount = 0;

// archcheck-go
if (go.violations?.length) {
  failCount += go.violations.length;
  for (const v of go.violations) {
    fail(
      `[GO] ${v.relDir || v.package} imports forbidden "${v.import}" (rule: \`${
        v.rule
      }\`)`
    );
  }
}
for (const w of go.warnings || []) {
  warn(`[GO] ${w}`);
  warnCount++;
}

// archcheck-frontend
if (fe.violations?.length) {
  failCount += fe.violations.length;
  for (const v of fe.violations) {
    const place = v.file ? `\`${v.file}\`` : "";
    fail(`[FE] ${place} ${v.kind}: ${v.detail}`);
  }
}
for (const w of fe.warnings || []) {
  const place = w.file ? `\`${w.file}\`` : "";
  warn(`[FE] ${place} ${w.kind}: ${w.detail}`);
  warnCount++;
}

// semgrep
const semgrepRows = [];
for (const r of sg.results || []) {
  const sev = (r.extra?.severity || "WARNING").toUpperCase();
  const file = r.path;
  const line = r.start?.line;
  const msg =
    r.extra?.message || r.extra?.metavars
      ? JSON.stringify(r.extra?.metavars)
      : r.check_id;
  semgrepRows.push(
    `| ${sev} | \`${file}:${line}\` | ${msg} | \`${r.check_id}\` |`
  );
  if (sev === "ERROR") {
    fail(`[Semgrep] ${file}:${line} ${msg} (${r.check_id})`);
    failCount++;
  } else {
    warn(`[Semgrep] ${file}:${line} ${msg} (${r.check_id})`);
    warnCount++;
  }
}

// grep
if (grepTxt && grepTxt.trim().length > 0) {
  warn(`Grep checks:\n\n${"```\n" + grepTxt + "\n```"}`);
  warnCount++;
}

// summary table
const summary = [
  `**Architecture Guard Summary**`,
  ``,
  `- GO violations: **${go.violations?.length || 0}**, warnings: ${
    go.warnings?.length || 0
  }`,
  `- FE violations: **${fe.violations?.length || 0}**, warnings: ${
    fe.warnings?.length || 0
  }`,
  `- Semgrep findings: **${(sg.results || []).length}**`,
  ``,
  semgrepRows.length
    ? `### Semgrep Findings\n| Sev | Location | Message | Rule |\n|---|---|---|---|\n${semgrepRows.join(
        "\n"
      )}`
    : ``,
].join("\n");

markdown(summary);

// If nothing failed, leave CI green. Danger will fail the job automatically when any `fail()` was called.
