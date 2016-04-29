///<reference path="../typings/angularjs/angular.d.ts"/>
///<reference path="../typings/angular-material/angular-material.d.ts"/>

import {App} from "./app";
import {Header} from "./header";
import {Settings} from "./settings";
import {AogList} from "./aogList/aogList";

const app = angular.module("app", ["ngComponentRouter", "ngMaterial", "angularMoment"]);

app
    .config(["$mdThemingProvider", "$locationProvider", ($mdThemingProvider:angular.material.IThemingProvider, $locationProvider:ng.ILocationProvider) => {
        $mdThemingProvider
            .theme("default")
            .primaryPalette("grey", {
                "default": "50"
            })
            .accentPalette("orange");

        // $locationProvider.html5Mode(true);
    }])

    .component(App.name, App)
    .value("$routerRootComponent", App.name)

    .component(Header.name, Header)

    .component(Settings.name, Settings)

    .component(AogList.name, AogList);

angular.bootstrap(document, [app.name]);
