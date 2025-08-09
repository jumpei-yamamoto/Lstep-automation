import fs from "node:fs";
import { execSync } from "node:child_process";

function diff(base, head) {
  return execSync(`git diff --name-only ${base} ${head}`, {
    encoding: "utf8",
  }).trim();
}

let base = "HEAD~1";
let head = "HEAD";

try {
  const evPath = process.env.GITHUB_EVENT_PATH;
  if (evPath && fs.existsSync(evPath)) {
    const ev = JSON.parse(fs.readFileSync(evPath, "utf8"));
    base = ev.pull_request?.base?.sha || base;
    head = ev.pull_request?.head?.sha || head;
  }
} catch {}

const files = diff(base, head).split("\n").filter(Boolean);

const coreChanged = files.some((f) =>
  /^backend\/internal\/(domain|usecase)\//.test(f)
);
if (!coreChanged) {
  console.log(
    "No core (backend/internal/{domain|usecase}) changes; skip test enforcement"
  );
  process.exit(0);
}

const hasGoTests = files.some((f) => f.endsWith("_test.go"));
if (!hasGoTests) {
  console.error(
    "❌ domain/usecase に変更がありますが、Goのユニットテスト差分が見当たりません。"
  );
  process.exit(2);
}
console.log("✅ tests changed along with core changes");
