
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
            const color = randomHexColor()
            matchesInput.push({ text: match.input, color })
            matchesRef.push({ text: match.ref, color })
        })

        highlightMatchesText(matchesInput, divInputEl)
        highlightMatchesText(matchesRef, divRefEl)

        divInputEl.classList.add("red")
        divRefEl.classList.add("red")

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

function outputMarkText(matchInput, element, color) {
    const idx = element.innerHTML.indexOf(matchInput.text);

    const span = document.createElement("span");
    span.innerHTML = matchInput.text;
    span.style.color = color

    const front = element.innerHTML.substring(0, idx)
    const back = element.innerHTML.substring(idx + matchInput.text.length)
    return front + span.outerHTML + back;
}

function clearOutput() {
    const output = document.getElementById("output");
    output.innerHTML = "";
    document.getElementById("clear").hidden = true;
}

function randomHexColor() {
    return `#${Math.random().toString(16).slice(2, 8)}`;
}

function highlightMatchesText(matches, element) {
    matches.sort((a, b) => b.text.start_idx - a.text.start_idx)
    matches.forEach(match => {
        const color = match.color
        const text = match.text

        element.innerHTML = outputMarkText(text, element, color)
    })
}