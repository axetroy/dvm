const path = require("path");
const os = require("os");
const download = require("download");
const ProgressBar = require("progress");

const pkg = require("../package.json");

function getArch() {
  const arch = os.arch();
  switch (arch) {
    case "ia32":
    case "x32":
      return "386";
    case "x64":
      return "amd64";
    case "arm":
      // @ts-expect-error ignore error
      const armv = process.config.variables.arm_version;

      if (!armv) return "armv7";

      return `armv${armv}`;
    default:
      return arch;
  }
}

function getPlatform() {
  const platform = os.platform();
  switch (platform) {
    case "win32":
      return "windows";
    default:
      return platform;
  }
}

function getDownloadURL(version) {
  const url = `https://github.com/axetroy/dvm/releases/download/${version}/dvm_${getPlatform()}_${getArch()}.tar.gz`;
  return url;
}

async function install(version) {
  const url = getDownloadURL(version);

  console.log(`Downloading '${url}'`);

  const bar = new ProgressBar("[:bar] :percent :etas", {
    complete: "=",
    incomplete: " ",
    width: 20,
    total: 0,
  });

  await download(url, path.join(__dirname, "..", "download"), {
    extract: true,
  }).on("response", (res) => {
    bar.total = res.headers["content-length"];
    res.on("data", (data) => bar.tick(data.length));
    res.on("error", (err) => {
      console.error("error");
    });
  });
}

install("v" + pkg.version).catch((err) => {
  throw err;
});
