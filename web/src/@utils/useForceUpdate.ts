import React from "react";

export function useForceUpdate() {
  const [updateForced, updateState] = React.useState({});
  return [updateForced as never, React.useCallback(() => updateState({}), [])]
}
