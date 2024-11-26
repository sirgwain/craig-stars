export function asyncToVoidWrapper<T extends unknown[]>(
	asyncFn: (...args: T) => Promise<void>
): (...args: T) => void {
	return (...args: T) => {
		// Call the asynchronous function and handle the promise internally
		asyncFn(...args).catch((err) => {
			console.error('Error in async function:', err);
		});
	};
}
