// CMPT315 - Pastebin
// Macewan University
// Jayden Laturnus

setupUpdatePost();

function setupUpdatePost() {
    const updateSubmit: HTMLInputElement = document.querySelector("#update-paste-form");

    updateSubmit.addEventListener("submit", updatePaste);
}

async function updatePaste(event: Event) {
    // Prevents the default form submission as we handle that here
    event.preventDefault();

    const updateSubmit: HTMLInputElement = document.querySelector("#admin-buttons");
    let adminAccessID = updateSubmit.dataset.adminId;

    let url: string = `http://localhost:3333/api/v1/posts/${adminAccessID}`;
    let post: string = newPaste();

    let options: Object = {
        method: "PUT",
        headers: {
            "Accept": "application/json",
            "Content-Type": "application/json"
        },
        body: post
    };

    let updateMessage: HTMLParagraphElement = document.querySelector("#update-paste-message");

    try {
        let response = await fetch(url, options);
        let responseOK: boolean = response && response.ok;
        if (responseOK) {
            updateMessage.style.opacity = "1";
            updateMessage.innerHTML = "Successfully updated paste!"
        } else {
            updateMessage.style.opacity = "1";
            updateMessage.innerHTML = "Sorry, we could not update your paste."
        }
    }
    catch(err) {
        console.log(err);
    }
}
