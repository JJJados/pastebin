// CMPT315 - Pastebin
// Macewan University
// Jayden Laturnus

setupGetPosts();

function setupGetPosts() {
    // Add listener for next page button
    let nextPageButton: HTMLButtonElement = document.querySelector("#next-page-button");
    nextPageButton.addEventListener("click", incrementPage);

    // Add listener for previous page button
    let prevPageButton: HTMLButtonElement = document.querySelector("#prev-page-button");
    prevPageButton.disabled = true;
    prevPageButton.addEventListener("click", decrementPage);

    // Get data from buttons div
    let pageButtonsDiv: HTMLDivElement = document.querySelector("#page-buttons");
    let limit: string = pageButtonsDiv.dataset.limit;
    let offset: string = pageButtonsDiv.dataset.offset;

    getPosts(parseInt(limit), parseInt(offset));
}


function incrementPage() {
    // Get data from buttons div
    let pageButtonsDiv: HTMLDivElement = document.querySelector("#page-buttons");
    let curPage: string | undefined = pageButtonsDiv.dataset.curPage;
    let limit: string | undefined = pageButtonsDiv.dataset.limit;
    let offset: string | undefined = pageButtonsDiv.dataset.offset;

    // Calculate new dataset values
    let newCurPage = parseInt(curPage) + 1;
    let newOffset = parseInt(offset) + parseInt(limit);

    let prevPageButton: HTMLButtonElement = document.querySelector("#prev-page-button");
    prevPageButton.disabled = false;

    // Assign new offset and page
    pageButtonsDiv.dataset.curPage = newCurPage.toString();
    pageButtonsDiv.dataset.offset = newOffset.toString();

    // Get the next set of posts
    getPosts(parseInt(limit), newOffset);
}

function decrementPage() {
    // Get data from buttons div
    let pageButtonsDiv: HTMLDivElement = document.querySelector("#page-buttons");
    let curPage: string = pageButtonsDiv.dataset.curPage;
    let limit: string = pageButtonsDiv.dataset.limit;
    let offset: string = pageButtonsDiv.dataset.offset;

    // Calculate new dataset values
    let newCurPage = parseInt(curPage) - 1 > 0 ? parseInt(curPage) - 1 : 1;
    let newOffset = newCurPage > 1 ? parseInt(offset) - parseInt(limit) : 0;

    let prevPageButton: HTMLButtonElement = document.querySelector("#prev-page-button");
    if (newCurPage === 1) {
        prevPageButton.disabled = true;
    }
    
    // Assign new offset and page
    pageButtonsDiv.dataset.curPage = newCurPage.toString();
    pageButtonsDiv.dataset.offset = newOffset.toString();
    // Get the next set of posts
    getPosts(parseInt(limit), newOffset);
}

function convertTimestamps(data: Post[]): Post[] {
    for (let i: number = 0; i < data.length; i++) {
        let timestampCreated: string = data[i].created;
        let timestampUpdated: string = data[i].updated;

        //let createdDate: string = new Date(timestampCreated).toLocaleString();
        //let updatedDate: string = new Date(timestampUpdated).toLocaleString();
        data[i].created = new Date(timestampCreated).toLocaleString();
        data[i].updated = new Date(timestampUpdated).toLocaleString();
    }
    return data;
}

async function getPosts(limit: number, offset: number) {
    try {
        let response = await fetch(`http://localhost:3333/api/v1/posts?limit=${limit}&offset=${offset}`);
        let responseOK: boolean = response && response.ok;

        if (responseOK) {
            try {
                let data = await response.json();
                // Convert postgres timestamps to date
                convertTimestamps(data);
                // Update posts template
                updatePostTable(data);
            }
            catch(err) {
                console.log("here")
                console.log(err);
            }
        } else {
            // Call on empty data set to populate last page elements
            updatePostTable([]);
            let lastPageMessage: HTMLParagraphElement = document.querySelector("#last-page-message");
            lastPageMessage.innerHTML = "Sorry, we couldn't retrieve any pastes."
        }
    }
    catch(err) {
        console.log(err);
    }
}

function updatePostTable(data: Object[]) {
    // Get the template from the DOM.
    let postsTemplate: HTMLScriptElement = document.querySelector("#pastes-template");

    let snippet: string = "";

    snippet = postsTemplate.text;

    // Create a render function for the template with doT.template.
    let renderFn = doT.template(snippet);
    // Use the render function to render the data.
    let renderedData = renderFn(data);

    let nextPageButton: HTMLButtonElement = document.querySelector("#next-page-button");
    if (data.length === 0) {
        nextPageButton.disabled = true;
    } else {
        nextPageButton.disabled = false;
    }

    // Insert the result into the DOM (inside the <div> with the ID log.
    let posts: HTMLDivElement = document.querySelector("#pastes");
    posts.innerHTML = renderedData;
}

