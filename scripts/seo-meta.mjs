// Deploy-time SEO head injection for prerendered Go pages.
import { readFileSync, writeFileSync, readdirSync, existsSync } from "node:fs";

const domain = process.argv[2] || "";
if (!domain || !existsSync("site")) {
  process.exit(0);
}

const pkg = JSON.parse(readFileSync("package.json", "utf8"));
const desc = String(pkg.description || "").replace(/\s+/g, " ").trim();

const esc = (value) =>
  String(value)
    .replace(/&/g, "&amp;")
    .replace(/"/g, "&quot;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");

for (const file of readdirSync("site").filter((entry) => entry.endsWith(".html"))) {
  const path = `site/${file}`;
  let html = readFileSync(path, "utf8");
  if (html.includes('property="og:title"')) {
    continue;
  }

  const match = html.match(/<title>([^<]*)<\/title>/);
  const title = match ? match[1].trim() : domain;
  const page = file === "index.html" ? "" : file;
  const url = `https://${domain}/${page}`;
  const tags = [
    `<meta name="description" content="${esc(desc)}">`,
    `<meta property="og:type" content="website">`,
    `<meta property="og:title" content="${esc(title)}">`,
    `<meta property="og:description" content="${esc(desc)}">`,
    `<meta property="og:url" content="${url}">`,
    `<meta property="og:site_name" content="Kinetic Gain">`,
    `<meta name="twitter:card" content="summary">`,
    `<meta name="twitter:title" content="${esc(title)}">`,
    `<meta name="twitter:description" content="${esc(desc)}">`
  ].join("\n  ");
  html = html.replace("</head>", `  ${tags}\n</head>`);
  writeFileSync(path, html);
}

