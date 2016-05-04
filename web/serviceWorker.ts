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
    , false)
});