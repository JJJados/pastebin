// CMPT315 - Pastebin
// Macewan University
// Jayden Laturnus

setupCreatePost();

function setupCreatePost() {
    const postSubmitButton: HTMLInputElement = document.querySelector("#create-paste-form");

    postSubmitButton.addEventListener("submit", createPaste);
}

async function createPaste(event: Event) {
    // Prevents the default form submission as we handle that here
    event.preventDefault();

    let url: string = "http://localhost:3333/api/v1/posts";
    let paste: string = newPaste();

    let options: Object = {
        method: "POST",
        headers: {
            "Accept": "application/json",
            "Content-Type": "application/json"
        },
        body: paste
    };

    let createMessage: HTMLParagraphElement = document.querySelector("#create-paste-message");

    try {
        let response = await fetch(url, options);
        let responseOK: boolean = response && response.ok;
        if (responseOK) {
            try {
                let data = await response.json();
                // Redirect to the post admin page
                window.location.replace(`http://localhost:3333/pastes/${data.adminAccessId}`);
            }
            catch(err) {
                console.log(err);
            }
        } else {
            createMessage.style.opacity = "1";
            createMessage.innerHTML = "Sorry, we could not create your paste."
        }
    }
    catch(err) {
        console.log(err);
    }
}