// CMPT315 - Pastebin
// Macewan University
// Jayden Laturnus

interface Post {
    readAccessId?: string,
    postTitle?: string;
    postContent?: string;
    publicAccess?: boolean;
    created?: string,
    updated?: string
}

interface PostReport {
    reported: boolean;
    reportedReason: string;
}

function newPaste(): string {
    let paste: Post = {
        postTitle: "",
        postContent: "",
        publicAccess: true
    }

    let pasteTitleElement: HTMLInputElement = document.querySelector("#paste-title");
    paste.postTitle = sanitize(pasteTitleElement.value);

    let pasteContentElement: HTMLTextAreaElement = document.querySelector("#paste-content-text");
    paste.postContent = sanitize(pasteContentElement.value);

    let privateCheckElement: HTMLInputElement = document.querySelector("#private-check");
    privateCheckElement.checked ? paste.publicAccess = false : paste.publicAccess = true;

    return JSON.stringify(paste);
}

function newReport(): string {
    let pasteReport: PostReport = {
        reported: true,
        reportedReason: ""
    }

    let reportContentElement: HTMLTextAreaElement = document.querySelector("#paste-report-text");
    pasteReport.reportedReason = sanitize(reportContentElement.value);

    return JSON.stringify(pasteReport);
}

/* 
sanitize sanitizes the provided string so it's safe to include it in the HTML
document.

Adapted from:
https://stackoverflow.com/questions/14129953/
how-to-encode-a-string-in-javascript-for-displaying-in-html 
*/
function sanitize(str: String): string {
    return String(str)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;');
}