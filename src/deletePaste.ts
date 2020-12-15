// CMPT315 - Pastebin
// Macewan University
// Jayden Laturnus

setupDeletePost();

function setupDeletePost() {
    const pasteDeleteButton: HTMLInputElement = document.querySelector("#paste-delete");

    pasteDeleteButton.addEventListener("click", deletePaste);
}

async function deletePaste(event: Event) {
    // Prevents the default form submission as we handle that here
    event.preventDefault();

    const adminButtons: HTMLDivElement = document.querySelector("#admin-buttons");
    let adminAccessID = adminButtons.dataset.adminId;

    if (confirm("Are you sure you want to delete this paste?")) {
        let url: string = `http://localhost:3333/api/v1/posts/${adminAccessID}`;
        let options: Object = {
            method: "DELETE"
        };

        let updateMessage: HTMLParagraphElement = document.querySelector("#update-paste-message");

        try {
            let response = await fetch(url, options);
            let responseOK: boolean = response && response.ok;
            if (responseOK) {
                // Go back to pastes after delete
                window.location.replace(`http://localhost:3333/pastes`);
            } else {
                // Use the same update message element for delete
                updateMessage.style.opacity = "1";
                updateMessage.innerHTML = "Sorry, we could not delete your paste."
            }
        }
        catch(err) {
            console.log(err);
        }
    } 
}
