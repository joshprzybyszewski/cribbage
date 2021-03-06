export type AlertType = 'error' | 'warning' | 'success' | 'info';

export const alertTypeToStyle = (t: AlertType) => `alert-${t}`;
