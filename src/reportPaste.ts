// CMPT315 - Pastebin
// Macewan University
// Jayden Laturnus

setupReportPost();

function setupReportPost() {
    const pasteReportButton: HTMLInputElement = document.querySelector("#report-paste-form");

    pasteReportButton.addEventListener("submit", reportPaste);
}

async function reportPaste(event: Event) {
    // Prevents the default form submission as we handle that here
    event.preventDefault();

    const reportForm: HTMLInputElement = document.querySelector("#report-paste-form");
    let readAccessID = reportForm.dataset.readId;

    let url: string = `http://localhost:3333/api/v1/posts/${readAccessID}/reports`;
    let report: string = newReport();

    let options: Object = {
        method: "POST",
        headers: {
            "Accept": "application/json",
            "Content-Type": "application/json"
        },
        body: report
    };

    let reportMessage: HTMLParagraphElement = document.querySelector("#report-paste-message");

    try {
        let response = await fetch(url, options);
        let responseOK: boolean = response && response.ok;
        if (responseOK) {
            // Go back to pastes after report
            window.location.replace(`http://localhost:3333/pastes`);
        } else {
            reportMessage.style.opacity = "1";
            reportMessage.innerHTML = "Sorry, your report did not go through."
        }
    }
    catch(err) {
        console.log(err);
    }
}