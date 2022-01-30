
async function analyzeText() {
    const inputText = document.getElementById("inputText");
    const refText = document.getElementById("refText");

    const output = document.getElementById("output");

    const inputTextValue = inputText.value;
    const refTextValue = refText.value;

    const divInputEl = document.createElement("div");
    divInputEl.innerHTML = inputTextValue;
    divInputEl.classList.add("output-item");

    const divRefEl = document.createElement("div");
    divRefEl.innerHTML = refTextValue;
    divRefEl.classList.add("output-item");

    output.appendChild(divInputEl);
    output.appendChild(divRefEl);

    const response = await fetch("/analysis", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            "input_text": inputTextValue,
            "ref_text": refTextValue,
        })
    });

    const apiResp = await response.json();
    if (apiResp.data.matches.length > 0) {
        divInputEl.classList.add("red");
        divRefEl.classList.add("red");
        alert("plagiarism found!");
    } else {
        divInputEl.classList.add("green");
        divRefEl.classList.add("green");
        alert("no plagiarism");
    }
}