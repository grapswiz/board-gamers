self.addEventListener("push", (evt:any) => {
    console.log("push");
    
    if (!evt.data) {
        return;
    }

    const data = evt.data.json();
    evt.waitUntil(
        (<any>self).registration.showNotification(
            data.title,
            {
                body: data.body
            }
        )
    , false);
});

var clients;
self.addEventListener("notificationclick", (evt:any) => {
    evt.notification.close();

    evt.waitUntil(
        clients.matchAll({type: "window"})
            .then(() => {
                if (clients.openWindow) {
                    return clients.openWindow("https://board-gamers.appspot.com");
                }
            })
    );
});