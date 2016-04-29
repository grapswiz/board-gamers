export class SettingsController {
    constructor() {
    }
}

export const Settings = {
    name: "settings",
    controller: SettingsController,
    template: `
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
    `
};
