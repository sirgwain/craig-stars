import { userNotFound, UserSession, UserStatuses } from '$lib/types/User';
import { me } from './services/Stores';

export async function authGuard(): Promise<UserSession | undefined> {
	const response = await fetch(`/api/me`, {
		method: 'GET',
		headers: {
			accept: 'application/json'
		}
	});

	if (!response.ok) {
		// no user
		me.update(() => userNotFound);
	} else {
		// update the logged in user in the context
		const userSession = (await response.json()) as UserSession;
		const user = Object.assign(new UserSession(), userSession);

		user.status = UserStatuses.LoggedIn;
		me.update(() => user);

		return user;
	}
}
