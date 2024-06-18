export function HSVtoRGB(h: number, s: number, v: number): [number, number, number] {
  if (h < 0 || h >= 360 || s < 0 || s > 1 || v < 0 || v > 1) {
    return [0, 0, 0];
  }

  const C = v * s;
  const X = C * (1 - Math.abs((h / 60) % 2 - 1));
  const m = v - C;

  let Rnot, Gnot, Bnot;
  if (h >= 0 && h < 60) {
    [Rnot, Gnot, Bnot] = [C, X, 0];
  } else if (h >= 60 && h < 120) {
    [Rnot, Gnot, Bnot] = [X, C, 0];
  } else if (h >= 120 && h < 180) {
    [Rnot, Gnot, Bnot] = [0, C, X];
  } else if (h >= 180 && h < 240) {
    [Rnot, Gnot, Bnot] = [0, X, C];
  } else if (h >= 240 && h < 300) {
    [Rnot, Gnot, Bnot] = [X, 0, C];
  } else if (h >= 300 && h < 360) {
    [Rnot, Gnot, Bnot] = [C, 0, X];
  } else {
    throw new Error('unreachable')
  }

  const r = Math.round((Rnot + m) * 255);
  const g = Math.round((Gnot + m) * 255);
  const b = Math.round((Bnot + m) * 255);

  return [r, g, b];
}
