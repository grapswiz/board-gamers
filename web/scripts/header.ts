import IHttpService = angular.IHttpService;
import {Auth} from "./interfaces";

export class HeaderController {
    auth: Auth;

    constructor(private $mdSidenav:angular.material.ISidenavService, private $http: IHttpService) {
        this.$http.get<Auth>("/api/v1/auth")
            .then((res) => {
                this.auth = res.data;
            });
    }

    toggleSidenav() {
        this.$mdSidenav("sidenav").toggle();
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
            <span>ボ</span>
        </h1>
        <span flex></span>
        <md-menu>
            <img ng-src="http://www.paper-glasses.com/api/twipi/{{::$ctrl.auth.user.screenName}}/normal" width="37px" height="37px" style="border-radius: 37px" ng-click="$mdOpenMenu($event)">
            <md-menu-content>
                <md-menu-item ng-if="!$ctrl.auth.isLoggedIn">
                    <a class="md-button" href="twitter/login">
                        ログイン
                    </a>
                </md-menu-item>
                <md-menu-item ng-if="$ctrl.auth.isLoggedIn">
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
