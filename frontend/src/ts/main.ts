console.log('Hello World!');

document.body.addEventListener('htmx:beforeSwap', function (evt: any) {
    // Allow 422 and 400 responses to swap
    // We treat these as form validation errors
    if (evt.detail.xhr.status === 422 || evt.detail.xhr.status === 400) {
        evt.detail.shouldSwap = true;
        evt.detail.isError = false;
    }
});