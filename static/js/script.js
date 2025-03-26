document.addEventListener("DOMContentLoaded", function() {
    const qrForm = document.getElementById("qrForm");
    const urlForm = document.getElementById("urlForm");

    // QR Code Form Submission
    if (qrForm) {
        qrForm.addEventListener("submit", function(event) {
            const input = document.getElementById("qrInput").value;
            if (input.trim() === "") {
                alert("Please enter text to generate QR code.");
                event.preventDefault();
            }
        });
    }

    // URL Shortener Form Submission
    if (urlForm) {
        urlForm.addEventListener("submit", function(event) {
            const input = document.getElementById("urlInput").value;
            if (!input.startsWith("http://") && !input.startsWith("https://")) {
                alert("Please enter a valid URL starting with http:// or https://");
                event.preventDefault();
            }
        });
    }
});
