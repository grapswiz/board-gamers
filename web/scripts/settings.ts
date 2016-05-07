import IWindowService = angular.IWindowService;
import IToastService = angular.material.IToastService;
import IHttpService = angular.IHttpService;
import IDialogService = angular.material.IDialogService;
export class SettingsController {
    isSupportedServiceWorker = true;

    isSubscribed: boolean;
    isSubscribedTrickplay = true;
    isSubscribedTendays = true;

    isSubscribedWithPush7 = false;

    isPushEnabled = false;
    useNotifications = false;

    constructor(private $http: IHttpService, private $window:IWindowService) {
        if (!("serviceWorker" in navigator)) {
            this.isSupportedServiceWorker = false;
            console.log("Service workers aren't supported in this browser.");
            return;
        }

        (<any>navigator).serviceWorker.register("serviceWorker.js").then((reg:any) => {
            if (reg.installing) {
                console.log("Service worker installing");
            } else if (reg.waiting) {
                console.log("Service worker installed");
            } else if (reg.active) {
                console.log("Service worker active");
            }

            this.initializeState(reg);
        });
    }

    $onInit() {
    }

    initializeState(reg: any) {
        if (!reg.showNotification) {
            console.log("Notifications aren't supported on service workers.");
            this.useNotifications = false;
        } else {
            this.useNotifications = true;
        }

        if ((<any>window).Notification.permission == "denied") {
            console.log("The user has blocked notification");
            return;
        }

        if (!("PushManager" in window)) {
            console.log("Push messaging isn't supported.");
            return;
        }

        (<any>navigator).serviceWorker.ready.then((reg:any) => {
            reg.pushManager.getSubscription()
                .then((subscription:any) => {
                    if (!subscription) {
                        console.log("Not yet subscribed to Push");
                        return;
                    }

                    console.log(subscription.toJSON());
                    const endpoint = subscription.endpoint;
                    const p256dh = subscription.getKey("p256dh");
                    const auth = subscription.getKey("auth");
                    this.updateStatus(endpoint, p256dh, auth, "init", this.shops());
                })
                .catch((err:any) => {
                    console.log("Error during getSubscription()", err);
                });

            //TODO
        });
    }

    subscribe() {
        (<any>navigator).serviceWorker.ready.then((reg:any) => {
            reg.pushManager.subscribe({userVisibleOnly: true})
                .then((subscription:any) => {
                    this.isPushEnabled = true;

                    const endpoint = subscription.endpoint;
                    const p256dh = subscription.getKey("p256dh");
                    const auth = subscription.getKey("auth");
                    this.updateStatus(endpoint, p256dh, auth, "subscribe", this.shops());
                })
                .catch((e:any) => {
                    if ((<any>window).Notification.permission == "denied") {
                        console.log("Permission for Notifications was denied");
                    } else {
                        console.log("Unable to subscribe to push.", e);
                    }
                });
        });
    }

    unsubscribe() {
        (<any>navigator).serviceWorker.ready.then((reg:any) => {
            reg.pushManager.getSubscription()
                .then((subscription:any) => {
                    if (!subscription) {
                        return;
                    }

                    const endpoint = subscription.endpoint;
                    const p256dh = subscription.getKey("p256dh");
                    const auth = subscription.getKey("auth");
                    this.updateStatus(endpoint, p256dh, auth, "unsubscribe", this.shops());

                    subscription.unsubscribe()
                        .then((successful:any) => {
                            this.isPushEnabled = false;
                        })
                        .catch((e:any) => {
                            console.log("unsubscription error: ", e);
                        });
                })
                .catch((e:any) => {
                    console.log("Error thrown while unsubscribing from push messaging.", e);
                });
        });
    }

    updateStatus(endpoint: string, p256dh: ArrayBuffer, auth: ArrayBuffer, statusType: string, shops: string[]) {
        if (statusType == "subscribe") {
            this.postSubscribeObj(statusType, endpoint, p256dh, auth, shops);
        } else if (statusType == "unsubscribe") {
            this.postSubscribeObj(statusType, endpoint, p256dh, auth, shops);
        } else if (statusType == "init") {
            this.isSubscribed = true;
        }
    }

    postSubscribeObj(statusType: string, endpoint: string, p256dh: ArrayBuffer, auth: ArrayBuffer, shops: string[]) {
        let data = {
            statusType: statusType,
            endpoint: endpoint,
            keys: {
                p256dh: btoa(String.fromCharCode.apply(null, new Uint8Array(p256dh))).replace(/\+/g, '-').replace(/\//g, '_'),
                auth: btoa(String.fromCharCode.apply(null, new Uint8Array(auth))).replace(/\+/g, '-').replace(/\//g, '_')
            },
            shops: shops
        };
        this.$http.post("/api/v1/subscription", JSON.stringify(data));
    }

    click() {
        if (this.isSubscribed) {
            this.subscribe();
        } else {
            this.unsubscribe();
        }
    }

    shops(): string[] {
        let shops = <string[]>[];
        if (this.isSubscribedTrickplay) {
            shops.push("トリックプレイ");
        }
        if (this.isSubscribedTendays) {
            shops.push("テンデイズ")
        }

        return shops;
    }

    goToPush7() {
        this.$window.location.href = "//board-gamers.app.push7.jp";
    }
}

export const Settings = {
    name: "settings",
    controller: SettingsController,
    template: `
        <md-list>
            <md-subheader class="md-no-sticky">設定</md-subheader>
            <md-list-item ng-show="$ctrl.isSupportedServiceWorker">
                <md-icon md-svg-src="img/icons/bell.svg" class="avatar"></md-icon>
                <p>通知を受け取る</p>
                <md-checkbox class="md-secondary" ng-model="$ctrl.isSubscribed" ng-change="$ctrl.click()"></md-checkbox>
            </md-list-item>
            <md-list-item ng-show="$ctrl.isSupportedServiceWorker && $ctrl.isSubscribed">
                <md-icon></md-icon>
                <p>トリックプレイ</p>
                <md-checkbox class="md-secondary" ng-model="$ctrl.isSubscribedTrickplay"></md-checkbox>
            </md-list-item>
            <md-list-item ng-show="$ctrl.isSupportedServiceWorker && $ctrl.isSubscribed">
                <md-icon></md-icon>
                <p>テンデイズ</p>
                <md-checkbox class="md-secondary" ng-model="$ctrl.isSubscribedTendays"></md-checkbox>
            </md-list-item>
            <md-list-item ng-show="!$ctrl.isSupportedServiceWorker" ng-click="$ctrl.goToPush7()">
                <md-icon></md-icon>
                <p>Push7で通知を受け取る</p>
            </md-list-item>
        </md-list>
    `
};
