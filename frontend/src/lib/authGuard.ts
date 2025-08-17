import { userClient } from './services/connect';
import { me } from './services/Stores';
import { userNotFound, UserSession, UserStatuses } from '$lib/types/User';

export async function authGuard(): Promise<UserSession | undefined> {
	try {
		const { user } = await userClient.getMe({});
		if (!user) {
			me.update(() => userNotFound);
			return;
		}
		// update the logged in user in the context
		const serverMe = Object.assign(new UserSession(), user);

		serverMe.status = UserStatuses.LoggedIn;
		me.update(() => serverMe);

		return serverMe;
	} catch {
		me.update(() => userNotFound);
	}
}
