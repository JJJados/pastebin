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
setupGetPosts();
function setupGetPosts() {
    let nextPageButton = document.querySelector("#next-page-button");
    nextPageButton.addEventListener("click", incrementPage);
    let prevPageButton = document.querySelector("#prev-page-button");
    prevPageButton.disabled = true;
    prevPageButton.addEventListener("click", decrementPage);
    let pageButtonsDiv = document.querySelector("#page-buttons");
    let limit = pageButtonsDiv.dataset.limit;
    let offset = pageButtonsDiv.dataset.offset;
    getPosts(parseInt(limit), parseInt(offset));
}
function incrementPage() {
    let pageButtonsDiv = document.querySelector("#page-buttons");
    let curPage = pageButtonsDiv.dataset.curPage;
    let limit = pageButtonsDiv.dataset.limit;
    let offset = pageButtonsDiv.dataset.offset;
    let newCurPage = parseInt(curPage) + 1;
    let newOffset = parseInt(offset) + parseInt(limit);
    let prevPageButton = document.querySelector("#prev-page-button");
    prevPageButton.disabled = false;
    pageButtonsDiv.dataset.curPage = newCurPage.toString();
    pageButtonsDiv.dataset.offset = newOffset.toString();
    getPosts(parseInt(limit), newOffset);
}
function decrementPage() {
    let pageButtonsDiv = document.querySelector("#page-buttons");
    let curPage = pageButtonsDiv.dataset.curPage;
    let limit = pageButtonsDiv.dataset.limit;
    let offset = pageButtonsDiv.dataset.offset;
    let newCurPage = parseInt(curPage) - 1 > 0 ? parseInt(curPage) - 1 : 1;
    let newOffset = newCurPage > 1 ? parseInt(offset) - parseInt(limit) : 0;
    let prevPageButton = document.querySelector("#prev-page-button");
    if (newCurPage === 1) {
        prevPageButton.disabled = true;
    }
    pageButtonsDiv.dataset.curPage = newCurPage.toString();
    pageButtonsDiv.dataset.offset = newOffset.toString();
    getPosts(parseInt(limit), newOffset);
}
function convertTimestamps(data) {
    for (let i = 0; i < data.length; i++) {
        let timestampCreated = data[i].created;
        let timestampUpdated = data[i].updated;
        data[i].created = new Date(timestampCreated).toLocaleString();
        data[i].updated = new Date(timestampUpdated).toLocaleString();
    }
    return data;
}
function getPosts(limit, offset) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            let response = yield fetch(`http://localhost:3333/api/v1/posts?limit=${limit}&offset=${offset}`);
            let responseOK = response && response.ok;
            if (responseOK) {
                try {
                    let data = yield response.json();
                    convertTimestamps(data);
                    updatePostTable(data);
                }
                catch (err) {
                    console.log("here");
                    console.log(err);
                }
            }
            else {
                updatePostTable([]);
                let lastPageMessage = document.querySelector("#last-page-message");
                lastPageMessage.innerHTML = "Sorry, we couldn't retrieve any pastes.";
            }
        }
        catch (err) {
            console.log(err);
        }
    });
}
function updatePostTable(data) {
    let postsTemplate = document.querySelector("#pastes-template");
    let snippet = "";
    snippet = postsTemplate.text;
    let renderFn = doT.template(snippet);
    let renderedData = renderFn(data);
    let nextPageButton = document.querySelector("#next-page-button");
    if (data.length === 0) {
        nextPageButton.disabled = true;
    }
    else {
        nextPageButton.disabled = false;
    }
    let posts = document.querySelector("#pastes");
    posts.innerHTML = renderedData;
}
