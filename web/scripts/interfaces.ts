export interface Auth {
    isLoggedIn: boolean;
    user: User;
}

export interface User {
    userId: string;
    screenName: string;
    // shops: any;
    notificationKey: string;
    profileImageUrl: string;
    profileImageUrlHttps: string;
}

export interface ArrivalOfGame {
    id: string;
    shop: string;
    games: string[];
    createdAt: string;
    url: string;
}
