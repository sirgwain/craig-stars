import { MethodOptions_IdempotencyLevel } from '@bufbuild/protobuf/wkt';
import { Code, ConnectError, type Interceptor } from '@connectrpc/connect';
import { addError } from './Errors';

// the player num game requests act as, for hot seat games or admins viewing another player's game
let asPlayerNum: number | undefined;
// true when viewing a game without being able to make changes, i.e. viewing a submitted turn
let readOnly = false;

export function setAsPlayerNum(num: number | undefined) {
	asPlayerNum = num;
}

export function setReadOnly(value: boolean) {
	readOnly = value;
}

export const asPlayerInterceptor: Interceptor = (next) => async (req) => {
	// requests can act as a specific player, i.e. unsubmitting another hot seat player's turn
	if (asPlayerNum && !req.header.has('X-As-Player')) {
		req.header.set('X-As-Player', `${asPlayerNum}`);
	}
	if (readOnly) {
		// the server rejects changes in read only mode, but don't bother sending them
		if (req.method.idempotency !== MethodOptions_IdempotencyLevel.NO_SIDE_EFFECTS) {
			const err = new ConnectError('Read-only view, changes are not saved', Code.PermissionDenied);
			addError(err);
			throw err;
		}
		req.header.set('X-Read-Only', 'true');
	}
	return next(req);
};
