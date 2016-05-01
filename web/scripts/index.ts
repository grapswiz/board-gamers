///<reference path="../typings/angularjs/angular.d.ts"/>
///<reference path="../typings/angular-material/angular-material.d.ts"/>
///<reference path="../typings/angulartics/angulartics.d.ts"/>

import {App} from "./app";
import {Header} from "./header";
import {Settings} from "./settings";
import {AogList} from "./aogList/aogList";

const app = angular.module("app", ["ngComponentRouter", "ngMaterial", "angularMoment", "angulartics", "angulartics.google.analytics"]);

app
    .config(["$mdThemingProvider", "$locationProvider", "$analyticsProvider", ($mdThemingProvider:angular.material.IThemingProvider, $locationProvider:ng.ILocationProvider, $analyticsProvider:angulartics.IAnalyticsServiceProvider) => {
        $mdThemingProvider
            .theme("default")
            .primaryPalette("grey", {
                "default": "50"
            })
            .accentPalette("orange");

        // $locationProvider.html5Mode(true);

        $analyticsProvider.firstPageview(true); /* Records pages that don't use $state or $route */
        $analyticsProvider.withAutoBase(true);  /* Records full path */
    }])

    .component(App.name, App)
    .value("$routerRootComponent", App.name)

    .component(Header.name, Header)

    .component(Settings.name, Settings)

    .component(AogList.name, AogList);

angular.bootstrap(document, [app.name]);
