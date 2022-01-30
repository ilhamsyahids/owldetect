
async function analyzeText() {
    const inputText = document.getElementById("inputText");
    const refText = document.getElementById("refText");

    const response = await fetch("/analysis", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            "input_text": inputText.value,
            "ref_text": refText.value,
        })
    });

    const apiResp = await response.json();
    if (apiResp.data.matches) {
        alert("plagiarism found!");
    } else {
        alert("no plagiarism");
    }
}