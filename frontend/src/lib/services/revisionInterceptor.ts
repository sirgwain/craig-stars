import { MethodOptions_IdempotencyLevel } from '@bufbuild/protobuf/wkt';
import { Code, ConnectError, type Interceptor } from '@connectrpc/connect';
import { staleData } from './Stale';

// This interceptor keeps a player from overwriting changes they made on another device.
//
// The server keeps a revision number for each player in a game. It increments the revision every
// time the player changes something (fleet orders, production, research, submitting a turn...) and
// returns the current revision in the X-Player-Revision header of every game response.
//
// We remember the revision we last saw, which tells us how new the data in this browser is:
//
// - When we read (i.e. the game status check every 30 seconds), we compare the revision in the
//   response to ours. If the server's is different, the player changed something somewhere else
//   and what we're showing is out of date, so we set staleData to show the reload prompt.
// - When we make a change, we send our revision with it. If it isn't the server's current revision
//   the server rejects the change with an Aborted error instead of saving it, and we show the
//   reload prompt. If it is, the server saves the change and returns the new revision, which
//   becomes ours.
//
// For example, with the game open on a phone and an iPad that both loaded revision 5:
//
//   phone:  update fleet orders, sending 5  -> saved, the phone is now at revision 6
//   iPad:   game status check               -> the response says 6 but we have 5, show the prompt
//   iPad:   update a planet, sending 5      -> rejected, the server is at 6
//   iPad:   reload                          -> loads everything, now at revision 6
const revisionHeader = 'X-Player-Revision';

type PlayerRevision = {
	// the last revision we saw from the server, undefined until the first response
	revision?: string;
	// the number of our changes that have been sent (or are queued) but haven't finished
	pending: number;
	// the number of our changes that have finished, successfully or not
	finished: number;
	// resolves when the last change we queued has finished, so the next one can wait for it
	queue: Promise<unknown>;
};

// one entry for each game and player we've made requests for
const revisions = new Map<string, PlayerRevision>();

// Forget the revisions we've seen and hide the reload prompt. This is called before (re)loading a
// game. We're about to load everything, so whatever revision the server has then is our new start.
export function resetPlayerRevisions() {
	revisions.clear();
	staleData.set(false);
}

export const revisionInterceptor: Interceptor = (next) => async (req) => {
	// only requests for a game have a player revision, the rest (races, users, techs) pass through
	const gameId = (req.message as { gameId?: bigint }).gameId;
	if (req.stream || !gameId) {
		return next(req);
	}

	// Hot seat players are different players on the server, each with their own revision, so track
	// them separately. This interceptor runs after the asPlayerInterceptor so the header is set.
	const key = `${gameId}:${req.header.get('X-As-Player') ?? ''}`;
	let state = revisions.get(key);
	if (!state) {
		state = { pending: 0, finished: 0, queue: Promise.resolve() };
		revisions.set(key, state);
	}
	// we hold onto this object for the rest of the request, even if resetPlayerRevisions() clears
	// the map in the meantime. A response from before a reload then can't affect the new state.
	const player = state;

	// Reads, i.e. loading the game, universe or a fleet. These never change the revision, we just
	// look at the one that comes back.
	if (req.method.idempotency === MethodOptions_IdempotencyLevel.NO_SIDE_EFFECTS) {
		// Our own changes move the revision too. If one is in flight while this read is, the read
		// could return the revision from before or after our change and we can't tell which, so it
		// would look like someone else changed something. Note whether any of our changes are
		// pending now and how many have finished, so we can tell after if one overlapped this read.
		const idle = player.pending === 0;
		const finished = player.finished;

		const res = await next(req);

		const revision = res.header.get(revisionHeader);
		// only trust the revision if none of our changes were pending when the read started, none
		// are pending now, and none started and finished in between
		if (revision && idle && player.pending === 0 && player.finished === finished) {
			if (player.revision === undefined) {
				// first response for this game, this is the revision of the data we're loading
				player.revision = revision;
			} else if (player.revision !== revision) {
				// the revision moved and it wasn't us, the player changed something on another device
				staleData.set(true);
			}
		}
		return res;
	}

	// Everything else is a change. Each change has to be sent with the revision returned by the
	// change before it, or the server would reject our own second change as stale. The UI doesn't
	// always wait for one change to finish before making another (i.e. quickly editing waypoints),
	// so changes wait their turn in a queue and only read player.revision once it's their turn.
	player.pending++;
	const result = player.queue.then(async () => {
		// the previous change has finished, so this is the latest revision
		if (player.revision !== undefined) {
			req.header.set(revisionHeader, player.revision);
		}
		try {
			const res = await next(req);
			// saved, the server incremented the revision and this is the new one
			player.revision = res.header.get(revisionHeader) ?? player.revision;
			return res;
		} catch (error) {
			if (error instanceof ConnectError) {
				if (error.code === Code.Aborted) {
					// The server has a newer revision than we sent, so it didn't save this change.
					// Keep our old revision so further changes are rejected too until we reload.
					staleData.set(true);
				} else {
					// the change failed for another reason (i.e. invalid orders), errors still
					// include the current revision
					player.revision = error.metadata.get(revisionHeader) ?? player.revision;
				}
			}
			throw error;
		} finally {
			player.pending--;
			player.finished++;
		}
	});
	// The next change waits for this one whether it succeeds or fails. The catch keeps a failed
	// change from failing the ones queued after it, the caller still gets the error from result.
	player.queue = result.catch(() => {});
	return result;
};
