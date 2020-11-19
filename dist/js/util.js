"use strict";
function newPaste() {
    let paste = {
        postTitle: "",
        postContent: "",
        publicAccess: true
    };
    let pasteTitleElement = document.querySelector("#paste-title");
    paste.postTitle = sanitize(pasteTitleElement.value);
    let pasteContentElement = document.querySelector("#paste-content-text");
    paste.postContent = sanitize(pasteContentElement.value);
    let privateCheckElement = document.querySelector("#private-check");
    privateCheckElement.checked ? paste.publicAccess = false : paste.publicAccess = true;
    return JSON.stringify(paste);
}
function newReport() {
    let pasteReport = {
        reported: true,
        reportedReason: ""
    };
    let reportContentElement = document.querySelector("#paste-report-text");
    pasteReport.reportedReason = sanitize(reportContentElement.value);
    return JSON.stringify(pasteReport);
}
function sanitize(str) {
    return String(str)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;');
}
