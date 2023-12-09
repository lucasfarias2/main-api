const socket = new WebSocket('ws://localhost:8080/ws');

socket.onmessage = function (event) {
    const message = event.data;
    fetch(url).then(response => response.text()).then(data => {
        updateContent("a", data);
    });
};

function updateContent(templateID, newContent) {
    // Update the specific part of your DOM with the new content
    // The logic here depends on how your templates map to your DOM
    // Example:
    document.getElementById("content").innerHTML = newContent;
}