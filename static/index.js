
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
        let offsetInput = 0;
        let offsetRef = 0;

        apiResp.data.matches.forEach(match => {
            const lengthInputBefore = divInputEl.innerHTML.length;
            outputMarkText(offsetInput, match.input, divInputEl);
            const lengthInputAfter = divInputEl.innerHTML.length;
            offsetInput += lengthInputAfter - lengthInputBefore

            const lengthRefBefore = divRefEl.innerHTML.length;
            outputMarkText(offsetRef, match.ref, divRefEl);
            const lengthRefAfter = divRefEl.innerHTML.length;
            offsetRef += lengthRefAfter - lengthRefBefore
        })

        alert("plagiarism found!");
    } else {
        divInputEl.classList.add("green");
        divRefEl.classList.add("green");
        alert("no plagiarism");
    }

    output.prepend(divRefEl);
    output.prepend(divInputEl);
}

function outputMarkText(offset, matchInput, divElement) {
    const start = matchInput.start_idx + offset;
    const end = matchInput.end_idx + offset;

    const textValue = divElement.innerHTML

    const span = document.createElement("span");
    span.classList.add("red");
    span.innerHTML = textValue.substring(start, end);

    divElement.innerHTML = divElement.innerHTML.substring(0, start) + span.outerHTML + divElement.innerHTML.substring(end);
}