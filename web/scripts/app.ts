import {Settings} from "./settings";
import {AogList} from "./aogList/aogList";

export class AppController {
    constructor() {
    }
}

export const App = {
    name: "app",
    controller: AppController,
    $routeConfig: [
        { path: "/list", name:"AogList", component: AogList.name, useAsDefault: true },
        { path: "/settings", name: "Settings", component: Settings.name }
    ],
    template: `
    <app-header></app-header>
    <div flex layout="row">
    <md-sidenav class="md-sidenav-left" md-is-locked-open="$mdMedia('gt-md')" md-component-id="sidenav">
        <md-content layout-padding>
            <md-list>
            <a ng-link="['AogList']">
                <md-list-item>
                    <md-button aria-label="Home">
                        <md-icon md-svg-src="img/icons/ic_home_black_24px.svg" class="avatar"></md-icon>
                        ホーム
                    </md-button>
                </md-list-item>
                </a>
                <a ng-link="['Settings']">
                <md-list-item>
                    <md-button aria-label="Settings">
                        <md-icon md-svg-src="img/icons/ic_settings_black_24px.svg" class="avatar"></md-icon>
                        設定                        
                    </md-button>                    
                </md-list-item>
                </a>
            </md-list>
        </md-content>
    </md-sidenav>
    <md-content flex id="content">
        <ng-outlet></ng-outlet>
    </md-content>
</div>
    `
};
