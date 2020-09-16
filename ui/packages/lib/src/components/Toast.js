import React, { useState } from "react";
import { EuiGlobalToastList } from "@elastic/eui";
import isEqual from "react-fast-compare";

let addToastHandler;
let removeAllToastsHandler;

// Other components can add toasts directly
export const addToast = toast => {
  addToastHandler(toast);
};

// Other components can remove all toasts directly
export const removeAllToasts = () => {
  removeAllToastsHandler();
};

export const Toast = ({ toastLifeTimeMs = 15000 }) => {
  const [toasts, setToasts] = useState([]);

  const addToast = toast => {
    !toasts.find(t => isEqual(t, toast)) && setToasts(toasts.concat(toast));
  };

  const removeToast = removedToast =>
    setToasts(toasts.filter(toast => toast.id !== removedToast.id));

  const removeAllToasts = () => setToasts([]);

  addToastHandler = addToast;
  removeAllToastsHandler = removeAllToasts;

  return (
    <EuiGlobalToastList
      toasts={toasts}
      dismissToast={removeToast}
      toastLifeTimeMs={toastLifeTimeMs}
    />
  );
};
