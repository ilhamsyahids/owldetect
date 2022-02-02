
# Owl Detect

Two text plagiarism checker.

## Table of Contents

- [Demo](#demo)
- [Features](#features)
- [Prerequisites](#Prerequisites)
- [Improvement](#improvement)
    - [Server Side](#server-side)
    - [Client Side](#client-side)
- [Next Improvement](#next-improvement)

## Demo

[Demo here](http://owl-detect.herokuapp.com/)

## Features

- Check similarity sentence by tokenizing and fuzziness search within two text.
- Highlight plagiarism sentences.
- Start and end index of plagiarism sentences.
- Checking history
- Unit test and deploy on Heroku

## Prerequisites

- [Golang](https://go.dev/dl/) ^1.15

## Getting Started

- Install [prerequisites](#prerequisites)
- Start project
    ```
    go run main.go
    ```

## Improvement

### Server Side

High-level improvement:
1. For each document will tokenize to sentences, simply using regex to separate by dot, question mark, and exclamation mark.
2. Iterate through all the input sentences and compare them with every reference sentence.
3. Two sentences will check the similarity (detail below).
4. Reference sentences that are already stated as plagiarism will not be computed in other input sentences.
5. Lastly, will merge overlapped index intervals

For checking the similarity of two sentences, we could compare directly with equal, but it won't work well if at least one char is different.
Instead, here's the detailed algorithm:
1. Each sentence tokenizes it by word (split by space).
2. Compare every input word with all reference words.
Simply rule: we can use equals, but I change with fuzzy search.
The reason is to approach alter some char in a word.
To prevent the same word double compute, we store the result on map index.
3. Count matches words and percentage them.
4. If below with threshold (50%) then the sentences are not similar, otherwise similar.

This approach is far from perfect, kindly give me feedback.

### Client Side

1. Overall UI, adding style and re-structure HTML.
2. Highlight sentence that is similar with the same color both in input and reference.
3. Use a notification to give feedback that will discard 2s later.

## Next Improvement

Checker Algorithm:
- Using stop words to not check "unimportant word". But we need a lot of stop words.
- Split sentence with other considerations, not only dot, question mark, and exclamation mark.
- ~~Remove symbols such as dash, underscore, etc.~~

Frontend:
- Integration test
- Use local storage to store history
- Mobile friendly

Backend:
- ~~Integration test~~

Deployment:
- Add more deployment, such as using Docker
- Add more release platform, such as GCP, AWS, DO
