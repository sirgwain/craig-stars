import { MethodOptions_IdempotencyLevel } from '@bufbuild/protobuf/wkt';
import { Code, ConnectError, type Interceptor } from '@connectrpc/connect';

// Ride out a server restart during deploy, when nginx returns 502/503 (Unavailable).
// Only RPCs marked with an idempotency_level in the proto are retried, so a
// request that did reach the server can't be applied twice.
export function createRetryInterceptor(delaysMs = [500, 1000, 2000, 2000]): Interceptor {
	return (next) => async (req) => {
		if (
			req.stream ||
			req.method.idempotency === MethodOptions_IdempotencyLevel.IDEMPOTENCY_UNKNOWN
		) {
			return next(req);
		}

		for (const delay of delaysMs) {
			try {
				return await next(req);
			} catch (error) {
				if (
					!(error instanceof ConnectError) ||
					error.code !== Code.Unavailable ||
					req.signal.aborted
				) {
					throw error;
				}
			}
			await new Promise((resolve) => setTimeout(resolve, delay));
		}
		return next(req);
	};
}

export const retryInterceptor = createRetryInterceptor();
