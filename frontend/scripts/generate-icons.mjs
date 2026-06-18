// Generates the PWA PNG icons (no external deps): a dark tile with an amber
// barbell. Run with `node scripts/generate-icons.mjs`.
import { deflateSync } from "node:zlib";
import { writeFileSync, mkdirSync } from "node:fs";

const BG = [0x0b, 0x0b, 0x0f];
const FG = [0xf5, 0xa6, 0x23];

function crc32(buf) {
  let c = ~0;
  for (let i = 0; i < buf.length; i++) {
    c ^= buf[i];
    for (let k = 0; k < 8; k++) c = (c >>> 1) ^ (0xedb88320 & -(c & 1));
  }
  return ~c >>> 0;
}

function chunk(type, data) {
  const len = Buffer.alloc(4);
  len.writeUInt32BE(data.length, 0);
  const typeBuf = Buffer.from(type, "ascii");
  const crc = Buffer.alloc(4);
  crc.writeUInt32BE(crc32(Buffer.concat([typeBuf, data])), 0);
  return Buffer.concat([len, typeBuf, data, crc]);
}

function makePng(size) {
  const px = (x, y, c) => {
    const o = y * (size * 3 + 1) + 1 + x * 3;
    raw[o] = c[0];
    raw[o + 1] = c[1];
    raw[o + 2] = c[2];
  };
  // Raw image data: each scanline prefixed with a filter byte (0).
  const raw = Buffer.alloc(size * (size * 3 + 1));
  for (let y = 0; y < size; y++) for (let x = 0; x < size; x++) px(x, y, BG);

  const s = size;
  const barY1 = Math.round(s * 0.46), barY2 = Math.round(s * 0.54);
  const barX1 = Math.round(s * 0.22), barX2 = Math.round(s * 0.78);
  // Bar.
  for (let y = barY1; y < barY2; y++) for (let x = barX1; x < barX2; x++) px(x, y, FG);
  // Plates near each end.
  const plateY1 = Math.round(s * 0.32), plateY2 = Math.round(s * 0.68);
  const drawPlate = (cx) => {
    const px1 = Math.round(cx - s * 0.05), px2 = Math.round(cx + s * 0.05);
    for (let y = plateY1; y < plateY2; y++) for (let x = px1; x < px2; x++) px(x, y, FG);
  };
  drawPlate(barX1 + s * 0.04);
  drawPlate(barX2 - s * 0.04);

  const ihdr = Buffer.alloc(13);
  ihdr.writeUInt32BE(size, 0);
  ihdr.writeUInt32BE(size, 4);
  ihdr[8] = 8; // bit depth
  ihdr[9] = 2; // color type: truecolor RGB
  const png = Buffer.concat([
    Buffer.from([0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a]),
    chunk("IHDR", ihdr),
    chunk("IDAT", deflateSync(raw)),
    chunk("IEND", Buffer.alloc(0)),
  ]);
  return png;
}

mkdirSync(new URL("../public", import.meta.url), { recursive: true });
for (const size of [192, 512]) {
  const out = new URL(`../public/icon-${size}.png`, import.meta.url);
  writeFileSync(out, makePng(size));
  console.log(`wrote icon-${size}.png`);
}
