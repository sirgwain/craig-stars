
// return a promise that waits a certain number of ms
// use this to ensure some operations take a minimum amount of time
// (like very fast server responses)
export const wait = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));
