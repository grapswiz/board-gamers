import IHttpService = angular.IHttpService;
import {Auth, User} from "./interfaces";
import ILocationService = angular.ILocationService;

export class HeaderController {
    auth: Auth;

    constructor(private $mdSidenav:angular.material.ISidenavService, private $http: IHttpService, private $location: ILocationService) {
        this.$http.get<Auth>("/api/v1/auth")
            .then((res) => {
                this.auth = res.data;
            });
    }

    toggleSidenav() {
        this.$mdSidenav("sidenav").toggle();
    }

    profileImage(user: User): string {
        if (!user) {
            return;
        }

        return "https" == this.$location.protocol() ? user.profileImageUrlHttps : user.profileImageUrl;
    }
}

export const Header = {
    name: "appHeader",
    controller: HeaderController,
    template: `
<md-toolbar layout="row">
    <div class="md-toolbar-tools">
        <md-button ng-click="$ctrl.toggleSidenav()" aria-label="toggle Sidenav" md-ink-ripple>
            <md-icon md-svg-src="img/icons/ic_menu_black_24px.svg" aria-label="Menu" hide-gt-md></md-icon>
        </md-button>
        <h1>
            <span>ボドゲ入荷速報(β)</span>
        </h1>
        <span flex></span>
        <md-button ng-if="!$ctrl.auth.isLoggedIn" aria-label="Login with Twitter"><a href="twitter/login">Twitterでログイン</a></md-button>
        <md-menu ng-if="$ctrl.auth.isLoggedIn">
            <img ng-src="{{::$ctrl.profileImage($ctrl.auth.user)}}" width="37px" height="37px" style="border-radius: 37px" ng-click="$mdOpenMenu($event)" aria-label="Open Menu">
            <md-menu-content>
                <md-menu-item>
                    <a class="md-button" href="twitter/logout">
                        ログアウト
                    </a>
                </md-menu-item>
            </md-menu-content>
        </md-menu>
    </div>
</md-toolbar>
    `
};
