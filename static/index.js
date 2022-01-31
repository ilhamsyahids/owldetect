
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
        const matchesInput = []
        const matchesRef = []

        apiResp.data.matches.forEach(match => {
            matchesInput.push(match.input)
            matchesRef.push(match.ref)
        })

        matchesInput.sort((a, b) => b.start_idx - a.start_idx)
        matchesInput.forEach(matchInput => {
            divInputEl.innerHTML = outputMarkText(matchInput, divInputEl)
        })

        matchesRef.sort((a, b) => b.start_idx - a.start_idx)
        matchesRef.forEach(matchInput => {
            divRefEl.innerHTML = outputMarkText(matchInput, divRefEl)
        })

        alert("plagiarism found!");
    } else {
        divInputEl.classList.add("green");
        divRefEl.classList.add("green");
        alert("no plagiarism");
    }

    output.prepend(divRefEl);
    output.prepend(divInputEl);

    document.getElementById("clear").hidden = false;
}

function outputMarkText(matchInput, element) {
    const idx = element.innerHTML.indexOf(matchInput.text);

    const span = document.createElement("span");
    span.classList.add("red");
    span.innerHTML = matchInput.text;

    const front = element.innerHTML.substring(0, idx)
    const back = element.innerHTML.substring(idx + matchInput.text.length)
    return front + span.outerHTML + back;
}

function clearOutput() {
    const output = document.getElementById("output");
    output.innerHTML = "";
    document.getElementById("clear").hidden = true;
}
