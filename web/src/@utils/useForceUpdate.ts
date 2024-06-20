import React from "react";

export function useForceUpdate() {
  const [, updateState] = React.useState({});
  return React.useCallback(() => updateState({}), []);
}
