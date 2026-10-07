import { ConnectError, type Interceptor } from '@connectrpc/connect';
import { versionUpdate } from './Version';

export const versionInterceptor: Interceptor = (next) => async (req) => {
	try {
		const res = await next(req);
		versionUpdate.observeServerVersion(res.header.get('x-app-version'));
		return res;
	} catch (error) {
		if (error instanceof ConnectError) {
			versionUpdate.observeServerVersion(error.metadata.get('x-app-version'));
		}
		throw error;
	}
};
