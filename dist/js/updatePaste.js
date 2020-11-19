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
setupUpdatePost();
function setupUpdatePost() {
    const updateSubmit = document.querySelector("#update-paste-form");
    updateSubmit.addEventListener("submit", updatePaste);
}
function updatePaste(event) {
    return __awaiter(this, void 0, void 0, function* () {
        event.preventDefault();
        const updateSubmit = document.querySelector("#admin-buttons");
        let adminAccessID = updateSubmit.dataset.adminId;
        let url = `http://localhost:3333/api/v1/posts/${adminAccessID}`;
        let post = newPaste();
        let options = {
            method: "PUT",
            headers: {
                "Accept": "application/json",
                "Content-Type": "application/json"
            },
            body: post
        };
        let updateMessage = document.querySelector("#update-paste-message");
        try {
            let response = yield fetch(url, options);
            let responseOK = response && response.ok;
            if (responseOK) {
                updateMessage.style.opacity = "1";
                updateMessage.innerHTML = "Successfully updated paste!";
            }
            else {
                updateMessage.style.opacity = "1";
                updateMessage.innerHTML = "Sorry, we could not update your paste.";
            }
        }
        catch (err) {
            console.log(err);
        }
    });
}
