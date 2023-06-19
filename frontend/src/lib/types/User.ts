export enum UserRole {
	admin,
	user
}

export interface User {
	id: number | undefined;
	username: string;
	role: UserRole;
}
