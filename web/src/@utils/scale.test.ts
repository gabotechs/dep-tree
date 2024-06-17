import { describe, expect, test } from "vitest";
import { scale } from "./scale.ts";

describe('scale', () => {
  it("1", [1, 0.5, 1.5, 0, 2], 1)
  it("2", [-2, 0.5, 1.5, 0, 2], 0)
  it("3", [10, 0.5, 1.5, 0, 2], 2)
  it("4", [.75, 0.5, 1, 1, 2], 1.5)
})

function it (name: string, v: [number, number, number, number, number], expected: number): void {
  test(name, () => {
    expect(scale(...v)).toEqual(expected);
  })
}
