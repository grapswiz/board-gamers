///<reference path="../typings/angularjs/angular.d.ts"/>
///<reference path="../typings/angular-material/angular-material.d.ts"/>

namespace app {
    "use strict";

    angular.module("app", ["ngComponentRouter", "ngMaterial"])
        .config(["$mdThemingProvider", "$locationProvider", ($mdThemingProvider:angular.material.IThemingProvider, $locationProvider:ng.ILocationProvider) => {
            $mdThemingProvider
                .theme("default")
                .primaryPalette("grey", {
                    "default": "50"
                })
                .accentPalette("orange");

            $locationProvider.html5Mode(true);
        }])

        .value("$routerRootComponent", "app")

        .component("app", {
            template: `
<app-header></app-header>
<div flex layout="row">
    <md-sidenav class="md-sidenav-left" md-is-locked-open="$mdMedia('gt-md')" md-component-id="sidenav">
        <md-content layout-padding>
            <md-list>
                <md-list-item>
                    <md-button aria-label="Home">
                        <md-icon md-svg-src="img/icons/ic_home_black_24px.svg" class="avatar"></md-icon>
                        ホーム
                    </md-button>
                </md-list-item>
                <md-list-item>
                    <md-button aria-label="Settings">
                        <md-icon md-svg-src="img/icons/ic_settings_black_24px.svg" class="avatar"></md-icon>
                        設定
                    </md-button>
                </md-list-item>
            </md-list>
        </md-content>
    </md-sidenav>
    <md-content flex id="content">
        <md-list>
            <md-subheader class="md-no-sticky">設定</md-subheader>
            <md-list-item>
                <md-icon md-svg-src="img/icons/bell.svg" class="avatar"></md-icon>
                <p>通知を受け取る</p>
                <md-checkbox class="md-secondary"></md-checkbox>
            </md-list-item>
            <md-list-item>
                <md-icon></md-icon>
                <p>トリックプレイ</p>
                <md-checkbox class="md-secondary"></md-checkbox>
            </md-list-item>
            <md-list-item>
                <md-icon></md-icon>
                <p>テンデイズ</p>
                <md-checkbox class="md-secondary"></md-checkbox>
            </md-list-item>
        </md-list>

        <md-list ng-if="arrivalOfGames" layout-wrap>
            <md-subheader class="md-no-sticky">最新の入荷情報</md-subheader>
            <md-list-item class="md-3-line md-long-text" ng-repeat="aog in ::arrivalOfGames">
                <div class="md-list-item-text" ng-cloak>
                    <h3>{{::aog.shop}}</h3>
                    <p><span ng-repeat="game in ::aog.games">{{::game | hankakuspace}}{{$last ? "" : ", "}}</span></p>
                    <div>{{::aog.createdAt}}</div>
                </div>
            </md-list-item>
        </md-list>
    </md-content>
</div>
            `,
            $routeConfig: []
        })

        .component("app-header", {
            template: `
    <md-toolbar layout="row">
    <div class="md-toolbar-tools">
        <md-button ng-click="toggleSidenav()" aria-label="toggle Sidenav" md-ink-ripple>
            <md-icon md-svg-src="img/icons/ic_menu_black_24px.svg" aria-label="Menu" hide-gt-md></md-icon>
        </md-button>
        <h1>
            <span>ボ</span>
        </h1>
        <span flex></span>
        <md-menu>
            <img ng-src="http://www.paper-glasses.com/api/twipi/{{::auth.user.screenName}}/normal" width="37px" height="37px" style="border-radius: 37px" ng-click="$mdOpenMenu($event)">
            <md-menu-content>
                <md-menu-item ng-if="!auth.isLoggedIn">
                    <a class="md-button" href="twitter/login">
                        ログイン
                    </a>
                </md-menu-item>
                <md-menu-item ng-if="auth.isLoggedIn">
                    <a class="md-button" href="twitter/logout">
                        ログアウト
                    </a>
                </md-menu-item>
            </md-menu-content>
        </md-menu>
    </div>
</md-toolbar>
        `
        });
}
