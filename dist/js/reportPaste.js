"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
setupReportPost();
function setupReportPost() {
    const pasteReportButton = document.querySelector("#report-paste-form");
    pasteReportButton.addEventListener("submit", reportPaste);
}
function reportPaste(event) {
    return __awaiter(this, void 0, void 0, function* () {
        event.preventDefault();
        const reportForm = document.querySelector("#report-paste-form");
        let readAccessID = reportForm.dataset.readId;
        let url = `http://localhost:3333/api/v1/posts/${readAccessID}/reports`;
        let report = newReport();
        let options = {
            method: "POST",
            headers: {
                "Accept": "application/json",
                "Content-Type": "application/json"
            },
            body: report
        };
        let reportMessage = document.querySelector("#report-paste-message");
        try {
            let response = yield fetch(url, options);
            let responseOK = response && response.ok;
            if (responseOK) {
                window.location.replace(`http://localhost:3333/pastes`);
            }
            else {
                reportMessage.style.opacity = "1";
                reportMessage.innerHTML = "Sorry, your report did not go through.";
            }
        }
        catch (err) {
            console.log(err);
        }
    });
}
