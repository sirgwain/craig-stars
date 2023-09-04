export enum UserRole {
	user = 'user',
	admin = 'admin',
	guest = 'guest'
}

// The status of the user
// This is used by the login redirect process to determine if the user has been loaded
// from the server yet
export enum UserStatus {
	Unknown,
	LoggedIn,
	NotFound
}
export interface User {
	id?: number;
	username: string;
	password?: string; // guest users have a hash for the password
	role: UserRole;
	status: UserStatus;
	discordId?: string;
	discordAvatar?: string;
}

export const emptyUser: User = {
	username: '',
	role: UserRole.user,
	status: UserStatus.Unknown
};

export const userNotFound: User = {
	username: '',
	role: UserRole.user,
	status: UserStatus.NotFound
};
