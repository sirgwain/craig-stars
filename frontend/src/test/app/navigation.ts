// Stub for $app/navigation in browser component tests, which run without SvelteKit.
// Tests that care about navigation can still vi.mock('$app/navigation').
const noop = () => {};

export const goto = async () => {};
export const invalidate = async () => {};
export const invalidateAll = async () => {};
export const preloadData = async () => ({ type: 'loaded', status: 200, data: {} });
export const preloadCode = async () => {};
export const refreshAll = async () => {};
export const pushState = noop;
export const replaceState = noop;
export const afterNavigate = noop;
export const beforeNavigate = noop;
export const onNavigate = noop;
export const disableScrollHandling = noop;
