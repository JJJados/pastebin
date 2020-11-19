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
setupDeletePost();
function setupDeletePost() {
    const pasteDeleteButton = document.querySelector("#paste-delete");
    pasteDeleteButton.addEventListener("click", deletePaste);
}
function deletePaste(event) {
    return __awaiter(this, void 0, void 0, function* () {
        event.preventDefault();
        const adminButtons = document.querySelector("#admin-buttons");
        let adminAccessID = adminButtons.dataset.adminId;
        if (confirm("Are you sure you want to delete this paste?")) {
            let url = `http://localhost:3333/api/v1/posts/${adminAccessID}`;
            let options = {
                method: "DELETE"
            };
            let updateMessage = document.querySelector("#update-paste-message");
            try {
                let response = yield fetch(url, options);
                let responseOK = response && response.ok;
                if (responseOK) {
                    window.location.replace(`http://localhost:3333/pastes`);
                }
                else {
                    updateMessage.style.opacity = "1";
                    updateMessage.innerHTML = "Sorry, we could not delete your paste.";
                }
            }
            catch (err) {
                console.log(err);
            }
        }
    });
}
