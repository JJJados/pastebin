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
setupCreatePost();
function setupCreatePost() {
    const postSubmitButton = document.querySelector("#create-paste-form");
    postSubmitButton.addEventListener("submit", createPaste);
}
function createPaste(event) {
    return __awaiter(this, void 0, void 0, function* () {
        event.preventDefault();
        let url = "http://localhost:3333/api/v1/posts";
        let paste = newPaste();
        let options = {
            method: "POST",
            headers: {
                "Accept": "application/json",
                "Content-Type": "application/json"
            },
            body: paste
        };
        let createMessage = document.querySelector("#create-paste-message");
        try {
            let response = yield fetch(url, options);
            let responseOK = response && response.ok;
            if (responseOK) {
                try {
                    let data = yield response.json();
                    window.location.replace(`http://localhost:3333/pastes/${data.adminAccessId}`);
                }
                catch (err) {
                    console.log(err);
                }
            }
            else {
                createMessage.style.opacity = "1";
                createMessage.innerHTML = "Sorry, we could not create your paste.";
            }
        }
        catch (err) {
            console.log(err);
        }
    });
}
