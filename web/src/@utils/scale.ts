export function scale(n: number, lo: number, hi: number, tlo: number, thi: number): number {
  if (n < lo) {
    n = lo;
  }
  if (n > hi) {
    n = hi;
  }
  return (n - lo) / (hi - lo) * (thi - tlo) + tlo;
}
