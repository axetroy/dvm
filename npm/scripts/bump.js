const fs = require("fs");
const path = require("path");

const ref = process.env.GITHUB_REF;

const version = ref.replace(/^refs\/tags\/v/, "");

const pkgPath = path.join(__dirname, "..", "package.json");

const pkg = require(pkgPath);

pkg.version = version;

fs.writeFileSync(pkgPath, JSON.stringify(pkg, null, 2));
